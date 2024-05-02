package service

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
	"database/sql"
)

type UserService interface {
	Register(db *sql.DB, user *domain.User) error
	Login(db *sql.DB, user *domain.User) error
}

type userService struct {
	db *sql.DB
	ur repository.UserRepository
}

func NewUserService(ur repository.UserRepository, db *sql.DB) userService {
	return userService{
		ur: ur,
		db: db,
	}
}

func (us *userService) Register(db *sql.DB, user *domain.User) error {
	return us.ur.CreateNewUser(db, user)
}
func (us *userService) Login(db *sql.DB, user *domain.User) error {
	return us.ur.CreateNewUser(db, user)
}
