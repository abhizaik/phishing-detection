package handler

import (
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/service/cache"
	"github.com/gin-gonic/gin"
)

func FlushCacheHandler(c *gin.Context) {
	// Initialize cache
	cacheInstance, err := cache.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  "failed to connect to cache service",
		})
		return
	}
	defer cacheInstance.Close()

	// Flush all keys
	if err := cacheInstance.FlushAll(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  "failed to flush cache",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": "all cache has been flushed successfully",
	})
}
