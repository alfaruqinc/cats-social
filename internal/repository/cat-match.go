package repository

import (
	"cats-social/internal/domain"
	"context"
	"database/sql"
)

type CatMatchRepository interface {
	CreateCatMatch(ctx context.Context, tx *sql.DB, catMatch *domain.CatMatch) (*domain.CatMatch, error)
	GetCatMatches(ctx context.Context, tx *sql.DB) ([]domain.CatMatch, error)
	GetCatMatchByID(ctx context.Context, tx *sql.DB, id string) (*domain.CatMatch, error)
	UpdateCatMatchByID(ctx context.Context, tx *sql.DB, id string, catMatch *domain.CatMatch) error
	DeleteCatMatch(id string) error
}

type catMatchRepository struct {
}

func NewCatMatchRepository() CatMatchRepository {
	return &catMatchRepository{}
}

func (c *catMatchRepository) CreateCatMatch(ctx context.Context, tx *sql.DB, catMatch *domain.CatMatch) (*domain.CatMatch, error) {
	return nil, nil
}

func (c *catMatchRepository) GetCatMatches(ctx context.Context, tx *sql.DB) ([]domain.CatMatch, error) {
	return nil, nil
}

func (c *catMatchRepository) GetCatMatchByID(ctx context.Context, tx *sql.DB, id string) (*domain.CatMatch, error) {
	return nil, nil
}

func (c *catMatchRepository) UpdateCatMatchByID(ctx context.Context, tx *sql.DB, id string, catMatch *domain.CatMatch) error {
	return nil
}

func (c *catMatchRepository) DeleteCatMatch(id string) error {
	return nil
}
