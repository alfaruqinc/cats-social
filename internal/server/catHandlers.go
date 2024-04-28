package server

import "github.com/gin-gonic/gin"

func handleAddNewCat(c *gin.Context) {
	c.String(201, "ADD NEW CAT")
}
