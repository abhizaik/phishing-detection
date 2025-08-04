package handler

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// RootHandler returns basic info about the SafeSurf API service
	r.GET("/", RootHandler)

	// Unversioned global health check
	r.GET("/health", HealthHandler)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", HealthHandler)

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
		v1.GET("/whois", WhoisHandler)

	}

	return r
}
