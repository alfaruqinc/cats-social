package token

import "os"

type Token struct {
	JWTSecret  string
	BcryptSalt string
}

func NewToken() *Token {
	return &Token{
		JWTSecret:  os.Getenv("JWT_Secret"),
		BcryptSalt: os.Getenv("Bcrypt_Salt"),
	}
}

func (t *Token) GetJWTSecret() string {

	return t.JWTSecret
}

func (t *Token) GetBcryptSalt() string {

	return t.BcryptSalt
}
