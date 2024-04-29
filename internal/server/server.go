package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type TokenService interface {
	GetJWTSecret() string
	GetBcryptSalt() string
}

type token struct {
	JWTSecret  string
	BcryptSalt string
}

func NewTokenService() *token {
	return &token{}
}

func (t *token) GetJWTSecret() string {
	t.JWTSecret = os.Getenv("JWT_Secret")

	return t.JWTSecret
}

func (t *token) GetBcryptSalt() string {
	t.BcryptSalt = os.Getenv("Bcrypt_Salt")

	return t.BcryptSalt
}

type Server struct {
	port int

	db *sql.DB
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	connStr := os.Getenv("DB_URL")

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	NewServer := &Server{
		port: port,

		db: db,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
