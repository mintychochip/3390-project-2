package user

import (
	"database/sql"
	"time"
)

type File struct {
	ID         uint32    `json:"id"`
	UserID     uint32    `json:"user_id"`
	Name       string    `json:"name"`
	UploadTime time.Time `json:"upload_time"`
}

type FileService struct {
	*genericService[File, uint32]
}

func NewFileService(db *sql.DB) *FileService {
	return &FileService{
		&genericService[File, uint32]{
			db: db,
		},
	}
}
func (fs *FileService) GetUserFile(userId uint32, fileName string) (*File, error) {
	return fs.getItem("SELECT id, upload_time FROM user_files WHERE name = ? AND user_id = ?",
		[]interface{}{fileName, userId}, func(t *File, rows *sql.Rows) error {
			err := rows.Scan(&t.ID, &t.UploadTime)
			if err != nil {
				return err
			}
			t.Name = fileName
			t.UserID = userId
			return nil
		})
}
func (fs *FileService) GetAllUserFiles(userId uint32) ([]*File, error) {
	return fs.getAllItems("SELECT id,name,upload_time FROM user_files WHERE user_id = ?", []interface{}{userId}, func(t *File, rows *sql.Rows) error {
		err := rows.Scan(&t.ID, &t.Name, &t.UploadTime)
		if err != nil {
			return err
		}
		t.UserID = userId
		return nil
	})
}
