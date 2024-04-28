package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	// version 1
	apiV1 := r.Group("/v1")

	// user
	user := apiV1.Group("/user")
	user.GET("", func(c *gin.Context) { c.String(200, "HALO USER") })

	// cat
	cat := apiV1.Group("/cat")
	cat.GET("", func(c *gin.Context) { c.String(200, "HALO CAT") })

	// cat match
	catMatch := cat.Group("/match")
	catMatch.GET("", func(c *gin.Context) { c.String(200, "HALO CAT MATCH") })

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
