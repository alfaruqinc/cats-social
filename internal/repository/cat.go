package repository

import (
	"cats-social/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type CatRepository interface {
	CreateCat(db *sql.DB, cat *domain.Cat) error
	GetAllCats(db *sql.DB) ([]domain.Cat, error)
	UpdateCat(db *sql.DB, cat *domain.Cat) error
	DeleteCat(db *sql.DB, catId uuid.UUID) error
	CheckCatIdExists(db *sql.DB, catId uuid.UUID, userId uuid.UUID) error
	CheckEditableSex(db *sql.DB, cat *domain.Cat) error
	CheckOwnerCat(ctx context.Context, tx *sql.Tx, catId uuid.UUID, userId uuid.UUID) (bool, error)
	CheckCatHasSameSex(ctx context.Context, tx *sql.Tx, cat1Id uuid.UUID, cat2Id uuid.UUID) (bool, error)
	CheckCatHasMatched(ctx context.Context, tx *sql.Tx, cat1Id uuid.UUID, cat2Id uuid.UUID) (bool, error)
	CheckCatFromSameOwner(ctx context.Context, tx *sql.Tx, cat1Id uuid.UUID, cat2Id uuid.UUID) (bool, error)
	CheckBothCatExists(ctx context.Context, tx *sql.Tx, cat1Id uuid.UUID, cat2Id uuid.UUID) (bool, error)
}

type catRepository struct{}

func NewCatRepository() CatRepository {
	return &catRepository{}
}

func (c *catRepository) CreateCat(db *sql.DB, catBody *domain.Cat) error {
	query := `INSERT INTO cats (id, created_at, name, race, sex, age_in_month, description, image_urls, owned_by_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
	_, err := db.Exec(query, catBody.ID, catBody.CreatedAt, catBody.Name, catBody.Race, catBody.Sex, catBody.AgeInMonth, catBody.Description, catBody.ImageUrls, catBody.OwnedById)
	if err != nil {
		return err
	}

	return nil
}

func (c *catRepository) GetAllCats(db *sql.DB) ([]domain.Cat, error) {
	return nil, nil
}

func (c *catRepository) UpdateCat(db *sql.DB, cat *domain.Cat) error {
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

func (c *catRepository) DeleteCat(db *sql.DB, catId uuid.UUID) error {
	query := `
		UPDATE cats
		SET deleted = true
		WHERE id = $1
	`

	_, err := db.Exec(query, catId)
	if err != nil {
		return err
	}

	return nil
}

func (c *catRepository) CheckCatIdExists(db *sql.DB, catId uuid.UUID, userId uuid.UUID) error {
	queryCheckCatId := `
		SELECT EXISTS (
			SELECT 1
			FROM cats
			WHERE id = $1 
				AND owned_by_id = $2
				AND deleted = false
		)
	`
	var catIdExists bool
	row := db.QueryRow(queryCheckCatId, catId, userId)
	err := row.Scan(&catIdExists)
	if err != nil {
		return err
	}

	if !catIdExists {
		return errors.New("cat is not found")
	}

	return nil
}

func (c *catRepository) CheckEditableSex(db *sql.DB, cat *domain.Cat) error {
	queryCheckCatId := `
		SELECT (sex != $1) as sex_diff, NOT EXISTS (
			SELECT 1
			FROM cat_matches
			WHERE user_cat_id = $2
		) as can_edit
		FROM cats
		WHERE id = $2
	`
	var sexDiff bool
	var canEdit bool
	row := db.QueryRow(queryCheckCatId, cat.Sex, cat.ID)
	err := row.Scan(&sexDiff, &canEdit)
	if err != nil {
		return err
	}

	if sexDiff && !canEdit {
		return errors.New("cannot edit sex when already requested to match")
	}

	return nil
}

func (c *catRepository) CheckOwnerCat(ctx context.Context, tx *sql.Tx, catId uuid.UUID, userId uuid.UUID) (bool, error) {
	queryCheckCatId := `
		SELECT EXISTS (
			SELECT 1
			FROM cats
			WHERE id = $1 
				AND owned_by_id = $2
				AND deleted = false
		)
	`
	var owner bool
	row := tx.QueryRowContext(ctx, queryCheckCatId, catId, userId)
	err := row.Scan(&owner)
	if err != nil {
		return false, err
	}

	return owner, nil
}

func (c *catRepository) CheckCatHasSameSex(ctx context.Context, tx *sql.Tx, cat1Id uuid.UUID, cat2Id uuid.UUID) (bool, error) {
	query := `
		SELECT c1.sex = c2.sex
		FROM cats c1
			JOIN cats c2 on c2.id != c1.id
		WHERE c1.id = $1
			AND c2.id = $2
	`
	var hasSameSex bool
	err := tx.QueryRowContext(ctx, query, cat1Id, cat2Id).Scan(&hasSameSex)
	if err != nil {
		return false, err
	}

	return hasSameSex, nil
}

func (c *catRepository) CheckCatHasMatched(ctx context.Context, tx *sql.Tx, cat1Id uuid.UUID, cat2Id uuid.UUID) (bool, error) {
	query := `
		SELECT c1.has_matched OR c2.has_matched
		FROM cats c1
			JOIN cats c2 on c2.id != c1.id
		WHERE c1.id = $1
			AND c2.id = $2
	`
	var hasMatched bool
	err := tx.QueryRowContext(ctx, query, cat1Id, cat2Id).Scan(&hasMatched)
	if err != nil {
		return false, err
	}

	return hasMatched, nil
}

func (c *catRepository) CheckCatFromSameOwner(ctx context.Context, tx *sql.Tx, cat1Id uuid.UUID, cat2Id uuid.UUID) (bool, error) {
	query := `
		SELECT c1.owned_by_id = c2.owned_by_id
		FROM cats c1
			JOIN cats c2 on c2.id != c1.id
		WHERE c1.id = $1
			AND c2.id = $2
	`
	var sameOwner bool
	err := tx.QueryRowContext(ctx, query, cat1Id, cat2Id).Scan(&sameOwner)
	if err != nil {
		return false, err
	}

	return sameOwner, nil
}

func (c *catRepository) CheckBothCatExists(ctx context.Context, tx *sql.Tx, cat1Id uuid.UUID, cat2Id uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM cats
			WHERE id = $1
				AND deleted = false
		) AND EXISTS (
			SELECT 1
			FROM cats
			WHERE id = $2
				AND deleted = false
		) as bothExists
	`
	var bothExists bool
	err := tx.QueryRowContext(ctx, query, cat1Id, cat2Id).Scan(&bothExists)
	if err != nil {
		return false, err
	}

	return bothExists, nil
}
