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

func Register(db *sql.DB, user User) error {
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

/*
*
Returns nil if the user does exist
*/
func UserExists(db *sql.DB, u *User) (bool, error) {
	var exists bool
	if u.Name == "" && u.Email == "" {
		return false, errors.New("username and email were not provided")
	}

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = ? OR email = ?)", u.Name, u.Email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %s", err)
	}

	return exists, nil
}
func UserAuthenticate(db *sql.DB, u *User) (bool, error) {
	b, err := compareHashedPassword(db, u)
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
