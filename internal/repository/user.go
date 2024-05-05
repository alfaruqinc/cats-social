package repository

import (
	"cats-social/internal/domain"
	"database/sql"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateNewUser(db *sql.DB, userPayload *domain.User) error
	GetById(db *sql.DB, userId uuid.UUID) (*domain.User, error)
	GetByEmail(db *sql.DB, userEmail string) (*domain.User, error)
}

type userRepository struct{}

func NewUserPg() UserRepository {
	return &userRepository{}
}

func (u *userRepository) CreateNewUser(db *sql.DB, userPayload *domain.User) error {
	query := `INSERT INTO users (id, name, email, password)
		VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(query, userPayload.Id, userPayload.Name, userPayload.Email, userPayload.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) GetById(db *sql.DB, userId uuid.UUID) (*domain.User, error) {
	query := `SELECT id, email, name, password
		FROM users WHERE id = $1
	`
	user := domain.User{}

	err := db.QueryRow(query, userId).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) GetByEmail(db *sql.DB, userEmail string) (*domain.User, error) {
	query := `SELECT id, email, name, password
		FROM users WHERE email = $1
	`
	user := domain.User{}

	err := db.QueryRow(query, userEmail).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
