package handler

import (
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/gin-gonic/gin"
)

func CheckStatusCodeHandler(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url query param is required"})
		return
	}

	_, isValid, err := checks.IsValidURL(rawURL)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	statusCode, statusText, isSuccess, isRedirect, err := checks.GetStatusCode(rawURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": statusCode,
		"status_text": statusText,
		"is_success":  isSuccess,
		"is_redirect": isRedirect,
	})
}
