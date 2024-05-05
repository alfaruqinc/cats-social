package service

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
	"database/sql"
)

type CatService interface {
	CreateCat(db *sql.DB, cat *domain.Cat) error
	GetAllCats(db *sql.DB) ([]domain.Cat, error)
}

type catService struct {
	db            *sql.DB
	catRepository repository.CatRepository
}

func NewCatService(db *sql.DB, catRepository repository.CatRepository) CatService {
	return &catService{
		catRepository: catRepository,
		db:            db,
	}
}

func (c *catService) CreateCat(db *sql.DB, cat *domain.Cat) error {
	return c.catRepository.CreateCat(db, cat)
}

func (c *catService) GetAllCats(db *sql.DB) ([]domain.Cat, error) {
	return c.catRepository.GetAllCats(db)
}
