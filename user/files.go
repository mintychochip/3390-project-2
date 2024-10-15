package user

import (
	"database/sql"
	"errors"
	"time"
)

type File struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	Name       string    `json:"name"`
	UploadTime time.Time `json:"upload_time"`
}

type FileService struct {
	DB *sql.DB
}

func (fs *FileService) GetFiles(userId uint64) ([]*File, error) {
	stmt, err := fs.DB.Prepare("SELECT id, user_id, name, upload_time FROM user_files WHERE user_id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File

	for rows.Next() {
		var File File
		if err := rows.Scan(&File.ID, &File.UserID, &File.Name, &File.UploadTime); err != nil {
			return nil, err
		}
		files = append(files, &File)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}
func (fs *FileService) UserHasFile(f *File) (bool, error) {
	var exists bool
	row := fs.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user_files WHERE name = ?)", f.Name).Scan(&exists)
	if row != nil {
		return false, errors.New("user has no file")
	}
	return exists, nil
}
func (fs *FileService) AssignUserFile(f *File, u *User) (bool, error) {
	stmt, err := fs.DB.Prepare("INSERT INTO user_files (user_id,name) VALUES (?,?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.ID, f.Name)
	if err != nil {
		return false, err
	}
	return true, nil
}
