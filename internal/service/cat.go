package service

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
	"database/sql"
)

type CatService interface {
	CreateCat(cat *domain.Cat) error
	GetAllCats() ([]domain.Cat, error)
}

type catService struct {
	db            *sql.DB
	catRepository repository.CatRepository
}

func NewCatService(catRepository repository.CatRepository, db *sql.DB) CatService {
	return &catService{
		catRepository: catRepository,
		db:            db,
	}
}

func (c *catService) CreateCat(cat *domain.Cat) error {
	return c.catRepository.CreateCat(cat)
}

func (c *catService) GetAllCats() ([]domain.Cat, error) {
	return c.catRepository.GetAllCats()
}
