package models

import (
	"time"

	"github.com/google/uuid"
)

type Cat struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	Name        string    `json:"name" db:"name"`
	Race        string    `json:"race" db:"race"`
	Sex         string    `json:"sex" db:"sex"`
	AgeInMonth  int32     `json:"ageInMonth" db:"age_in_month"`
	Description string    `json:"description" db:"description"`
	ImageUrls   []string  `json:"imageUrls" db:"image_urls"`
	HasMatched  bool      `json:"hasMatched" db:"has_matched"`
	OwnedBy     uuid.UUID `json:"-" db:"owned_by"`
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
	"isAlreadyMatched",
	"ageInMonth",
	"owned",
	"search",
}
