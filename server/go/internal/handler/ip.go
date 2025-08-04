package handler

import (
	"log"
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/gin-gonic/gin"
)

func CheckIfUsesIPHandler(c *gin.Context) {
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

	isIP, err := checks.UsesIPInsteadOfDomain(url)
	if err != nil {
		log.Printf("IP check failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uses_ip": isIP,
	})
}

func ResolveIPHandler(c *gin.Context) {
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
		log.Printf("Domain extraction failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	ips, err := checks.GetIPAddress(domain)
	if err != nil {
		log.Printf("IP resolution failed for domain %s: %v", domain, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ip_addresses": ips,
	})
}
