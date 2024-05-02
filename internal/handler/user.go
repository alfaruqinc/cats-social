package handler

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type userResponse struct {
	Email       string `json:"email" binding:"required,email"`
	Name        string `json:"name" binding:"required"`
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

		err = repository.NewUserPg().CreateNewUser(db, userBody)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, domain.NewInternalServerError("SOMETHING WENT WRONG"))
			panic(err)
		}

		token, err := userBody.GenerateToken()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "something went wrong")
			panic(err)
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

		err := validateLoginUser(*userBody, db)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		user, err := repository.NewUserPg().GetByEmail(db, userBody.Email)
		if err != nil {
			ctx.JSON(http.StatusNotFound, domain.NewNotFoundError("email has been used"))
			return
		}

		isValidPassword := user.ComparePassword(userBody.Password)
		if !isValidPassword {
			ctx.JSON(http.StatusBadRequest, domain.NewBadRequest("invalid email or password"))
			return
		}

		token, err := userBody.GenerateToken()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, domain.NewInternalServerError(err.Error()))
			panic(err)
		}

		res := &userResponse{
			Email:       user.Email,
			Name:        user.Name,
			AccessToken: token,
		}

		ctx.JSON(http.StatusOK, domain.NewStatusOk("User Logged Successfully", res))
	}
}

func validateLoginUser(user domain.User, db *sql.DB) error {
	if len(user.Email) < 1 {
		err := errors.New("email not be empty")
		return err
	}

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", user.Email)
	if err := row.Scan(&count); err != nil {
		return err
	}

	if !validEmail(user.Email) {
		err := errors.New("invalid email format")
		return err
	}
	return nil
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
