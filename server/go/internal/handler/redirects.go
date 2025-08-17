package handler

import (
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/gin-gonic/gin"
)

func CheckRedirectsHandler(c *gin.Context) {
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

	redirected, chain, finalURL, chainLength, err := checks.CheckRedirects(rawURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_redirected": redirected,
		"chain":         chain,
		"finalURL":      finalURL,
		"chainLength":   chainLength,
	})
}
