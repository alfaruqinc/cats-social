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

var (
	CatMatchStatuses = []string{
		CatMatchStatusPending,
		CatMatchStatusAccepted,
		CatMatchStatusRejected,
	}
)

type CreateCatMatchRequest struct {
	UserCatID  uuid.UUID `json:"userCatId" validate:"required"`
	MatchCatID uuid.UUID `json:"matchCatId" validate:"required"`
	Message    string    `json:"message" validate:"required"`
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
	ID        uuid.UUID   `json:"id"`
	MatchCat  CatResponse `json:"matchCatDetail"`
	UserCat   CatResponse `json:"userCatDetail"`
	Message   string      `json:"message"`
	CreatedAt time.Time   `json:"createdAt"`
}
