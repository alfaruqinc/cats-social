package repository

import (
	"cats-social/internal/domain"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateNewUser(db *sql.DB, userPayload *domain.User) error
	GetById(db *sql.DB, userId uuid.UUID) error
	GetByEmail(db *sql.DB, userEmail string) error
}

type userImpl struct {
}

func NewUserPg() UserRepository {
	return &userImpl{}
}

func (u *userImpl) CreateNewUser(db *sql.DB, userPayload *domain.User) error {
	query := `INSERT INTO users (id, name, email, password)
		VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(query, userPayload.Id, userPayload.Name, userPayload.Email, userPayload.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *userImpl) GetById(db *sql.DB, userId uuid.UUID) error {
	query := `SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE id = $1
		)
	`
	var idExists bool
	err := db.QueryRow(query, userId).Scan(&idExists)
	if err != nil {
		return err
	}

	if !idExists {
		return errors.New("userId is not found")
	}

	return nil
}

func (u *userImpl) GetByEmail(db *sql.DB, userEmail string) error {
	query := `SELECT EXISTS(
		SELECT 1 
		FROM users 
		WHERE email = $1
		)
	`

	var exists bool
	err := db.QueryRow(query, userEmail).Scan(&exists, &domain.NewUser().Password)
	if err != nil {
		return err
	}

	if !exists {

		return errors.New("user not found")
	}

	return nil
}
