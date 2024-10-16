package service

import (
	"api-3390/container"
	"database/sql"
)

type FileService struct {
	*genericService[container.File, uint32]
}

func NewFileService(db *sql.DB) *FileService {
	return &FileService{
		&genericService[container.File, uint32]{
			db: db,
		},
	}
}

func (fs *FileService) UpdateUserById(f *container.File) error {
	return fs.updateItem("UPDATE user_files SET user_id = ?, name = ?, upload_time = ? WHERE id = ?",
		[]interface{}{f.UserID, f.Name, f.UploadTime, f.ID})
}

//	func (fs *FileService) DeleteUserById(k uint32) error {
//		return us.deleteItems("DELETE FROM users WHERE id = ?", []interface{}{k})
//	}
//
//	func (fs *FileService) GetUserById(k uint32) (*container.User, error) {
//		return us.getItem("SELECT name,email,password FROM users WHERE id = ?", []interface{}{k},
//			func(t *container.User, rows *sql.Rows) error {
//				return rows.Scan(&t.Name, &t.Email, &t.Password)
//			})
//	}
func (fs *FileService) GetAllFiles() ([]*container.File, error) {
	return fs.getAllItems("SELECT * FROM user_files", []interface{}{}, func(t *container.File, rows *sql.Rows) error {
		return rows.Scan(&t.ID, &t.UserID, &t.Name, &t.UploadTime)
	})
}

//func (fs *FileService) CreateUser(u *container.User) error {
//	if u.Email == "" || u.Name == "" || u.Password == "" {
//		return errors.New("fields were not completed")
//	}
//	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//	return us.insertItem(u, func(u *container.User) (string, []interface{}) {
//		return "INSERT INTO users (name,email,password) VALUES (?,?,?)", []interface{}{u.Name, u.Email, string(hashed)}
//	})
//}
