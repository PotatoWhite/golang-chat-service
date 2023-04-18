package api_restful

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func badRequest(c *gin.Context, err error) {
	if err == nil {
		err = fmt.Errorf("bad request")
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

func noContent(c *gin.Context, err error) {
	if err == nil {
		err = fmt.Errorf("no content")
	}

	c.JSON(http.StatusNoContent, gin.H{
		"error": err.Error(),
	})
}

func internalServerError(c *gin.Context, err error) {
	if err == nil {
		err = fmt.Errorf("internal server error")
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": fmt.Sprintf("internal server error: %v", err),
	})
}
