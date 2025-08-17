package handler

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/abhizaik/SafeSurf/internal/service/domaininfo"
	"github.com/abhizaik/SafeSurf/internal/service/rank"
	"github.com/gin-gonic/gin"
)

func AnalyzeURLHandler(c *gin.Context) {
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

	// Record start time for accurate total execution measurement
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex

	errors := make(chan string, 10)
	timings := make(map[string]string)

	var (
		rankValue         int
		supportsHSTS      bool
		isIP              bool
		ips               []string
		hasPunycode       bool
		redirected        bool
		chain             []string
		trusted           bool
		risky             bool
		trusted_tld_icann bool
		risky_tld_icann   bool
		isShortener       bool
		statusCode        int
		statusText        string
		isSuccess         bool
		isRedirect        bool
		tooLong           bool
		tooDeep           bool
		age               string
		domainInfo        *domaininfo.RegistrationData
	)

	// ---------- Concurrent Checks ----------

	wg.Add(11)

	// Check: Domain Rank
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			r := rank.DomainRankLookup(domain)
			mu.Lock()
			rankValue = r
			timings["domain_rank"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: HSTS Support
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			hsts, err := checks.SupportsHSTS(url)
			if err != nil {
				errors <- "hsts check failed"
				return
			}
			mu.Lock()
			supportsHSTS = hsts
			timings["hsts_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: Uses IP Instead of Domain
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			ipUsed, err := checks.UsesIPInsteadOfDomain(url)
			if err != nil {
				errors <- "uses_ip check failed"
				return
			}
			mu.Lock()
			isIP = ipUsed
			timings["ip_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: Get IP addresses
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			resolvedIPs, err := checks.GetIPAddress(domain)
			if err != nil {
				errors <- "ip resolution failed"
				return
			}
			mu.Lock()
			ips = resolvedIPs
			timings["ip_resolution"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: Contains Punycode
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			puny, err := checks.ContainsPunycode(url)
			if err != nil {
				errors <- "punycode check failed"
				return
			}
			mu.Lock()
			hasPunycode = puny
			timings["punycode_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: Redirect Chain
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			isRedir, redirChain, err := checks.CheckRedirects(url)
			if err != nil {
				errors <- "redirect check failed"
				return
			}
			mu.Lock()
			redirected = isRedir
			chain = redirChain
			timings["redirect_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: Trusted and Risky TLD
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			t, ict := checks.IsTrustedTld(domain)
			r, icr := checks.IsRiskyTld(domain)
			mu.Lock()
			trusted = t
			trusted_tld_icann = ict
			risky = r
			risky_tld_icann = icr
			timings["tld_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: Is URL Shortener
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			short := checks.IsUrlShortener(domain)
			mu.Lock()
			isShortener = short
			timings["shortener_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: HTTP Status Code
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			code, text, success, redirect, err := checks.GetStatusCode(url)
			if err != nil {
				errors <- "status code check failed"
				return
			}
			mu.Lock()
			statusCode = code
			statusText = text
			isSuccess = success
			isRedirect = redirect
			timings["status_code_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: URL Structure - Length & Depth
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			mu.Lock()
			tooLong = checks.TooLongUrl(url)
			tooDeep = checks.TooDeepUrl(url)
			timings["url_structure_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// ---------- WHOIS Call (not concurrent-safe or fast) ----------
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			d, a, err := domaininfo.Lookup(domain)
			if err != nil {
				errors <- "whois lookup failed"
				return
			}
			mu.Lock()
			domainInfo = d
			age = a
			timings["whois_lookup"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Wait for all concurrent checks (including WHOIS)
	wg.Wait()

	// Calculate total execution time
	totalTime := time.Since(startTime)

	response := gin.H{
		"rank":                 rankValue,
		"supports_hsts":        supportsHSTS,
		"uses_ip":              isIP,
		"ip_addresses":         ips,
		"contains_punycode":    hasPunycode,
		"is_redirected":        redirected,
		"chain":                chain,
		"is_trusted_tld":       trusted,
		"trusted_tld_is_icann": trusted_tld_icann,
		"is_risky_tld":         risky,
		"risky_tld_is_icann":   risky_tld_icann,
		"is_url_shortener":     isShortener,
		"status_code":          statusCode,
		"status_text":          statusText,
		"is_success":           isSuccess,
		"is_redirect":          isRedirect,
		"too_long":             tooLong,
		"too_deep":             tooDeep,
		"domain":               domain,
		"age":                  age,
		"domainInfo":           domainInfo,
		"performance": gin.H{
			"total_time": totalTime.String(),
			"timings":    timings,
		},
	}

	// Flag incomplete result if any error occurred or timeout triggered
	if ctx.Err() != nil || len(errors) > 0 {
		response["incomplete"] = true
		errList := []string{}
		for len(errors) > 0 {
			errList = append(errList, <-errors)
		}
		response["errors"] = errList
	}

	c.JSON(http.StatusOK, response)
}
