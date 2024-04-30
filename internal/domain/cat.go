package domain

import (
	"time"

	"github.com/google/uuid"
)

type CreateCatRequest struct {
	Name string `json:"name" validate:"required"`
	Race string `json:"race" validate:"required"`
}

type Cat struct {
	ID          uuid.UUID `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Name        string    `db:"name"`
	Race        string    `db:"race"`
	Sex         string    `db:"sex"`
	AgeInMonth  int32     `db:"age_in_month"`
	Description string    `db:"description"`
	ImageUrls   []string  `db:"image_urls"`
	HasMatched  bool      `db:"has_matched"`
	OwnedBy     uuid.UUID `db:"owned_by"`
}

type CatResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
}

func NewCat() *Cat {
	id := uuid.New()
	createdAt := time.Now().Format(time.RFC3339)
	parsedCreatedAt, _ := time.Parse(time.RFC3339, createdAt)

	return &Cat{
		ID:         id,
		CreatedAt:  parsedCreatedAt,
		HasMatched: false,
	}
}

var CatRace = []string{
	"Persian",
	"Maine Coon",
	"Siamese",
	"Ragdoll",
	"Bengal",
	"Sphynx",
	"British Shorthair",
	"Abyssinian",
	"Scottish Fold",
	"Birman",
}

var CatSex = []string{"male", "female"}

var CatQueryParams = []string{
	"id",
	"limit",
	"offset",
	"race",
	"sex",
	"hasMatched",
	"ageInMonth",
	"owned",
	"search",
}
