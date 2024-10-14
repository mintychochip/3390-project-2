package user

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

type AuthenticationService struct {
	DB          *sql.DB
	SigningKey  []byte
	TokenExpire time.Duration
}

func AuthService(db *sql.DB, duration time.Duration) (*AuthenticationService, error) {
	key, err := generateSignedKey()
	if err != nil {
		return nil, err
	}
	auth := AuthenticationService{
		DB:          db,
		TokenExpire: duration,
		SigningKey:  key,
	}
	return &auth, nil
}
func (auth *AuthenticationService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return auth.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func generateSignedKey() ([]byte, error) {
	key := make([]byte, SigningKeyLength)
	if _, err := rand.Read(key); err != nil {
		return []byte{}, err
	}
	return key, nil
}

func (auth *AuthenticationService) RegisterUser(u *User) error {
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
func (auth *AuthenticationService) ExistsUser(u *User) (bool, error) {
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

func (auth *AuthenticationService) AuthenticateUser(u *User) (string, error) {
	if b, err := compareHashedPassword(auth.DB, u); err != nil {
		return "", err
	} else if !b {
		return "", errors.New("invalid credentials")
	}
	token, err := auth.generateToken(u)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (auth *AuthenticationService) generateToken(u *User) (string, error) {
	claims := jwt.MapClaims{
		"id":    u.Name,
		"email": u.Email,
		"exp":   time.Now().Add(auth.TokenExpire).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(auth.SigningKey)
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
