// This handler is for development/testing purposes only.
// This code is intended for local experiments and feature prototyping.
// Do NOT commit permanent changes or production features to this handler.

package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/abhizaik/SafeSurf/internal/service/cache"
	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/gin-gonic/gin"
)

func TestApiHandler(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not extract domain from url"})
		return
	}

	////////////////////////////////////////////////////////////////////
	// THIS IS A TEMPORARY CODE SECTION
	// MAKE CHANGES ONLY INSIDE THIS BOX
	// SAFE TO DELETE LATER IF PUSHED BY MISTAKE
	ctx := context.Background()

	conn, err := cache.New()
if err != nil {
    log.Fatalf("failed to connect to redis: %v", err)
	return 
}


	_ = conn.Set(ctx, "foo", "bar",time.Minute)
	val, _ := conn.Get(ctx, "foo")

	c.JSON(http.StatusOK, gin.H{
		"domain": domain,
		"val": val,
	})

	////////////////////////////////////////////////////////////////////

}
