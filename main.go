package main

import (
	"api-3390/config"
	"api-3390/user"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const UploadPath = "./uploads/"

func main() {
	os.MkdirAll(UploadPath, os.ModePerm)
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := cfg.Connection()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(user.UserTable)
	if err != nil {
		log.Fatal(err)
	}
	as := user.AuthService{DB: db}
	//fs := user.FileService{DB: db}
	r := chi.NewRouter()
	cfg.ApplicationMiddleWare(r)
	r.Get("/", renderUploadForm)
	r.Post("/api/files/", handleFileUpload)
	r.Post("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var u user.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		b, err := as.ExistsUser(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !b {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		authenticate, err := as.AuthenticateUser(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !authenticate {
			http.Error(w, "user not authorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	r.Post("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		var u user.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, fmt.Errorf("failed to register user: %w", err).Error(), http.StatusBadRequest)
			return
		}
		b, err := as.ExistsUser(&u)
		if err != nil {
			http.Error(w, fmt.Errorf("failed to register user: %w", err).Error(), http.StatusBadRequest)
			return
		}
		if b {
			http.Error(w, "failed to register user: user already exists", http.StatusConflict)
			return
		}
		if err := as.RegisterUser(&u); err != nil {
			http.Error(w, fmt.Errorf("failed to register user: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
	log.Println(fmt.Sprintf("Starting server on: '%s'", cfg.Address()))
	if err := http.ListenAndServe(cfg.Address(), r); err != nil {
		log.Fatal(err)
	}
}
func renderUploadForm(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering upload form")
	tmpl := `
        <!DOCTYPE html>
        <html>
        <body>
            <h2>Upload File</h2>
            <form action="/upload" method="post" enctype="multipart/form-data">
                Select file: <input type="file" name="file"><br><br>
                <input type="submit" value="Upload">
            </form>
        </body>
        </html>
    `
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// handleFileUpload processes the uploaded file
func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := saveFile(file, handler.Filename); err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)
}

// saveFile writes the uploaded file to the server
func saveFile(file multipart.File, filename string) error {
	// Read the file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Create and write the file
	outFile, err := os.Create(UploadPath + filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := outFile.Write(fileBytes); err != nil {
		return err
	}
	return nil
}
func getConfig() (*config.Config, error) {
	if len(os.Args) > 1 {
		cfg, err := config.Load(os.Args[1])
		return cfg, err
	}
	cfg, err := config.LoadFromEnv()
	return cfg, err

}
