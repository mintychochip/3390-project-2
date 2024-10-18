package user

//
//import (
//	"api-3390/service"
//	"database/sql"
//	"time"
//)
//

//
//type FileService struct {
//	*service.genericService[File, uint32]
//}
//
//func NewFileService(db *sql.DB) *FileService {
//	return &FileService{
//		&service.genericService[File, uint32]{
//			db: db,
//		},
//	}
//}
//func (fs *FileService) GetUserFile(userId uint32, fileName string) (*File, error) {
//	return fs.getItem("SELECT id, upload_time FROM user_files WHERE name = ? AND user_id = ?",
//		[]interface{}{fileName, userId}, func(t *File, rows *sql.Rows) error {
//			err := rows.Scan(&t.ID, &t.UploadTime)
//			if err != nil {
//				return err
//			}
//			t.Name = fileName
//			t.UserID = userId
//			return nil
//		})
//}
//func (fs *FileService) GetAllUserFiles(userId uint32) ([]*File, error) {
//	return fs.getAllItems("SELECT id,name,upload_time FROM user_files WHERE user_id = ?", []interface{}{userId}, func(t *File, rows *sql.Rows) error {
//		err := rows.Scan(&t.ID, &t.Name, &t.UploadTime)
//		if err != nil {
//			return err
//		}
//		t.UserID = userId
//		return nil
//	})
//}
