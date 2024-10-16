package handler

import (
	"api-3390/container"
	"encoding/json"
	"errors"
	_ "github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
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
	id, err := getUserId(r)
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

	err = a.Services.UserService.UpdateUserById(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (a *API) HandleDeleteUserById(w http.ResponseWriter, r *http.Request) {
	id, err := getUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := a.Services.UserService.DeleteUserById(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (a *API) HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := getUserId(r)
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
func getUserId(r *http.Request) (uint32, error) {
	val, ok := r.Context().Value("user_id").(string)
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