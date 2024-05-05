package server

import (
	"cats-social/internal/auth"
	"cats-social/internal/handler"
	"cats-social/internal/repository"
	"cats-social/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}

func (s *Server) RegisterRoutes() http.Handler {
	// dependency injection
	catRepository := repository.NewCatRepository()
	catMatchRepository := repository.NewCatMatchRepository()

	catService := service.NewCatService(s.db, catRepository)
	catMatchService := service.NewCatMatchService(s.db, catMatchRepository, catRepository)

	catHandler := handler.NewCatHandler(catService)
	catMatchHandler := handler.NewCatMatchHandler(catMatchService)

	r := gin.Default()

	// r := gin.New()
	// r.Use(gin.Recovery())
	// r.Use(jsonLoggerMiddleware())

	// version 1
	apiV1 := r.Group("/v1")

	// user
	user := apiV1.Group("/user")
	user.POST("/register", handler.HandleNewUser(s.db))
	user.POST("/login", handler.HandleLogin(s.db))

	// cat
	cat := apiV1.Group("/cat")
	cat.Use(auth.NewAuth(repository.NewUserPg()).Authentication(s.db))

	cat.POST("", catHandler.CreateCat())
	cat.GET("", handler.HandleGetAllCats(s.db))
	cat.PUT(":catId", handler.HandleUpdateCat(s.db))
	cat.DELETE(":catId", handler.HandleDeleteCat(s.db))

	// cat match
	catMatch := cat.Group("/match")
	catMatch.POST("", catMatchHandler.CreateCatMatch())
	catMatch.GET("", catMatchHandler.GetCatMatchesByIssuerOrReceiverID())
	catMatch.POST("/approve", catMatchHandler.ApproveCatMatch())
	catMatch.POST("/reject", catMatchHandler.RejectCatMatch())
	catMatch.DELETE(":id", catMatchHandler.DeleteCatMatchByID())

	return r
}
