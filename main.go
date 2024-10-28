package main

import (
	"api-3390/config"
	"api-3390/const"
	"api-3390/container/predicate"
	"api-3390/handler"
	"api-3390/handler/middleware"
	"api-3390/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
	"strings"
)

const UploadPath = "./uploads/"

func main() {
	os.MkdirAll(UploadPath, os.ModePerm)
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := cfg.Connection()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatalf("Error enabling foreign keys: %v", err)
	}
	_, err = db.Exec(constants.UserTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(constants.UserFileTable)
	if err != nil {
		log.Fatal(err)
	}
	services := handler.NewServices(service.NewAuthService(db), service.NewFileService(db), service.NewUserService(db))
	api := handler.API{Services: services}
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-API-KEY"},
		AllowCredentials: false,
		MaxAge:           300, // Max cache age in seconds ??
	}))
	r.Use(applicationMiddleWare(cfg))

	// Auth Routes
	r.Route("/login", func(r chi.Router) {
		r.With(middleware.InterceptJson(map[string][]predicate.Predicate[string]{
			"email":    {predicate.IsNotEmpty, predicate.EmailIsValid},
			"password": {predicate.IsNotEmpty},
		})).Post("/", api.HandleLogin)
	})

	// User Routes
	r.Route("/users", func(r chi.Router) {
		r.Get("/", api.HandleGetAllUsers)
		r.With(middleware.InterceptJson(map[string][]predicate.Predicate[string]{
			"email": {predicate.IsNotEmpty, predicate.EmailIsValid},
		})).Post("/", api.HandleCreateUser)
		r.Route("/{user_id}", func(r chi.Router) {
			r.Use(middleware.URLParam("user_id", predicate.AllowedCharacters, predicate.NonNegative))
			r.Get("/", api.HandleGetUserById)
			r.Delete("/", api.HandleDeleteUserById)
			r.With(middleware.InterceptJson(map[string][]predicate.Predicate[string]{
				"email": {predicate.EmailIsValid},
			})).Put("/", api.HandleUpdateUserById)
			r.Route("/files", func(r chi.Router) {
				r.Get("/", api.HandleGetUserFiles)
				r.Route("/{file_name}", func(r chi.Router) {
					r.Use(middleware.URLParam("file_name", predicate.AllowedCharacters))
					r.Delete("/", api.HandleDeleteUserFileByName)
					r.Get("/", api.HandleGetUserFileByName)
				})
			})
		})
	})

	// File Routes
	r.Route("/files", func(r chi.Router) {
		r.Get("/", api.HandleGetAllFiles)
		r.Post("/", api.HandleCreateFile(constants.FileMap,
			[]predicate.Predicate[string]{
				predicate.NonNegative,
				predicate.AllowedCharacters}))
		r.Route("/{file_id}", func(r chi.Router) {
			r.Use(middleware.URLParam("file_id", predicate.AllowedCharacters, predicate.NonNegative))
			r.Get("/", api.HandleGetFileById)
			r.Put("/", api.HandleUpdateFileById)
		})
	})
	log.Println(fmt.Sprintf("Starting server on: '%s'", cfg.Address))
	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		log.Fatal(err)
	}
}

func applicationMiddleWare(cfg *config.Config) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(strings.ToUpper(cfg.ReferenceHeader))
			if apiKey != cfg.ReferenceKey {
				http.Error(w, "Forbidden", http.StatusForbidden)
				if apiKey == "" {
					return
				}
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
