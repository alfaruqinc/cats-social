package repository

import (
	"cats-social/internal/domain"
	"database/sql"
)

type CatRepository interface {
	CreateCat(cat *domain.Cat) error
	GetAllCats() ([]domain.Cat, error)
}

type CatRepositoryImpl struct {
}

func NewCatRepository(db *sql.DB) CatRepository {
	return &CatRepositoryImpl{}
}

func (c *CatRepositoryImpl) CreateCat(cat *domain.Cat) error {
	return nil
}

func (c *CatRepositoryImpl) GetAllCats() ([]domain.Cat, error) {
	return nil, nil
}
