package domain

import (
	"time"

	"github.com/google/uuid"
)

var (
	CatMatchStatusPending  = "pending"
	CatMatchStatusAccepted = "accepted"
	CatMatchStatusRejected = "rejected"
)

var CatMatchStatuses = []string{
	CatMatchStatusPending,
	CatMatchStatusAccepted,
	CatMatchStatusRejected,
}

type CreateCatMatchRequest struct {
	UserCatID  string `json:"userCatId" validate:"required"`
	MatchCatID string `json:"matchCatId" validate:"required"`
	Message    string `json:"message" validate:"required"`
}

func NewCatMatchFromBody(body CreateCatMatchRequest) *CatMatch {
	id := uuid.New()
	createdAt := time.Now().Format(time.RFC3339)
	parsedCreatedAt, _ := time.Parse(time.RFC3339, createdAt)

	parsedMatchCatId, _ := uuid.Parse(body.MatchCatID)
	parsedUserCatId, _ := uuid.Parse(body.UserCatID)

	return &CatMatch{
		ID:         id,
		CreatedAt:  parsedCreatedAt,
		MatchCatID: parsedMatchCatId,
		UserCatID:  parsedUserCatId,
		Message:    body.Message,
	}
}

type CatMatch struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	IssuedByID uuid.UUID `db:"issued_by_id"`
	IssuedBy   User
	MatchCatID uuid.UUID `db:"match_cat_id"`
	MatchCat   Cat
	UserCatID  uuid.UUID `db:"user_cat_id"`
	UserCat    Cat
	Message    string
	Status     string
}

type CatMatchResponse struct {
	ID        uuid.UUID    `json:"id"`
	IssuedBy  UserResponse `json:"issuedBy"`
	MatchCat  CatResponse  `json:"matchCatDetail"`
	UserCat   CatResponse  `json:"userCatDetail"`
	Message   string       `json:"message"`
	CreatedAt time.Time    `json:"createdAt"`
}

func NewCatMatch() *CatMatch {
	id := uuid.New()
	createdAt := time.Now().Format(time.RFC3339)
	parsedCreatedAt, _ := time.Parse(time.RFC3339, createdAt)

	return &CatMatch{
		ID:        id,
		CreatedAt: parsedCreatedAt,
	}
}
