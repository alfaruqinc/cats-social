package handlers

import (
	"cats-social/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleAddNewCat() gin.HandlerFunc {
	return func(c *gin.Context) {
		catBody := models.NewCat()
		if err := c.ShouldBindJSON(&catBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"message": "success", "data": &catBody})
	}
}
