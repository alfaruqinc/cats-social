package handlers

import (
	"cats-social/internal/models"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type response struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func HandleAddNewCat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		catBody := models.NewCat()
		if err := c.ShouldBindJSON(&catBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(catBody.Name) < 1 || len(catBody.Name) > 30 {
			err := errors.New("name length should between 1 and 30 characters")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO cats (id, created_at, name, race, sex, age_in_month, description, image_urls)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err := db.Exec(query, catBody.ID, catBody.CreatedAt, catBody.Name, catBody.Race, catBody.Sex, catBody.AgeInMonth, catBody.Description, catBody.ImageUrls)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		res := &response{
			Id:        catBody.ID,
			CreatedAt: catBody.CreatedAt,
		}

		c.JSON(201, gin.H{"message": "success", "data": res})
	}
}

func HandleGetAllCats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT id, name, race, sex, age_in_month, image_urls, description, created_at
		FROM cats
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		cats := []*models.Cat{}
		m := pgtype.NewMap()

		for rows.Next() {
			cat := &models.Cat{}

			err = rows.Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, m.SQLScanner(&cat.ImageUrls), &cat.Description, &cat.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			cats = append(cats, cat)
		}

		c.JSON(http.StatusOK, gin.H{"message": "success", "data": &cats})
	}
}
