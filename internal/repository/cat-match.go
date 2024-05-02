package repository

import (
	"cats-social/internal/domain"
	"context"
	"database/sql"
)

type CatMatchRepository interface {
	CreateCatMatch(ctx context.Context, tx *sql.Tx, catMatch *domain.CatMatch) (*domain.CatMatch, error)
	GetCatMatches(ctx context.Context, tx *sql.Tx) ([]domain.CatMatch, error)
	GetCatMatchByID(ctx context.Context, tx *sql.Tx, id string) (*domain.CatMatch, error)
	UpdateCatMatchByID(ctx context.Context, tx *sql.Tx, id string, catMatch *domain.CatMatch) error
	DeleteCatMatch(ctx context.Context, tx *sql.Tx, id string) error
}

type catMatchRepository struct {
}

func NewCatMatchRepository() CatMatchRepository {
	return &catMatchRepository{}
}

func (c *catMatchRepository) CreateCatMatch(ctx context.Context, tx *sql.Tx, catMatch *domain.CatMatch) (*domain.CatMatch, error) {
	query := `INSERT INTO cat_matches (id, created_at, issued_by_id, match_cat_id, user_cat_id, message, status) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return &domain.CatMatch{}, err
	}

	return nil, nil
}

func (c *catMatchRepository) GetCatMatches(ctx context.Context, tx *sql.Tx) ([]domain.CatMatch, error) {
	query := `SELECT * FROM cat_matches`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var catMatches []domain.CatMatch
	for rows.Next() {
		var catMatch domain.CatMatch
		err := rows.Scan(&catMatch.ID, &catMatch.CreatedAt, &catMatch.IssuedByID, &catMatch.MatchCatID, &catMatch.UserCatID, &catMatch.Message, &catMatch.Status)
		if err != nil {
			return nil, err
		}
		catMatches = append(catMatches, catMatch)
	}

	return catMatches, nil
}

func (c *catMatchRepository) GetCatMatchByID(ctx context.Context, tx *sql.Tx, id string) (*domain.CatMatch, error) {
	query := `SELECT * FROM cat_matches WHERE id = $1`

	row := tx.QueryRowContext(ctx, query, id)

	var catMatch domain.CatMatch
	err := row.Scan(&catMatch.ID, &catMatch.CreatedAt, &catMatch.IssuedByID, &catMatch.MatchCatID, &catMatch.UserCatID, &catMatch.Message, &catMatch.Status)
	if err != nil {
		return nil, err
	}

	return &catMatch, nil
}

func (c *catMatchRepository) UpdateCatMatchByID(ctx context.Context, tx *sql.Tx, id string, catMatch *domain.CatMatch) error {
	// currently only update status
	query := `UPDATE cat_matches SET status = $1 WHERE id = $2`

	_, err := tx.ExecContext(ctx, query, catMatch.Status, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *catMatchRepository) DeleteCatMatch(ctx context.Context, tx *sql.Tx, id string) error {
	query := `DELETE FROM cat_matches WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
