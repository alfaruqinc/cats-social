package handler

import (
	"cats-social/internal/service"
)

type CatHandler interface{}

type catHandler struct {
	catSerivce service.CatService
}

func NewCatHandler(catService service.CatService) CatHandler {
	return &catHandler{
		catSerivce: catService,
	}
}
