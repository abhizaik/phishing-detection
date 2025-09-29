package handler

import (
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/abhizaik/SafeSurf/internal/service/screenshot"
	"github.com/gin-gonic/gin"
)

func ScreenshotHandler(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "success",
			"msg":    "url query param is required",
			"file":   "",
		})
		return
	}

	_, isValid, err := checks.IsValidURL(rawURL)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "success",
			"msg":    "invalid url",
			"file":   "",
		})
		return
	}

	filePath := screenshot.TakeScreenshot(rawURL)

	c.JSON(http.StatusOK, gin.H{
		"file":   filePath,
		"msg":    "screenshot captured successfully",
		"status": "success",
	})
	return
}
