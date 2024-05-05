package handler

import (
	"cats-social/internal/domain"
	"cats-social/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CatHandler interface {
	CreateCat() gin.HandlerFunc
}

type catHandler struct {
	catSerivce service.CatService
}

func NewCatHandler(catService service.CatService) CatHandler {
	return &catHandler{
		catSerivce: catService,
	}
}

func (c *catHandler) CreateCat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userReq, _ := ctx.Get("userData")
		user := userReq.(*domain.User)

		catBody := domain.NewCat()
		if err := ctx.ShouldBindJSON(&catBody); err != nil {
			ctx.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		err := validateRequestBody(*catBody)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		catBody.OwnedById = user.Id

		err = c.catSerivce.CreateCat(catBody)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "something went wrong")
			panic(err)
		}

		res := &domain.CreateCatResponse{
			ID:        catBody.ID,
			CreatedAt: catBody.CreatedAt,
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "success", "data": res})
	}
}
