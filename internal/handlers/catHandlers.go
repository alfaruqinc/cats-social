package handlers

import (
	"github.com/gin-gonic/gin"
)

func HandleAddNewCat(c *gin.Context) {
	c.String(201, "ADD NEW CAT")
}
