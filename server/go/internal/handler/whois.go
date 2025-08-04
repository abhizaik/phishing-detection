package handler

import (
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/gin-gonic/gin"
)

func WhoisHandler(c *gin.Context) {
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

	domain, err := checks.GetDomain(rawURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not extract domain"})
		return
	}

	whoisData, age, err := checks.GetWhoisData(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "whois lookup failed",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domain":     domain,
		"age":        age,
		"created_at": whoisData.Domain.CreatedDate,
		"expires_at": whoisData.Domain.ExpirationDate,
		"registrar":  whoisData.Registrar.Name,
		"raw_data":   whoisData,
	})
}
