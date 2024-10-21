package handler

import (
	"api-3390/container"
	"api-3390/container/predicate"
	"api-3390/user"
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

func (a *API) HandleCreateFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userid := r.FormValue("userid")
	if userid == "" {
		http.Error(w, "userid is required", http.StatusBadRequest)
		return
	}
	b := predicate.AllowedCharacters.Test(userid) && predicate.NonNegative.Test(userid)
	if !b {
		http.Error(w, "userid is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(userid, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "unable to create file", http.StatusInternalServerError)
		return
	}
	var f = &container.File{
		UserID: uint32(id),
		Name:   fileHeader.Filename,
	}
	filePath := filepath.Join("./uploads", userid) // Change as needed
	path := filepath.Join(filePath, f.Name)
	if !validFileExtension.Test(path) {
		http.Error(w, validFileExtension.ErrorMessage(f.Name), http.StatusBadRequest)
		return
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
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintf(w, "File uploaded successfully: %s", fileHeader.Filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
}
func parseQueryParams(r *http.Request) *FileQueryParams {
	return &FileQueryParams{
		Column:    r.URL.Query().Get("column"),
		Operation: r.URL.Query().Get("operation"),
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

	p := parseQueryParams(r)

	if p.IsEmpty() {
		a.handleGetFileById(w, r, file, filePath)
	} else {
		switch p.Operation {
		case "average":
			a.calculateAverage(w, p.Column, filePath)
		default:
			http.Error(w, "invalid operation", http.StatusBadRequest)
			return
		}
	}
}
func (a *API) calculateAverage(w http.ResponseWriter, columnName string, filePath string) {

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		http.Error(w, "error reading file", http.StatusInternalServerError)
		return
	}

	columnIndex := -1
	for i, header := range headers {
		if strings.EqualFold(header, columnName) {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		http.Error(w, "column not found", http.StatusBadRequest)
		return
	}

	var total float64
	var count int

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
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

	if count == 0 {
		http.Error(w, "no valid data to average", http.StatusBadRequest)
		return
	}

	average := total / float64(count)

	// Return the average as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"average": %f}`, average)
}
func (a *API) handleGetFileById(w http.ResponseWriter, r *http.Request, file *container.File, filePath string) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
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
}
func (a *API) HandleDeleteFileById(w http.ResponseWriter, r *http.Request) {
	id, err := getStringId("file_id", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := a.Services.FileService.DeleteFileById(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

var validFileExtension = predicate.Predicate[string]{
	Test: func(t string) bool {
		ext := strings.ToLower(filepath.Ext(t))
		for _, validExt := range user.ValidExtensions {
			if ext == validExt {
				return true
			}
		}
		return false
	},
	ErrorMessage: predicate.ErrorMessage(fmt.Sprintf("the file is not an accepted file format, these are the accepted files: %s", user.ValidExtensions)),
}
