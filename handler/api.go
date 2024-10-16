package handler

import "api-3390/service"

type API struct {
	Services *Services
}
type Services struct {
	UserService *service.UserService
	FileService *service.FileService
}

func NewServices(us *service.UserService, fs *service.FileService) *Services {
	return &Services{
		UserService: us,
		FileService: fs,
	}
}
