package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "SafeSurf API",
		"status":  "running",
		"version": "v1",
	})
}
