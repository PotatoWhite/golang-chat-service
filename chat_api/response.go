package chat_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

func noContent(c *gin.Context, err error) {
	c.JSON(http.StatusNoContent, gin.H{
		"error": err.Error(),
	})
}

func internalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": fmt.Sprintf("internal server error: %v", err),
	})
}
