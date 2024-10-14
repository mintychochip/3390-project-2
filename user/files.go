package user

import (
	"database/sql"
	"time"
)

type File struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"userId"`
	FileName   string    `json:"file_name"`
	UploadTime time.Time `json:"upload_time"`
}

type FileService struct {
	DB *sql.DB
}

func (fs *FileService) AllFiles(u *User) ([]*File, error) {
	stmt, err := fs.DB.Prepare("SELECT id, user_id, file_name, upload_time FROM user_files WHERE user_id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(u.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File

	for rows.Next() {
		var File File
		if err := rows.Scan(&File.ID, &File.UserID, &File.FileName, &File.UploadTime); err != nil {
			return nil, err
		}
		files = append(files, &File)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}
