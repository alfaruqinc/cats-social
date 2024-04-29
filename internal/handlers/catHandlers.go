package handlers

import (
	"cats-social/internal/models"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
