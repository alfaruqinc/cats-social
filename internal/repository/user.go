package repository

import (
	"cats-social/internal/domain"
	"database/sql"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateNewUser(userPayload *domain.User) (*domain.UserResponse, domain.MessageErr)
	GetById(userId uuid.UUID) (*domain.User, domain.MessageErr)
	GetByEmail(userEmail string) (*domain.User, domain.MessageErr)
}

type userImpl struct {
	db *sql.DB
}

func NewUserPg(db *sql.DB) UserRepository {
	return &userImpl{}
}

func (u *userImpl) CreateNewUser(userPayload *domain.User) (*domain.UserResponse, domain.MessageErr) {
	return nil, nil
}

func (u *userImpl) GetById(userId uuid.UUID) (*domain.User, domain.MessageErr) {
	return nil, nil
}

func (u *userImpl) GetByEmail(userEmail string) (*domain.User, domain.MessageErr) {

	return nil, nil
}
