package handler

import (
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/gin-gonic/gin"
)

func CheckTrustedTLDHandler(c *gin.Context) {
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

	trusted, icann := checks.IsTrustedTld(domain)
	c.JSON(http.StatusOK, gin.H{
		"is_trusted_tld": trusted,
		"is_icann":       icann,
	})
}

func CheckRiskyTLDHandler(c *gin.Context) {
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

	risky, icann := checks.IsRiskyTld(domain)
	c.JSON(http.StatusOK, gin.H{
		"is_risky_tld": risky,
		"is_icann":     icann,
	})
}

func CheckUrlShortenerHandler(c *gin.Context) {
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

	isShortener := checks.IsUrlShortener(domain)
	c.JSON(http.StatusOK, gin.H{
		"is_url_shortener": isShortener,
	})
}
