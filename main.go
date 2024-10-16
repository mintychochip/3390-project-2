package main

import (
	"api-3390/config"
	"api-3390/container/predicate"
	"api-3390/handler"
	"api-3390/handler/middleware"
	"api-3390/service"
	"api-3390/user"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
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

	_, err = db.Exec(user.UserTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(user.UserFileTable)
	if err != nil {
		log.Fatal(err)
	}
	services := handler.NewServices(service.NewUserService(db), service.NewFileService(db))
	api := handler.API{Services: services}
	r := chi.NewRouter()
	r.Use(applicationMiddleWare(cfg))
	r.Route("/users", func(r chi.Router) {
		r.Get("/", api.HandleGetAllUsers)
		r.With(middleware.InterceptJson(map[string][]predicate.Predicate[string]{
			"email": {predicate.IsNotEmpty, predicate.EmailIsValid},
		})).Post("/", api.HandleCreateUser)
		r.Route("/{user_id}", func(r chi.Router) {
			r.Use(middleware.URLParam("user_id", predicate.AllowedCharacters, predicate.NonNegativePredicate))
			r.Get("/", api.HandleGetUserById)
			r.Delete("/", api.HandleDeleteUserById)
			r.With(middleware.InterceptJson(map[string][]predicate.Predicate[string]{
				"email": {predicate.EmailIsValid},
			})).Put("/", api.HandleUpdateUserById)
		})
	})
	r.Route("/files", func(r chi.Router) {
		r.Get("/", api.HandleGetAllFiles)
	})
	log.Println(fmt.Sprintf("Starting server on: '%s'", cfg.Address()))
	if err := http.ListenAndServe(cfg.Address(), r); err != nil {
		log.Fatal(err)
	}
}

func applicationMiddleWare(cfg *config.Config) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-KEY")
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
