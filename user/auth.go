package user

// import (
//
//	"api-3390/service"
//	"crypto/rand"
//	"database/sql"
//	"errors"
//	"fmt"
//	"github.com/golang-jwt/jwt/v4"
//	"golang.org/x/crypto/bcrypt"
//	"time"
//
// )
//
//	type AuthenticationService struct {
//		*service.genericService[User, uint32]
//		SigningKey  []byte
//		TokenExpire time.Duration
//	}
//
//	func (a *AuthenticationService) UserIsRegistered(obj *User) (bool, error) {
//		if obj.Name == "" && obj.Email == "" {
//			return false, errors.New("username and email were not provided")
//		}
//		return a.itemExists(obj, func(obj *User) (string, []interface{}) {
//			return "SELECT EXISTS(SELECT 1 FROM users WHERE name = ? OR email = ?)", []interface{}{obj.Name, obj.Email}
//		})
//	}
//
//	func NewAuthenticationService(db *sql.DB, duration time.Duration) (*AuthenticationService, error) {
//		key, err := generateSignedKey()
//		if err != nil {
//			return nil, err
//		}
//		return &AuthenticationService{
//			genericService: &service.genericService[User, uint32]{
//				db: db,
//			},
//			TokenExpire: duration,
//			SigningKey:  key,
//		}, nil
//	}
//
//	func (a *AuthenticationService) ValidateToken(tokenString string) (*jwt.Token, error) {
//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, errors.New("unexpected signing method")
//			}
//			return a.SigningKey, nil
//		})
//
//		if err != nil {
//			return nil, err // You could also check for specific error types here
//		}
//
//		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//			// Check claims like exp, nbf, etc.
//			if exp, ok := claims["exp"].(float64); ok {
//				if time.Unix(int64(exp), 0).Before(time.Now()) {
//					return nil, errors.New("token is expired")
//				}
//			}
//			// Additional claims checks can be added here
//		} else {
//			return nil, errors.New("invalid token claims")
//		}
//
//		return token, nil
//	}
//
//	func generateSignedKey() ([]byte, error) {
//		key := make([]byte, SigningKeyLength)
//		if _, err := rand.Read(key); err != nil {
//			return []byte{}, err
//		}
//		return key, nil
//	}
//
//	func (a *AuthenticationService) AuthenticateUser(u *User) (bool, error) {
//		b, err := compareHashedPassword(a.db, u)
//		if err != nil {
//			return false, err
//		}
//		return b, nil
//	}
//
//	func (a *AuthenticationService) GenerateToken(u *User) (string, error) {
//		claims := jwt.MapClaims{
//			"id":    u.ID,
//			"email": u.Email,
//			"exp":   time.Now().Add(a.TokenExpire).Unix(),
//		}
//		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//		return token.SignedString(a.SigningKey)
//	}
//func compareHashedPassword(db *sql.DB, u *User) (bool, error) {
//	var hashed string
//	err := db.QueryRow("SELECT password FROM users WHERE name = ? OR email = ?", u.Name, u.Email).Scan(&hashed)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return false, nil
//		}
//		return false, fmt.Errorf("failed to locate u: %s", err)
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(u.Password))
//	if err != nil {
//		return false, nil
//	}
//	return true, nil
//}
