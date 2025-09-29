package handler

import (
	"log"
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/abhizaik/SafeSurf/internal/service/rank"
	"github.com/gin-gonic/gin"
)

type domainRequest struct {
	Domain string `json:"domain" binding:"required"`
}

func GetDomainRankHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url query param is required"})
		return
	}

	_, isValid, err := checks.IsValidURL(url)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	domain, err := checks.GetDomain(url)
	if err != nil {
		log.Printf("domain extraction failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not extract domain from url"})
		return
	}

	rankValue := rank.DomainRankLookup(domain)

	c.JSON(http.StatusOK, gin.H{
		"rank": rankValue,
	})
}
