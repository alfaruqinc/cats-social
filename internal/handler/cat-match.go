package handler

import (
	"cats-social/internal/domain"
	"cats-social/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CatMatchHandler interface {
	CreateCatMatch() gin.HandlerFunc
	GetCatMatchesByIssuerOrReceiverID() gin.HandlerFunc
	UpdateCatMatchByID() gin.HandlerFunc
	DeleteCatMatchByID() gin.HandlerFunc
	ApproveCatMatch() gin.HandlerFunc
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
		result, err := c.catMatchService.CreateCatMatch(ctx, nil)
		if err != nil {
			ctx.JSON(err.Status(), gin.H{
				"message": err.Message(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": result,
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
			"message": "Get Cat Match By ID",
		})
	}
}

func (c *catMatchHandler) UpdateCatMatchByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Update Cat Match By ID",
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

		var body domain.CatMatch
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusInternalServerError, domain.NewInternalServerError("something went wrong"))
			panic(err)
		}

		err := c.catMatchService.ApproveCatMatchByMatchCatID(ctx, user.Id.String(), body.MatchCatID.String())
		if err, ok := err.(domain.MessageErr); ok {
			ctx.JSON(err.Status(), err)
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "success approve cat match",
		})
	}
}
