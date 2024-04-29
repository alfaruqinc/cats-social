package models

import (
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

type User struct {
	Id           uuid.UUID    `json:"id"`
	Email        string       `json:"email"`
	Name         string       `json:"name"`
	Password     string       `json:"password"`
	TokenService TokenService // Menambahkan dependensi tokenService
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

func (u *User) GenerateToken(signingKey []byte) (string, error) {
	claims := jwt.MapClaims{
		"id":   u.Id,
		"name": u.Name,
		"exp":  time.Now().Add(time.Hour * 8).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func (u *User) signToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := u.TokenService.GetJWTSecret()

	tokenString, _ := token.SignedString([]byte(secretKey))

	return tokenString
}

func (u *User) parseToken(tokenString string) (*jwt.Token, MessageErr) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, invalidTokenErr
		}

		secretKey := u.TokenService.GetJWTSecret()

		return []byte(secretKey), nil
	})

	if err != nil {

		return nil, invalidTokenErr
	}

	return token, nil
}

func (u *User) bindTokenToUserEntity(claim jwt.MapClaims) MessageErr {

	if uuid, ok := claim["id"].(uuid.UUID); !ok {
		return invalidTokenErr
	} else {
		u.Id = uuid
	}

	if name, ok := claim["name"].(string); !ok {
		return invalidTokenErr
	} else {
		u.Name = name
	}

	return nil
}

func (u *User) ValidateToken(bearerToken string) MessageErr {
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

	return err
}
