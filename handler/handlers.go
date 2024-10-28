package handler

import (
	"api-3390/container"
	"api-3390/container/predicate"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// User Handlers
func (a *API) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	us, err := a.Services.UserService.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(us); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (a *API) HandleUpdateUserById(w http.ResponseWriter, r *http.Request) {
	id, err := getStringId("user_id", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := a.Services.UserService.GetUserById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if u == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updatedUser container.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedUser.ID = id
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "ErrorMessage hashing password", http.StatusInternalServerError)
			return
		}
		updatedUser.Password = string(hashedPassword)
	}

	err = a.Services.UserService.UpdateUser(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (a *API) HandleDeleteUserById(w http.ResponseWriter, r *http.Request) {
	id, err := getStringId("user_id", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := a.Services.UserService.DeleteUserById(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (a *API) HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := getStringId("user_id", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := a.Services.UserService.GetUserById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (a *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user container.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.Services.UserService.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Auth Handlers
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	u, err := a.Services.AuthService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// File Handlers
func (a *API) HandleGetAllFiles(w http.ResponseWriter, r *http.Request) {
	us, err := a.Services.FileService.GetAllFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(us); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (a *API) HandleCreateFile(fileTypeMap map[string][]predicate.Predicate[io.Reader], idPredicates []predicate.Predicate[string]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userid := r.FormValue("userid")
		for _, p := range idPredicates {
			if !p.Test(userid) {
				http.Error(w, p.ErrorMessage(userid), http.StatusBadRequest)
				return
			}
		}
		parsedId, err := strconv.ParseUint(userid, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "unable to create file", http.StatusBadRequest)
			return
		}
		var f = &container.File{
			UserID: uint32(parsedId),
			Name:   fileHeader.Filename,
		}
		if b, err := a.Services.FileService.UserHasFileEntry(f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if !b {
			if err := a.Services.FileService.CreateFileEntry(f); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if err := a.Services.FileService.UpdateFileEntry(f); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "unable to create file", http.StatusBadRequest)
			return
		}
		filePath := filepath.Join("./uploads", userid)
		path := filepath.Join(filePath, fileHeader.Filename)

		ext := filepath.Ext(path)
		if _, exists := fileTypeMap[ext]; !exists {
			http.Error(w, "file type not supported", http.StatusBadRequest)
			return
		}

		var predicates = fileTypeMap[ext]
		for _, p := range predicates {
			d := io.NopCloser(bytes.NewReader(data))
			if !p.Test(d) {
				http.Error(w, p.ErrorMessage(d), http.StatusBadRequest)
				return
			}
		}
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		outFile, err := os.Create(path)
		if err != nil {
			http.Error(w, "Unable to create file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		d := io.NopCloser(bytes.NewReader(data))

		_, err = io.Copy(outFile, d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = fmt.Fprintf(w, "File uploaded successfully: %s", fileHeader.Filename)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
func (a *API) HandleGetFileById(w http.ResponseWriter, r *http.Request) {
	id, err := getStringId("file_id", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file, err := a.Services.FileService.GetFileById(id)
	if file == nil || err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	filePath := filepath.Join("./uploads", strconv.Itoa(int(file.UserID)), file.Name)
	columnsParam := r.URL.Query().Get("columns")
	columns := strings.Split(columnsParam, ",")
	p := QueryParams{
		Operation: r.URL.Query().Get("operation"),
		Column:    columns,
	}
	qb := NewQueryBuilder().
		AddQuery("sum", a.calculateSum).
		SetDefaultCase(func(w http.ResponseWriter, _ []string, filePath string) {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filePath))
			w.Header().Set("Content-Type", "application/octet-stream")
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				http.Error(w, "file not found", http.StatusNotFound)
				return
			}
			f, err := os.Open(filePath)
			if err != nil {
				http.Error(w, "unable to open file", http.StatusInternalServerError)
				return
			}
			defer f.Close()
			if _, err := io.Copy(w, f); err != nil {
				http.Error(w, "unable to send file", http.StatusInternalServerError)
			}
		})
	qb.Build(w, p, filePath)
}
func (a *API) calculateSum(w http.ResponseWriter, columnName []string, filePath string) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read the header
	headers, err := reader.Read()
	if err != nil {
		http.Error(w, "error reading file", http.StatusInternalServerError)
		return
	}

	// Find the index of the specified column
	columnIndex := -1
	for i, header := range headers {
		if strings.EqualFold(header, columnName[0]) {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		http.Error(w, "column not found", http.StatusBadRequest)
		return
	}

	// Variables for total sum and count
	var total float64
	var count int

	// Read the records and calculate the total
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			http.Error(w, "error reading file", http.StatusInternalServerError)
			return
		}

		// Parse the value in the specified column
		value, err := strconv.ParseFloat(record[columnIndex], 64)
		if err != nil {
			http.Error(w, "invalid data format", http.StatusBadRequest)
			return
		}

		total += value
		count++
	}

	// Return the sum as JSON
	response := map[string]float64{"sum": total}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (a *API) HandleDeleteUserFileByName(w http.ResponseWriter, r *http.Request) {
	userid, err := getStringId("user_id", r)
	fileName := r.Context().Value("file_name").(string)
	file, err := a.Services.FileService.GetUserFileByName(userid, fileName)
	if file == nil || err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	filePath := filepath.Join("./uploads", strconv.Itoa(int(userid)), file.Name)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	err = os.Remove(filePath)
	if err != nil {
		http.Error(w, "unable to delete file", http.StatusInternalServerError)
		return
	}
	fmt.Println(file.ID)
	if err := a.Services.FileService.DeleteFileById(file.ID); err != nil {
		http.Error(w, "unable to delete file", http.StatusInternalServerError)
		return
	}
}

func (a *API) HandleGetUserFileByName(w http.ResponseWriter, r *http.Request) {
	userid, err := getStringId("user_id", r)
	fileName := r.Context().Value("file_name").(string)
	file, err := a.Services.FileService.GetUserFileByName(userid, fileName)
	if file == nil || err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	filePath := filepath.Join("./uploads", strconv.Itoa(int(userid)), file.Name)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	columnsParam := r.URL.Query().Get("columns")
	columns := strings.Split(columnsParam, ",")
	p := QueryParams{
		Operation: r.URL.Query().Get("operation"),
		Column:    columns,
	}
	qb := NewQueryBuilder().
		AddQuery("sum", a.calculateSum).
		SetDefaultCase(func(w http.ResponseWriter, _ []string, filePath string) {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filePath))
			w.Header().Set("Content-Type", "application/octet-stream")
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				http.Error(w, "file not found", http.StatusNotFound)
				return
			}
			f, err := os.Open(filePath)
			if err != nil {
				http.Error(w, "unable to open file", http.StatusInternalServerError)
				return
			}
			defer f.Close()
			if _, err := io.Copy(w, f); err != nil {
				http.Error(w, "unable to send file", http.StatusInternalServerError)
			}
		})
	qb.Build(w, p, filePath)
}
func (a *API) HandleUpdateFileById(w http.ResponseWriter, r *http.Request) {
	id, err := getStringId("file_id", r)
	fmt.Println(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := a.Services.FileService.GetFileById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if u == nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	var updatedFile container.File
	err = json.NewDecoder(r.Body).Decode(&updatedFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	updatedFile.ID = id
	if err := a.Services.FileService.UpdateFileEntry(&updatedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *API) HandleGetUserFiles(w http.ResponseWriter, r *http.Request) {
	id, err := getStringId("user_id", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	files, err := a.Services.FileService.GetUserFiles(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Helper Functions
func getStringId(key string, r *http.Request) (uint32, error) {
	val, ok := r.Context().Value(key).(string)
	if !ok {
		return 0, errors.New("un defined value")
	}
	id, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}
