package server

import (
	"cats-social/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// version 1
	apiV1 := r.Group("/v1")

	// user
	user := apiV1.Group("/user")
	user.POST("/register", handlers.HandleNewUser(s.db))

	// cat
	cat := apiV1.Group("/cat")
	cat.GET("", handlers.HandleGetAllCats(s.db))
	cat.POST("", handlers.HandleAddNewCat(s.db))

	// cat match
	catMatch := cat.Group("/match")
	catMatch.GET("", func(c *gin.Context) { c.String(200, "HALO CAT MATCH") })

	return r
}
