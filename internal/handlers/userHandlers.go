package handlers

import (
	"cats-social/internal/models"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userResponse struct {
	Email       string `json:"email" binding:"required,email"`
	Name        string `json:"name" binding:"required"`
	AccessToken string `json:"accessToken"`
}

func HandleNewUser(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userBody := models.NewUser()

		if err := ctx.ShouldBindJSON(&userBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if userBody.Email == models.NewUser().Email {
			err := errors.New("email has been used")
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if len(userBody.Name) < 1 || len(userBody.Name) > 100 {
			err := errors.New("name length should be between 1 and 100")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(userBody.Password) > 8 {
			err := errors.New("minimum password is 8 length")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userBody.HashPassword()

		query := `INSERT INTO users (id, name, email, password)
		VALUES ($1, $2, $3, $4)`

		_, err := db.Exec(query, userBody.Id, userBody.Name, userBody.Email, userBody.Password)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token, err := models.NewUser().GenerateToken()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		res := &userResponse{
			Email:       userBody.Email,
			Name:        userBody.Name,
			AccessToken: token,
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "User Register Succesfully", "data": res})
	}
}
