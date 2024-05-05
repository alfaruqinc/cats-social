package service

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
	"database/sql"
	"net/url"
)

type CatService interface {
	CreateCat(cat *domain.Cat) error
	GetAllCats(user *domain.User, queryParams url.Values) ([]domain.Cat, domain.MessageErr)
}

type catService struct {
	db            *sql.DB
	catRepository repository.CatRepository
}

func NewCatService(db *sql.DB, catRepository repository.CatRepository) CatService {
	return &catService{
		db:            db,
		catRepository: catRepository,
	}
}

func (c *catService) CreateCat(cat *domain.Cat) error {
	return c.catRepository.CreateCat(c.db, cat)
}

func (c *catService) GetAllCats(user *domain.User, queryParams url.Values) ([]domain.Cat, domain.MessageErr) {
	cats, err := c.catRepository.GetAllCats(c.db, user, queryParams)
	if err != nil {
		return nil, domain.NewInternalServerError("something went wrong")
	}

	return cats, nil
}
