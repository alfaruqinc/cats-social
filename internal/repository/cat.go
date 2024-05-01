package repository

import (
	"cats-social/internal/domain"
	"database/sql"
)

type CatRepository interface {
	CreateCat(db *sql.DB, cat *domain.Cat) error
	GetAllCats(db *sql.DB) ([]domain.Cat, error)
}

type CatRepositoryImpl struct{}

func NewCatRepository() CatRepository {
	return &CatRepositoryImpl{}
}

func (c *CatRepositoryImpl) CreateCat(db *sql.DB, catBody *domain.Cat) error {
	query := `INSERT INTO cats (id, created_at, name, race, sex, age_in_month, description, image_urls, owned_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
	_, err := db.Exec(query, catBody.ID, catBody.CreatedAt, catBody.Name, catBody.Race, catBody.Sex, catBody.AgeInMonth, catBody.Description, catBody.ImageUrls, catBody.OwnedBy)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	return nil
}

func (c *CatRepositoryImpl) GetAllCats(db *sql.DB) ([]domain.Cat, error) {
	return nil, nil
}
