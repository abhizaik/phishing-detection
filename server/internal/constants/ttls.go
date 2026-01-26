package constants

import "time"

// Cache TTL constants for analyzer tasks
// These values balance freshness with performance for cached network operations

// DomainRankTTL - Domain rank doesn't change frequently
const DomainRankTTL = 24 * time.Hour

// IPResolutionTTL - DNS IP addresses can change but not frequently
const IPResolutionTTL = 30 * time.Minute

// DNSValidityTTL - DNS records (NS/MX) can change occasionally
const DNSValidityTTL = 30 * time.Minute

// WHOISLookupTTL - Domain registration info is relatively stable
const WHOISLookupTTL = 24 * time.Hour

// HTTPCombinedTTL - HTTP responses (redirects, status codes, HSTS) can change more frequently
const HTTPCombinedTTL = 30 * time.Minute

// TLSCombinedTTL - TLS/SSL certificates don't change frequently
const TLSCombinedTTL = 24 * time.Hour

// PhishTankTTL - Threat intel data should be reasonably fresh
const PhishTankTTL = 1 * time.Hour

// ContentAnalysisTTL - Page content can change, but usually stays consistent for a while
const ContentAnalysisTTL = 1 * time.Hour
