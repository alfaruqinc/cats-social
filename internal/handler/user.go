package handler

import (
	"cats-social/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type userResponse struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func HandleNewUser(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userBody := domain.NewUser()

		if err := ctx.ShouldBindJSON(&userBody); err != nil {
			ctx.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		err := ValidateUserRequest(*userBody, db)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		userBody.HashPassword()

		query := `INSERT INTO users (id, name, email, password)
		VALUES ($1, $2, $3, $4)`

		_, err = db.Exec(query, userBody.Id, userBody.Name, userBody.Email, userBody.Password)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, domain.NewInternalServerError(err.Error()))
			return
		}

		token, err := domain.NewUser().GenerateToken()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, domain.NewInternalServerError(err.Error()))
		}

		res := &userResponse{
			Email:       userBody.Email,
			Name:        userBody.Name,
			AccessToken: token,
		}

		ctx.JSON(http.StatusCreated, domain.NewStatusCreated("User registered successfully", res))
	}
}

func HandleLogin(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userBody := domain.NewUser()

		if err := ctx.ShouldBindJSON(&userBody); err != nil {
			ctx.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		var user domain.User

		row := db.QueryRow(`SELECT id, name, email, password FROM users WHERE email = $1`, userBody.Email)
		err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.JSON(http.StatusNotFound, domain.NewNotFoundError("Invalid email or password"))
				return
			}
			ctx.JSON(http.StatusInternalServerError, domain.NewInternalServerError(err.Error()))
			return
		}

		if !user.ComparePassword(userBody.Password) {
			ctx.JSON(http.StatusUnauthorized, domain.NewUnauthorizedError("Invalid email or password"))
			return
		}

		token, err := domain.NewUser().GenerateToken()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, domain.NewInternalServerError(err.Error()))
			return
		}

		res := &userResponse{
			Email:       user.Email,
			Name:        user.Name,
			AccessToken: token,
		}

		ctx.JSON(http.StatusOK, domain.NewStatusOk("User Logged Successfully", res))
		fmt.Print(res)
	}
}

func ValidateUserRequest(userBody domain.User, db *sql.DB) error {
	if !validEmail(userBody.Email) {
		err := errors.New("invalid email format")
		return err
	}

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", userBody.Email)
	if err := row.Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		err := domain.NewConflictError("email has been used")
		return err
	}

	if userBody.Email == domain.NewUser().Email {
		err := domain.NewConflictError("email has been used")
		return err
	}

	if len(userBody.Name) < 1 || len(userBody.Name) > 100 {
		err := errors.New("name length should be between 1 and 100")
		return err
	}

	if len(userBody.Password) > 8 {
		err := errors.New("minimum password is 8 length")
		return err
	}
	return nil
}

func validEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
