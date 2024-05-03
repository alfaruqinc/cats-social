package service

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type CatMatchService interface {
	CreateCatMatch(ctx context.Context, user *domain.User, catMatchPayload *domain.CatMatch) domain.MessageErr
	GetCatMatchesByIssuerOrReceiverID(ctx context.Context, id string) ([]domain.CatMatchResponse, domain.MessageErr)
	UpdateCatMatchByID(ctx context.Context, id string, catMatchPayload *domain.CatMatch) (string, domain.MessageErr)
	DeleteCatMatchByID(ctx context.Context, id string, userId string) domain.MessageErr
	ApproveCatMatch(ctx context.Context, userId string, matchId string) domain.MessageErr
	RejectCatMatch(ctx context.Context, userId string, matchId string) domain.MessageErr
}

type catMatchService struct {
	db                 *sql.DB
	catMatchRepository repository.CatMatchRepository
	catRepository      repository.CatRepository
}

func NewCatMatchService(db *sql.DB, catMatchRepository repository.CatMatchRepository, catRespository repository.CatRepository) CatMatchService {
	return &catMatchService{
		db:                 db,
		catMatchRepository: catMatchRepository,
		catRepository:      catRespository,
	}
}

func (c *catMatchService) CreateCatMatch(ctx context.Context, user *domain.User, catMatchPayload *domain.CatMatch) domain.MessageErr {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.NewBadRequest("Failed to start transaction")
	}
	defer tx.Rollback()

	owner, err := c.catRepository.CheckOwnerCat(ctx, tx, catMatchPayload.UserCatID, user.Id)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}
	if !owner {
		return domain.NewNotFoundError("You are not the cat's owner")
	}

	_, err = c.catMatchRepository.CreateCatMatch(ctx, tx, catMatchPayload)
	if err != nil {
		return domain.NewBadRequest("Failed to create cat match")
	}

	err = tx.Commit()
	if err != nil {
		return domain.NewBadRequest("Failed to commit transaction")
	}

	return nil
}

func (c *catMatchService) GetCatMatchesByIssuerOrReceiverID(ctx context.Context, id string) ([]domain.CatMatchResponse, domain.MessageErr) {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, domain.NewBadRequest("Failed to start transaction")
	}
	defer tx.Rollback()

	catMatches, err := c.catMatchRepository.GetCatMatchesByIssuerOrReceiverID(ctx, tx, id)
	if err != nil {
		return nil, domain.NewBadRequest("Failed to get cat match")
	}
	tx.Commit()

	var catMatchResponses []domain.CatMatchResponse
	for _, catMatch := range catMatches {
		// TODO: create helper for mapping response
		catMatchResponses = append(catMatchResponses, domain.CatMatchResponse{
			ID:        catMatch.ID,
			CreatedAt: catMatch.CreatedAt,
			Message:   catMatch.Message,
			MatchCat: domain.CatResponse{
				ID:          catMatch.MatchCat.ID,
				Name:        catMatch.MatchCat.Name,
				Race:        catMatch.MatchCat.Race,
				Sex:         catMatch.MatchCat.Sex,
				AgeInMonth:  catMatch.MatchCat.AgeInMonth,
				Description: catMatch.MatchCat.Description,
				ImageUrls:   catMatch.MatchCat.ImageUrls,
				HasMatched:  catMatch.MatchCat.HasMatched,
				CreatedAt:   catMatch.MatchCat.CreatedAt,
			},
			UserCat: domain.CatResponse{
				ID:          catMatch.UserCat.ID,
				Name:        catMatch.UserCat.Name,
				Race:        catMatch.UserCat.Race,
				Sex:         catMatch.UserCat.Sex,
				AgeInMonth:  catMatch.UserCat.AgeInMonth,
				Description: catMatch.UserCat.Description,
				ImageUrls:   catMatch.UserCat.ImageUrls,
				HasMatched:  catMatch.UserCat.HasMatched,
				CreatedAt:   catMatch.UserCat.CreatedAt,
			},
		})
	}

	return catMatchResponses, nil
}

func (c *catMatchService) UpdateCatMatchByID(ctx context.Context, id string, catMatchPayload *domain.CatMatch) (string, domain.MessageErr) {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return "", domain.NewInternalServerError("Failed to start transaction")
	}
	defer tx.Rollback()

	// TODO check if the cat match request is exist

	matchCatID, err := uuid.Parse(id)
	if err != nil {
		return "", domain.NewBadRequest("Invalid cat match id")
	}

	err = c.catMatchRepository.UpdateCatMatchByID(ctx, tx, matchCatID.String(), catMatchPayload)
	if err != nil {
		return "", domain.NewInternalServerError("Failed to update cat match")
	}

	// TODO: Delete All cat match request that has status pending

	tx.Commit()

	return "successfully matches the cat match request", nil
}

func (c *catMatchService) DeleteCatMatchByID(ctx context.Context, id string, userId string) domain.MessageErr {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.NewInternalServerError("Failed to start transaction")
	}
	defer tx.Rollback()

	canDelete, err := c.catMatchRepository.CanDeleteCatMatch(ctx, tx, id, userId)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}

	if !canDelete {
		return domain.NewNotFoundError("Cat match request is not found")
	}

	status, err := c.catMatchRepository.GetStatusCatMatchByID(ctx, tx, id)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}

	if status != "waiting" {
		return domain.NewBadRequest("Cannot delete non waiting cat match request")
	}

	err = c.catMatchRepository.DeleteCatMatchByID(ctx, tx, id)
	if err != nil {
		return domain.NewBadRequest("Failed to create cat match")
	}

	err = tx.Commit()
	if err != nil {
		return domain.NewInternalServerError("Failed to commit transaction")
	}

	return nil
}

func (c *catMatchService) ApproveCatMatch(ctx context.Context, userId string, matchId string) domain.MessageErr {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.NewInternalServerError("Failed to start transaction")
	}
	defer tx.Rollback()

	userIsReceiver, err := c.catMatchRepository.CheckIfUserIsReceiver(ctx, tx, matchId, userId)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}
	if !userIsReceiver {
		return domain.NewNotFoundError("Cat match request is not found")
	}

	status, err := c.catMatchRepository.GetStatusCatMatchByID(ctx, tx, matchId)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}
	if status != "waiting" {
		return domain.NewBadRequest("Cat match request already approved or rejected")
	}

	err = c.catMatchRepository.ApproveCatMatch(ctx, tx, userId, matchId)
	if err != nil {
		return domain.NewBadRequest("Failed to approve cat match request")
	}

	err = tx.Commit()
	if err != nil {
		return domain.NewInternalServerError("Failed to commit transaction")
	}
	return nil
}

func (c *catMatchService) RejectCatMatch(ctx context.Context, userId string, matchId string) domain.MessageErr {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.NewInternalServerError("Failed to start transaction")
	}
	defer tx.Rollback()

	userIsReceiver, err := c.catMatchRepository.CheckIfUserIsReceiver(ctx, tx, matchId, userId)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}
	if !userIsReceiver {
		return domain.NewNotFoundError("Cat match request is not found")
	}

	status, err := c.catMatchRepository.GetStatusCatMatchByID(ctx, tx, matchId)
	if err != nil {
		return domain.NewInternalServerError("something went wrong")
	}
	if status != "waiting" {
		return domain.NewBadRequest("Cat match request already approved or rejected")
	}

	err = c.catMatchRepository.RejectCatMatch(ctx, tx, userId, matchId)
	if err != nil {
		return domain.NewBadRequest("Failed to reject cat match request")
	}

	err = tx.Commit()
	if err != nil {
		return domain.NewInternalServerError("Failed to commit transaction")
	}
	return nil
}
