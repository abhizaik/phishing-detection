package handler

import (
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/abhizaik/SafeSurf/internal/service/domaininfo"
	"github.com/gin-gonic/gin"
)

func DomainInfoHandler(c *gin.Context) {
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

	domainInfo, err := domaininfo.Lookup(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "domainInfo lookup failed",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domain":     domain,
		"domainInfo": domainInfo,
	})
}
