package service

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
	"database/sql"
	"net/url"
)

type CatService interface {
	CreateCat(cat *domain.Cat) domain.MessageErr
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

func (c *catService) CreateCat(cat *domain.Cat) domain.MessageErr {
	err := c.catRepository.CreateCat(c.db, cat)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}

	return nil
}

func (c *catService) GetAllCats(user *domain.User, queryParams url.Values) ([]domain.Cat, domain.MessageErr) {
	cats, err := c.catRepository.GetAllCats(c.db, user, queryParams)
	if err != nil {
		return nil, domain.NewInternalServerError("something went wrong")
	}

	return cats, nil
}
