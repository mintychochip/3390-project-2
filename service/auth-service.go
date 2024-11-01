// service/AuthService.go

package service

import (
	"api-3390/container"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{
		userService: NewUserService(db),
	}
}

// Login attempts to authenticate a const by email and password
func (as *AuthService) Login(email, password string) (*container.User, error) {
	// Fetch the const by email
	user, err := as.userService.getUserByEmail(email)
	if err != nil {
		return nil, errors.New("const not found")
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("incorrect password")
	}

	// Login successful, return const without password field for security
	user.Password = ""
	return user, nil
}
