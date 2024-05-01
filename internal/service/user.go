package service

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
)

type UserService interface {
	Register(userPayload *domain.User) (*domain.UserResponse, domain.MessageErr)
	Login(userPayload *domain.User) (*domain.UserResponse, domain.MessageErr)
}

type userService struct {
	ur repository.UserRepository
}

func NewUserService(ur repository.UserRepository) userService {
	return userService{ur: ur}
}

func (us *userService) Register(userPayload *domain.User) (*domain.UserResponse, domain.MessageErr) {
	return us.ur.CreateNewUser(userPayload)
}
func (us *userService) Login(userPayload *domain.User) (*domain.UserResponse, domain.MessageErr) {
	return us.ur.CreateNewUser(userPayload)
}
