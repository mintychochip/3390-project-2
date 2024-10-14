package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserRecord struct {
	ID       uint64 `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

func RegisterUser(db *sql.DB, user User) error {
	b := validEmail(user.Email)
	if !b {
		return fmt.Errorf("invalid email was attempted to be registered: '%s'", user.Email)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUserRecord := UserRecord{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPassword),
	}
	_, err = db.Exec("INSERT INTO users(name, email, password) VALUES ($1, $2, $3)", newUserRecord.Name, newUserRecord.Email, newUserRecord.Password)
	if err != nil {
		return err
	}
	return nil
}

func validEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}
func compareHashedPassword(db *sql.DB, field, identifier, password string) (bool, error) {
	var hashed string
	query := fmt.Sprintf("SELECT password FROM users WHERE %s = ?", field)

	err := db.QueryRow(query, identifier).Scan(&hashed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to locate user: %s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
func AuthenticateUserByName(db *sql.DB, name, password string) (bool, error) {
	b, err := compareHashedPassword(db, "name", name, password)
	if err != nil {
		return false, err
	}
	return b, nil
}
func AuthenticateUserByEmail(db *sql.DB, email, password string) (bool, error) {
	b, err := compareHashedPassword(db, "email", email, password)
	if err != nil {
		return false, err
	}
	return b, nil
}
