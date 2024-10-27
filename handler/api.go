package handler

import "api-3390/service"

type API struct {
	Services *Services
}
type Services struct {
	AuthService *service.AuthService
	FileService *service.FileService
	UserService *service.UserService
}

func NewServices(as *service.AuthService, fs *service.FileService, us *service.UserService) *Services {
	return &Services{
		AuthService: as,
		FileService: fs,
		UserService: us,
	}
}
