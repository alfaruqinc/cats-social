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

type CatQueryParams struct {
	Id               string `form:"id"`
	Limit            int32  `form:"limit"`
	Offset           int32  `form:"offset"`
	Race             string `form:"race"`
	Sex              string `form:"sex"`
	IsAlreadyMatched bool   `form:"isAlreadyMatched"`
	AgeInMonth       int32  `form:"ageInMonth"`
	Owned            bool   `form:"owned"`
	Search           string `form:"search"`
}
