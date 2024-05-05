package handler

import (
	"cats-social/internal/domain"
	"cats-social/internal/service"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CatHandler interface {
	CreateCat() gin.HandlerFunc
	GetAllCats() gin.HandlerFunc
	UpdateCat() gin.HandlerFunc
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
			err, _ := err.(domain.MessageErr)
			ctx.JSON(err.Status(), err)
			if err.Status() > 499 {
				panic(err)
			}
		}

		res := &domain.CreateCatResponse{
			ID:        catBody.ID,
			CreatedAt: catBody.CreatedAt,
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "success", "data": res})
	}
}

func (c *catHandler) GetAllCats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userReq, _ := ctx.Get("userData")
		user := userReq.(*domain.User)

		cats, err := c.catSerivce.GetAllCats(user, ctx.Request.URL.Query())
		if err != nil {
			ctx.JSON(err.Status(), err)
			if err.Status() > 499 {
				panic(err)
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "success", "data": &cats})
	}
}

func (c *catHandler) UpdateCat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		catId := ctx.Param("catId")
		parsedCatId, err := uuid.Parse(catId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, domain.NewNotFoundError("cat is not found"))
			return
		}

		// TODO: create dto for update cat request
		catBody := domain.NewCat()
		if err := ctx.ShouldBindJSON(catBody); err != nil {
			ctx.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		userReq, _ := ctx.Get("userData")
		user := userReq.(*domain.User)

		err = validateRequestBody(*catBody)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		catBody.ID = parsedCatId

		err = c.catSerivce.UpdateCat(user, catBody)
		if err != nil {
			err, _ := err.(domain.MessageErr)
			ctx.JSON(err.Status(), err)
			if err.Status() > 499 {
				panic(err)
			}
		}

		updatedAt := time.Now().Format(time.RFC3339)
		parsedUpdatedAt, _ := time.Parse(time.RFC3339, updatedAt)

		updatedCat := domain.UpdateCatResponse{
			ID:        parsedCatId,
			UpdatedAt: parsedUpdatedAt,
		}

		ctx.JSON(http.StatusOK, domain.NewStatusOk("success", updatedCat))
	}
}

func validateRequestBody(body domain.Cat) error {
	if len(body.Name) < 1 || len(body.Name) > 30 {
		err := errors.New("name length should be between 1 and 30 characters")
		return err
	}

	if slices.Contains(domain.CatRace, body.Race) != true {
		err := errors.New("accepted race is only Persian, Maine Coon, Siamese, Ragdoll, Bengal, Sphynx, British Shorthair, Abyssinian, Scottish Fold, Birman")
		return err
	}

	if slices.Contains(domain.CatSex, body.Sex) != true {
		err := errors.New("accepted sex is only male and female")
		return err
	}

	if body.AgeInMonth < 1 || body.AgeInMonth > 120082 {
		err := errors.New("your cat's age is minimum 1 month and maximum 120082 month")
		return err
	}

	if len(body.Description) < 1 || len(body.Description) > 200 {
		err := errors.New("description length should be between 1 and 200 characters")
		return err
	}

	if len(body.ImageUrls) < 1 {
		err := errors.New("image urls at least have 1 image")
		return err
	}

	for _, imageUrl := range body.ImageUrls {
		if len(imageUrl) < 1 {
			err := errors.New("image urls cannot have empty item")
			return err
		}

		_, err := url.ParseRequestURI(imageUrl)
		if err != nil {
			err := errors.New("image url should have valid url")
			return err
		}
	}
	return nil
}
