package models

import (
	"cats-social/internal/database"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uuid     string `json:"uuid"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

var invalidTokenErr = NewUnauthenticatedError("invalid token")

func CreateUserUUID() string {
	return uuid.New().String()
}

func (u *User) HashPassword() MessageErr {
	salt, err := strconv.Atoi(database.BcryptSalt)

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

func (u *User) GenerateToken() string {
	claims := u.tokenClaim()

	return u.signToken(claims)
}

func (u *User) tokenClaim() jwt.MapClaims {
	return jwt.MapClaims{
		"uuid": u.Uuid,
		"name": u.Name,
		"exp":  time.Now().Add(time.Hour * 8).Unix(),
	}
}

func (u *User) signToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := database.JWTSecret

	tokenString, _ := token.SignedString([]byte(secretKey))

	return tokenString
}

func (u *User) parseToken(tokenString string) (*jwt.Token, MessageErr) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, invalidTokenErr
		}

		secretKey := database.JWTSecret

		return []byte(secretKey), nil
	})

	if err != nil {

		return nil, invalidTokenErr
	}

	return token, nil
}

func (u *User) bindTokenToUserEntity(claim jwt.MapClaims) MessageErr {

	if uuid, ok := claim["uuid"].(string); !ok {
		return invalidTokenErr
	} else {
		u.Uuid = uuid
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
