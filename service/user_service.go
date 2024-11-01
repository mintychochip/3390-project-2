package service

import (
	"api-3390/container"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*genericService[container.User, uint32]
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		&genericService[container.User, uint32]{
			db: db,
		},
	}
}
func (us *UserService) UserExists(u *container.User) (bool, error) {
	return us.itemExists("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", []interface{}{u.ID})
}
func (us *UserService) UpdateUser(u *container.User) error {
	return us.updateItem("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?",
		[]interface{}{u.Name, u.Email, u.Password, u.ID})
}
func (us *UserService) DeleteUserById(k uint32) error {
	return us.deleteItems("DELETE FROM users WHERE id = ?", []interface{}{k})
}
func (us *UserService) getUserByEmail(email string) (*container.User, error) {
	return us.getItem("SELECT id, name, password FROM users WHERE email = ?", []interface{}{email},
		func(t *container.User, rows *sql.Rows) error {
			t.Email = email
			return rows.Scan(&t.ID, &t.Name, &t.Password)
		})
}
func (us *UserService) GetUserById(k uint32) (*container.User, error) {
	return us.getItem("SELECT name,email,password FROM users WHERE id = ?", []interface{}{k},
		func(t *container.User, rows *sql.Rows) error {
			return rows.Scan(&t.Name, &t.Email, &t.Password)
		})
}
func (us *UserService) GetAllUsers() ([]*container.User, error) {
	return us.getAllItems("SELECT * FROM users", []interface{}{}, func(t *container.User, rows *sql.Rows) error {
		return rows.Scan(&t.ID, &t.Name, &t.Email, &t.Password)
	})
}
func (us *UserService) CreateUser(u *container.User) error {
	if u.Email == "" || u.Name == "" || u.Password == "" {
		return errors.New("fields were not completed")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return us.insertItem("INSERT INTO users (name,email,password) VALUES (?,?,?)",
		[]interface{}{u.Name, u.Email, string(hashed)})
}
