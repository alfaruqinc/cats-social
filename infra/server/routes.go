package server

import (
	"cats-social/internal/auth"
	"cats-social/internal/handler"
	"cats-social/internal/repository"
	"cats-social/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	// dependency injection
	catMatchRepository := repository.NewCatMatchRepository()
	catMatchService := service.NewCatMatchService(catMatchRepository, s.db)
	catMatchHandler := handler.NewCatMatchHandler(catMatchService)

	r := gin.Default()

	// version 1
	apiV1 := r.Group("/v1")

	// user
	user := apiV1.Group("/user")
	user.POST("/register", handler.HandleNewUser(s.db))
	user.POST("/login", handler.HandleLogin(s.db))

	// cat
	cat := apiV1.Group("/cat")
	cat.Use(auth.NewAuth(repository.NewUserPg()).Authentication(s.db))
	cat.GET("", handler.HandleGetAllCats(s.db))
	cat.POST("", handler.HandleAddNewCat(s.db))
	cat.PUT(":catId", handler.HandleUpdateCat(s.db))
	cat.DELETE(":catId", handler.HandleDeleteCat(s.db))

	// cat match
	catMatch := cat.Group("/match")
	catMatch.POST("", catMatchHandler.CreateCatMatch())
	catMatch.GET("", catMatchHandler.GetCatMatchesByIssuerOrReceiverID())

	return r
}
