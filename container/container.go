package container

import "time"

type User struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type File struct {
	ID         uint32    `json:"id"`
	UserID     uint32    `json:"user_id"`
	Name       string    `json:"name"`
	UploadTime time.Time `json:"upload_time"`
}
