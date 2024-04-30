package main

import (
	"cats-social/infra/server"
	"database/sql"
	"fmt"
)

type Server struct {
	port int

	db *sql.DB
}

func main() {
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
