package domain

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenService interface {
	GetJWTSecret() string
	GetBcryptSalt() string
}

type tokenService struct {
	JWTSecret  string
	BcryptSalt string
}

func NewTokenService() TokenService {
	return &tokenService{
		JWTSecret:  os.Getenv("JWT_SECRET"),
		BcryptSalt: os.Getenv("BCRYPT_SALT"),
	}
}

func (t *tokenService) GetJWTSecret() string {
	return t.JWTSecret
}

func (t *tokenService) GetBcryptSalt() string {
	return t.BcryptSalt
}

type NewUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	TokenService string `json:"accessToken"`
}

type User struct {
	Id           uuid.UUID    `json:"id" db:"id"`
	Email        string       `json:"email" db:"email" validate:"required,email"`
	Name         string       `json:"name" db:"name" validate:"required,min=5,max=50"`
	Password     string       `json:"password" db:"password" validate:"required,min=5,max=15"`
	TokenService TokenService `json:"accessToken"`
	CreatedAt    time.Time    `json:"createdAt" db:"created_at"`
}

func NewUser() *User {
	id := uuid.New()
	token := NewTokenService()

	return &User{
		Id:           id,
		TokenService: token,
	}
}

var invalidTokenErr = NewUnauthenticatedError("invalid token")

func (u *User) HashPassword() MessageErr {
	salt, err := strconv.Atoi(u.TokenService.GetBcryptSalt())
	if err != nil {
		return NewInternalServerError("SOMETHING WENT WRONG")
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(u.Password), salt)
	if err != nil {
		return NewInternalServerError("SOMETHING WENT WRONG")
	}

	u.Password = string(bs)

	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"id":    u.Id.String(),
		"email": u.Email,
		"exp":   time.Now().Add(time.Hour * 8).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(u.TokenService.GetJWTSecret()))

	return tokenString, nil
}

func (u *User) parseToken(tokenString string) (*jwt.Token, MessageErr) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, invalidTokenErr
		}

		return []byte(u.TokenService.GetJWTSecret()), nil
	})
	if err != nil {
		return nil, invalidTokenErr
	}

	return token, nil
}

func (u *User) bindTokenToUserEntity(claim jwt.MapClaims) MessageErr {
	idString, ok := claim["id"].(string)
	if !ok {
		return invalidTokenErr
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return invalidTokenErr
	}
	u.Id = id

	email, ok := claim["email"].(string)
	if !ok {
		return invalidTokenErr
	}

	u.Email = email

	return nil
}

func (u *User) ValidateToken(bearerToken string) error {
	isBearer := strings.HasPrefix(bearerToken, "Bearer")

	if !isBearer {
		return NewUnauthenticatedError("token should be Bearer")
	}

	splitToken := strings.Fields(bearerToken)

	if len(splitToken) != 2 {
		return NewUnauthenticatedError("invalid token")
	}

	tokenString := splitToken[1]

	token, err := u.parseToken(tokenString)
	if err != nil {
		return err
	}

	var mapClaims jwt.MapClaims

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return invalidTokenErr
	} else {
		mapClaims = claims
	}
	err = u.bindTokenToUserEntity(mapClaims)
	if err != nil {
		return err
	}

	return nil
}
