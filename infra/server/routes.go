package server

import (
	"cats-social/internal/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// version 1
	apiV1 := r.Group("/v1")

	// user
	user := apiV1.Group("/user")
	user.GET("", func(c *gin.Context) { c.String(200, "HALO USER") })

	// cat
	cat := apiV1.Group("/cat")
	cat.GET("", handler.HandleGetAllCats(s.db))
	cat.POST("", handler.HandleAddNewCat(s.db))

	// cat match
	catMatch := cat.Group("/match")
	catMatch.GET("", func(c *gin.Context) { c.String(200, "HALO CAT MATCH") })

	return r
}
