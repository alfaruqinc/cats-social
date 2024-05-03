package handler

import (
	"cats-social/internal/domain"
	"cats-social/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CatMatchHandler interface {
	CreateCatMatch() gin.HandlerFunc
	GetCatMatchesByIssuerOrReceiverID() gin.HandlerFunc
	DeleteCatMatchByID() gin.HandlerFunc
	ApproveCatMatch() gin.HandlerFunc
	RejectCatMatch() gin.HandlerFunc
}

type catMatchHandler struct {
	catMatchService service.CatMatchService
}

func NewCatMatchHandler(catMatchService service.CatMatchService) CatMatchHandler {
	return &catMatchHandler{
		catMatchService: catMatchService,
	}
}

func (c *catMatchHandler) CreateCatMatch() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userReq, _ := ctx.Get("userData")
		user := userReq.(*domain.User)

		body := domain.NewCatMatch()
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusInternalServerError, domain.NewInternalServerError("something went wrong"))
		}

		body.IssuedByID = user.Id

		err := c.catMatchService.CreateCatMatch(ctx, body)
		if err != nil {
			ctx.JSON(err.Status(), gin.H{
				"message": err.Message(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "success create cat match",
		})
	}
}

func (c *catMatchHandler) GetCatMatchesByIssuerOrReceiverID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: get id of user from token
		result, err := c.catMatchService.GetCatMatchesByIssuerOrReceiverID(ctx, "")
		if err != nil {
			ctx.JSON(err.Status(), gin.H{
				"message": err.Message(),
			})

			return
		}

		ctx.JSON(200, gin.H{
			"data":    result,
			"message": "success",
		})
	}
}

func (c *catMatchHandler) DeleteCatMatchByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		catMatchId := ctx.Param("id")

		userReq, _ := ctx.Get("userData")
		user := userReq.(*domain.User)

		_, err := uuid.Parse(catMatchId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, domain.NewNotFoundError("Cat match request is not found"))
		}

		err = c.catMatchService.DeleteCatMatchByID(ctx, catMatchId, user.Id.String())
		if err, ok := err.(domain.MessageErr); ok {
			ctx.JSON(err.Status(), err)
			return
		}

		ctx.JSON(200, gin.H{
			"message": "success delete cat match request",
		})
	}
}

func (c *catMatchHandler) ApproveCatMatch() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userReq, _ := ctx.Get("userData")
		user := userReq.(*domain.User)

		body := struct {
			MatchId string `json:"matchId"`
		}{}
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusInternalServerError, domain.NewInternalServerError("something went wrong"))
			panic(err)
		}

		_, err := uuid.Parse(body.MatchId)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusNotFound, domain.NewNotFoundError("Cat match request is not found"))
			return
		}

		err = c.catMatchService.ApproveCatMatch(ctx, user.Id.String(), body.MatchId)
		if err, ok := err.(domain.MessageErr); ok {
			ctx.JSON(err.Status(), err)
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "success approve cat match",
		})
	}
}

func (c *catMatchHandler) RejectCatMatch() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userReq, _ := ctx.Get("userData")
		user := userReq.(*domain.User)

		body := struct {
			MatchId string `json:"matchId"`
		}{}
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusInternalServerError, domain.NewInternalServerError("something went wrong"))
			panic(err)
		}

		_, err := uuid.Parse(body.MatchId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, domain.NewNotFoundError("Cat match request is not found"))
			return
		}

		err = c.catMatchService.RejectCatMatch(ctx, user.Id.String(), body.MatchId)
		if err, ok := err.(domain.MessageErr); ok {
			ctx.JSON(err.Status(), err)
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "success reject cat match",
		})
	}
}
