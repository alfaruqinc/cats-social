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
	UpdateCat(user *domain.User, cat *domain.Cat) domain.MessageErr
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

func (c *catService) UpdateCat(user *domain.User, cat *domain.Cat) domain.MessageErr {
	catExists, err := c.catRepository.CheckCatExists(c.db, cat.ID, user.Id)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}
	if catExists {
		return domain.NewNotFoundError("cat does not exists")
	}

	canEdit, err := c.catRepository.CheckEditableSexV2(c.db, cat)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}
	if canEdit {
		return domain.NewBadRequest("cannot edit sex when already requested to match")
	}

	err = c.catRepository.UpdateCat(c.db, cat)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}

	return nil
}
