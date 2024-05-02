package service

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
	"context"
	"database/sql"
)

type CatMatchService interface {
	CreateCatMatch(ctx context.Context, catMatchPayload *domain.CatMatch) (string, domain.MessageErr)
	GetCatMatches(ctx context.Context) ([]domain.CatMatchResponse, domain.MessageErr)
	GetCatMatchByID(ctx context.Context, id string) (*domain.CatMatchResponse, domain.MessageErr)
	UpdateCatMatchByID(ctx context.Context, id string, catMatchPayload *domain.CatMatch) (*domain.CatMatchResponse, domain.MessageErr)
	DeleteCatMatch(ctx context.Context, id string) domain.MessageErr
}

type catMatchService struct {
	catMatchRepository repository.CatMatchRepository
	db                 *sql.DB
}

func NewCatMatchService(catMatchRepository repository.CatMatchRepository, db *sql.DB) CatMatchService {
	return &catMatchService{
		catMatchRepository: catMatchRepository,
		db:                 db,
	}
}

func (c *catMatchService) CreateCatMatch(ctx context.Context, catMatchPayload *domain.CatMatch) (string, domain.MessageErr) {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return "", domain.NewBadRequest("Failed to start transaction")
	}
	defer tx.Rollback()

	_, err = c.catMatchRepository.CreateCatMatch(ctx, tx, catMatchPayload)
	if err != nil {
		return "", domain.NewBadRequest("Failed to create cat match")
	}

	err = tx.Commit()
	if err != nil {
		return "", domain.NewBadRequest("Failed to commit transaction")
	}

	return "successfully send match request", nil
}

func (c *catMatchService) GetCatMatches(ctx context.Context) ([]domain.CatMatchResponse, domain.MessageErr) {
	return nil, nil
}

func (c *catMatchService) GetCatMatchByID(ctx context.Context, id string) (*domain.CatMatchResponse, domain.MessageErr) {
	return nil, nil
}

func (c *catMatchService) UpdateCatMatchByID(ctx context.Context, id string, catMatchPayload *domain.CatMatch) (*domain.CatMatchResponse, domain.MessageErr) {
	return nil, nil
}

func (c *catMatchService) DeleteCatMatch(ctx context.Context, id string) domain.MessageErr {
	return nil
}
