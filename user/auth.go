package user

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type AuthService struct {
	DB *sql.DB
}

func (auth *AuthService) RegisterUser(u *User) error {
	b := validEmail(u.Email)
	if !b {
		return fmt.Errorf("invalid email was attempted to be registered: '%s'", u.Email)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	stmt, err := auth.DB.Prepare("INSERT INTO users(name, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Name, u.Email, hashedPassword)
	return err
}
func (auth *AuthService) ExistsUser(u *User) (bool, error) {
	var exists bool
	if u.Name == "" && u.Email == "" {
		return false, errors.New("username and email were not provided")
	}

	err := auth.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = ? OR email = ?)", u.Name, u.Email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %s", err)
	}

	return exists, nil
}

func (auth *AuthService) AuthenticateUser(u *User) (bool, error) {
	b, err := compareHashedPassword(auth.DB, u)
	if err != nil {
		return false, err
	}
	return b, nil
}

func validEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

func compareHashedPassword(db *sql.DB, u *User) (bool, error) {
	var hashed string
	err := db.QueryRow("SELECT password FROM users WHERE name = ? OR email = ?", u.Name, u.Email).Scan(&hashed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to locate u: %s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(u.Password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
