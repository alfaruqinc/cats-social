package service

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
)

type CatMatchService interface {
	CreateCatMatch(catMatchPayload *domain.CatMatch) (*domain.CatMatchResponse, domain.MessageErr)
	GetCatMatches() ([]domain.CatMatchResponse, domain.MessageErr)
	GetCatMatchByID(id string) (*domain.CatMatchResponse, domain.MessageErr)
	UpdateCatMatchByID(id string, catMatchPayload *domain.CatMatch) (*domain.CatMatchResponse, domain.MessageErr)
	DeleteCatMatch(id string) domain.MessageErr
}

type catMatchService struct {
	catMatchRepository repository.CatMatchRepository
}

func NewCatMatchService(catMatchRepository repository.CatMatchRepository) CatMatchService {
	return &catMatchService{
		catMatchRepository: catMatchRepository,
	}
}

func (c *catMatchService) CreateCatMatch(catMatchPayload *domain.CatMatch) (*domain.CatMatchResponse, domain.MessageErr) {
	return nil, nil
}

func (c *catMatchService) GetCatMatches() ([]domain.CatMatchResponse, domain.MessageErr) {
	return nil, nil
}

func (c *catMatchService) GetCatMatchByID(id string) (*domain.CatMatchResponse, domain.MessageErr) {
	return nil, nil
}

func (c *catMatchService) UpdateCatMatchByID(id string, catMatchPayload *domain.CatMatch) (*domain.CatMatchResponse, domain.MessageErr) {
	return nil, nil
}

func (c *catMatchService) DeleteCatMatch(id string) domain.MessageErr {
	return nil
}
