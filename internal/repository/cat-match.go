package repository

import (
	"cats-social/internal/domain"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type CatMatchRepository interface {
	CreateCatMatch(ctx context.Context, tx *sql.Tx, catMatch *domain.CatMatch) (*domain.CatMatch, error)
	GetCatMatchByID(ctx context.Context, tx *sql.Tx, id string) (*domain.CatMatch, error)
	GetCatMatchesByIssuerOrReceiverID(ctx context.Context, tx *sql.Tx, id string) ([]domain.CatMatch, error)
	UpdateCatMatchByID(ctx context.Context, tx *sql.Tx, id string, catMatch *domain.CatMatch) error
	DeleteCatMatchByID(ctx context.Context, tx *sql.Tx, id string) error
	GetStatusCatMatchByID(ctx context.Context, tx *sql.Tx, id string) (string, error)
	ApproveCatMatch(ctx context.Context, tx *sql.Tx, userId string, matchId string) error
	RejectCatMatch(ctx context.Context, tx *sql.Tx, userId string, matchId string) error
	CanDeleteCatMatch(ctx context.Context, tx *sql.Tx, id string, userId string) (bool, error)
	CheckIfUserIsReceiver(ctx context.Context, tx *sql.Tx, id string, userId string) (bool, error)
}

type catMatchRepository struct{}

func NewCatMatchRepository() CatMatchRepository {
	return &catMatchRepository{}
}

func (c *catMatchRepository) CreateCatMatch(ctx context.Context, tx *sql.Tx, catMatch *domain.CatMatch) (*domain.CatMatch, error) {
	query := `INSERT INTO cat_matches (id, created_at, issued_by_id, match_cat_id, user_cat_id, message) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.ExecContext(ctx, query, catMatch.ID, catMatch.CreatedAt, catMatch.IssuedByID, catMatch.MatchCatID, catMatch.UserCatID, catMatch.Message)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *catMatchRepository) GetCatMatchByID(ctx context.Context, tx *sql.Tx, id string) (*domain.CatMatch, error) {
	query := `
		SELECT	cm.id, cm.created_at, cm.issued_by_id, cm.match_cat_id, cm.user_cat_id, cm.message, cm.status, 
				u.name as issued_by_name, u.email as issued_by_email, u.created_at as issued_by_created_at, 
				ca.name as match_cat_name, ca.race as match_cat_race, ca.sex as match_cat_sex, ca.description as match_cat_description, ca.age_in_month as match_cat_age_in_month, ca.image_urls as match_cat_image_urls, ca.has_matched as match_cat_has_matched , ca.created_at as match_cat_created_at,
				cb.name as user_cat_name, cb.race as user_cat_race, cb.sex as match_cat_sex, cb.description as user_cat_description, cb.age_in_month as user_cat_age_in_month, cb.image_urls as user_cat_image_urls, cb.has_matched as user_cat_has_matched , cb.created_at as user_cat_created_at
		FROM cat_matches cm
		INNER JOIN users u ON cat_matches.issued_by_id = users.id
		INNER JOIN cats ca ON cat_matches.match_cat_id = cats.id
		INNER JOIN cats cb ON cat_matches.user_cat_id = cats.id
		WHERE id = $1`

	row := tx.QueryRowContext(ctx, query, id)

	var catMatch domain.CatMatch
	err := row.Scan(
		&catMatch.ID,
		&catMatch.CreatedAt,
		&catMatch.IssuedByID,
		&catMatch.MatchCatID,
		&catMatch.UserCatID,
		&catMatch.Message,
		&catMatch.Status,
		&catMatch.IssuedBy.Name,
		&catMatch.IssuedBy.Email,
		&catMatch.IssuedBy.CreatedAt,
		&catMatch.MatchCat.Name,
		&catMatch.MatchCat.Race,
		&catMatch.MatchCat.Sex,
		&catMatch.MatchCat.Description,
		&catMatch.MatchCat.AgeInMonth,
		&catMatch.MatchCat.ImageUrls,
		&catMatch.MatchCat.HasMatched,
		&catMatch.MatchCat.CreatedAt,
		&catMatch.UserCat.Name,
		&catMatch.UserCat.Race,
		&catMatch.UserCat.Sex,
		&catMatch.UserCat.Description,
		&catMatch.UserCat.AgeInMonth,
		&catMatch.UserCat.ImageUrls,
		&catMatch.UserCat.HasMatched,
		&catMatch.UserCat.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if catMatch.ID == uuid.Nil {
		return &domain.CatMatch{}, nil
	}

	return &catMatch, nil
}

func (c *catMatchRepository) GetCatMatchesByIssuerOrReceiverID(ctx context.Context, tx *sql.Tx, id string) ([]domain.CatMatch, error) {
	query := `
		SELECT	cm.id, cm.created_at, cm.issued_by_id, cm.match_cat_id, cm.user_cat_id, cm.message, cm.status, 
			u.name as issued_by_name, u.email as issued_by_email, u.created_at as issued_by_created_at, 
			ca.name as match_cat_name, ca.race as match_cat_race, ca.sex as match_cat_sex, ca.description as match_cat_description, ca.age_in_month as match_cat_age_in_month, ca.image_urls as match_cat_image_urls, ca.has_matched as match_cat_has_matched , ca.created_at as match_cat_created_at,
			cb.name as user_cat_name, cb.race as user_cat_race, cb.sex as match_cat_sex, cb.description as user_cat_description, cb.age_in_month as user_cat_age_in_month, cb.image_urls as user_cat_image_urls, cb.has_matched as user_cat_has_matched , cb.created_at as user_cat_created_at
		FROM cat_matches cm
		INNER JOIN users u ON cat_matches.issued_by_id = users.id
		INNER JOIN cats ca ON cat_matches.match_cat_id = cats.id
		INNER JOIN cats cb ON cat_matches.user_cat_id = cats.id
		WHERE u.id = $1 OR ca.owned_by = $1
	`

	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var catMatches []domain.CatMatch
	for rows.Next() {
		var catMatch domain.CatMatch
		err := rows.Scan(
			&catMatch.ID,
			&catMatch.CreatedAt,
			&catMatch.IssuedByID,
			&catMatch.MatchCatID,
			&catMatch.UserCatID,
			&catMatch.Message,
			&catMatch.Status,
			&catMatch.IssuedBy.Name,
			&catMatch.IssuedBy.Email,
			&catMatch.IssuedBy.CreatedAt,
			&catMatch.MatchCat.Name,
			&catMatch.MatchCat.Race,
			&catMatch.MatchCat.Sex,
			&catMatch.MatchCat.Description,
			&catMatch.MatchCat.AgeInMonth,
			&catMatch.MatchCat.ImageUrls,
			&catMatch.MatchCat.HasMatched,
			&catMatch.MatchCat.CreatedAt,
			&catMatch.UserCat.Name,
			&catMatch.UserCat.Race,
			&catMatch.UserCat.Sex,
			&catMatch.UserCat.Description,
			&catMatch.UserCat.AgeInMonth,
			&catMatch.UserCat.ImageUrls,
			&catMatch.UserCat.HasMatched,
			&catMatch.UserCat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		catMatches = append(catMatches, catMatch)
	}

	return catMatches, nil
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

func (c *catMatchRepository) DeleteCatMatchByID(ctx context.Context, tx *sql.Tx, id string) error {
	query := `DELETE FROM cat_matches WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *catMatchRepository) GetStatusCatMatchByID(ctx context.Context, tx *sql.Tx, id string) (string, error) {
	query := `SELECT status FROM cat_matches WHERE id = $1`

	var status string
	err := tx.QueryRowContext(ctx, query, id).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (c *catMatchRepository) CanDeleteCatMatch(ctx context.Context, tx *sql.Tx, id string, userId string) (bool, error) {
	query := `SELECT issued_by_id = $2 FROM cat_matches WHERE id = $1`

	var canDelete bool
	err := tx.QueryRowContext(ctx, query, id, userId).Scan(&canDelete)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return canDelete, nil
}

func (c *catMatchRepository) ApproveCatMatch(ctx context.Context, tx *sql.Tx, userId string, matchId string) error {
	query := `UPDATE cat_matches SET status = 'approved' WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, matchId)
	if err != nil {
		return err
	}

	queryDeleteWaitingStatus := `
		DELETE FROM cat_matches
		WHERE status = 'waiting'
			AND match_cat_id = (
				SELECT match_cat_id
				FROM cat_matches
				WHERE id = $1
			)
	`
	_, err = tx.ExecContext(ctx, queryDeleteWaitingStatus, matchId)
	if err != nil {
		return err
	}

	querySetHasMatched := `
		UPDATE cats 
		SET has_matched = true 
		WHERE id IN (
			SELECT match_cat_id
			FROM cat_matches
			WHERE id = $1
			UNION
			SELECT user_cat_id
			FROM cat_matches
			WHERE id = $1
		)
	`
	_, err = tx.ExecContext(ctx, querySetHasMatched, matchId)
	if err != nil {
		return err
	}

	return nil
}

func (c *catMatchRepository) CheckIfUserIsReceiver(ctx context.Context, tx *sql.Tx, id string, userId string) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1 
		FROM cat_matches AS cm
		JOIN cats ON cats.id = cm.match_cat_id
		WHERE cm.id = $1
			AND cats.owned_by_id = $2
	)
	`
	var exists bool
	err := tx.QueryRowContext(ctx, query, id, userId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (c *catMatchRepository) RejectCatMatch(ctx context.Context, tx *sql.Tx, userId string, matchId string) error {
	query := `UPDATE cat_matches SET status = 'rejected' WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, matchId)
	if err != nil {
		return err
	}

	return nil
}
