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
func (fs *FileService) GetUserFiles(userId uint32) ([]*container.File, error) {
	return fs.getAllItems("SELECT id,name,upload_time FROM user_files WHERE user_id = ?", []interface{}{userId}, func(t *container.File, rows *sql.Rows) error {
		t.UserID = userId
		return rows.Scan(&t.ID, &t.Name, &t.UploadTime)
	})
}
func (fs *FileService) UpdateFile(f *container.File) error {
	return fs.updateItem("UPDATE user_files SET user_id = ?, name = ? WHERE id = ?",
		[]interface{}{f.UserID, f.Name, f.ID})
}

func (fs *FileService) DeleteFileById(k uint32) error {
	return fs.deleteItems("DELETE FROM user_files WHERE id = ?", []interface{}{k})
}
func (fs *FileService) GetFileById(k uint32) (*container.File, error) {
	return fs.getItem("SELECT user_id,name,upload_time FROM user_files WHERE id = ?", []interface{}{k},
		func(f *container.File, rows *sql.Rows) error {
			f.ID = k
			return rows.Scan(&f.UserID, &f.Name, &f.UploadTime)
		})
}
func (fs *FileService) GetAllFiles() ([]*container.File, error) {
	return fs.getAllItems("SELECT * FROM user_files", []interface{}{}, func(t *container.File, rows *sql.Rows) error {
		return rows.Scan(&t.ID, &t.UserID, &t.Name, &t.UploadTime)
	})
}

func (fs *FileService) CreateFile(f *container.File) error {
	return fs.insertItem("INSERT INTO user_files (user_id, name) VALUES (?,?)",
		[]interface{}{f.UserID, f.Name})
}
