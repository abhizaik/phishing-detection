package handler

import (
	"time"

	"github.com/abhizaik/SafeSurf/internal/handler/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Allow all origins
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Global Rate Limiter: 20 requests per minute per IP
	r.Use(middleware.RateLimiter(20, time.Minute))

	// RootHandler returns basic info about the SafeSurf API service
	r.GET("/", RootHandler)

	// Unversioned global health check
	r.GET("/health", HealthHandler)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", HealthHandler)
		v1.GET("/test", TestApiHandler)

		v1.GET("/analyze", AnalyzeURLHandler)
		v1.GET("/rank", GetDomainRankHandler)
		v1.GET("/ip/check", CheckIfUsesIPHandler)
		v1.GET("/ip/resolve", ResolveIPHandler)
		v1.GET("/length", CheckUrlLengthHandler)
		v1.GET("/depth", CheckUrlDepthHandler)
		v1.GET("/hsts", CheckHSTSHandler)
		v1.GET("/redirects", CheckRedirectsHandler)
		v1.GET("/punycode", CheckPunycodeHandler)
		v1.GET("/trusted-tld", CheckTrustedTLDHandler)
		v1.GET("/risky-tld", CheckRiskyTLDHandler)
		v1.GET("/url-shortener", CheckUrlShortenerHandler)
		v1.GET("/status-code", CheckStatusCodeHandler)
		v1.GET("/domain-info", DomainInfoHandler)

		v1.DELETE("/cache", FlushCacheHandler)

		v1.GET("/screenshot", ScreenshotHandler)

	}

	return r
}
