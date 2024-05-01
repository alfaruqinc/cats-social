package repository

import (
	"cats-social/internal/domain"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type CatRepository interface {
	CreateCat(db *sql.DB, cat *domain.Cat) error
	GetAllCats(db *sql.DB) ([]domain.Cat, error)
	UpdateCat(db *sql.DB, cat *domain.Cat) error
	CheckCatIdExists(db *sql.DB, catId uuid.UUID) error
	DeleteCat(db *sql.DB, catId uuid.UUID) error
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
		return err
	}

	return nil
}

func (c *CatRepositoryImpl) GetAllCats(db *sql.DB) ([]domain.Cat, error) {
	return nil, nil
}

func (c *CatRepositoryImpl) UpdateCat(db *sql.DB, cat *domain.Cat) error {
	query := `
		UPDATE cats
		SET name = $2,
			race = $3,
			sex = $4,
			age_in_month = $5,
			description = $6,
			image_urls = $7
		WHERE id = $1
	`
	_, err := db.Exec(query, cat.ID, cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, cat.ImageUrls)
	if err != nil {
		return err
	}

	return nil
}

func (c *CatRepositoryImpl) DeleteCat(db *sql.DB, catId uuid.UUID) error {
	query := `
		UPDATE cats
		SET deleted = $2
		WHERE id = $1
	`

	_, err := db.Exec(query, catId, true)
	if err != nil {
		return err
	}

	return nil
}

func (c *CatRepositoryImpl) CheckCatIdExists(db *sql.DB, catId uuid.UUID) error {
	queryCheckCatId := `
		SELECT EXISTS (
			SELECT 1
			FROM cats
			WHERE id = $1
		)
	`
	var catIdExists bool
	row := db.QueryRow(queryCheckCatId, catId)
	err := row.Scan(&catIdExists)
	if err != nil {
		return err
	}

	if !catIdExists {
		return errors.New("cat is not found")
	}

	return nil
}
