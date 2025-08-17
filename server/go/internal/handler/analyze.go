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

type Response struct {
	URL            string                       `json:"url"`
	Domain         string                       `json:"domain"`
	Features       Features                     `json:"features"`
	Infrastructure Infrastructure               `json:"infrastructure"`
	DomainInfo     *domaininfo.RegistrationData `json:"domain_info"`
	Analysis       Analysis                     `json:"analysis"`
	Performance    Performance                  `json:"performance"`
	Incomplete     bool                         `json:"incomplete"`
	Errors         []string                     `json:"errors"`
}

// all phishing-related derived features
type Features struct {
	Rank int       `json:"rank"`
	TLD  TLDInfo   `json:"tld"`
	URL  URLChecks `json:"url"`
}

// keywords reflated stuff
type Keywords struct {
	HasKeywords bool                `json:"has_keywords"`
	Found       []string            `json:"found"`
	Categories  map[string][]string `json:"categories"`
}

// infra-level stuff
type Infrastructure struct {
	IPAddresses      []string `json:"ip_addresses"`
	NameserversValid bool     `json:"nameservers_valid"`
	MXRecordsValid   bool     `json:"mx_records_valid"`
}

// extra artifacts
type Analysis struct {
	IsRedirected           bool       `json:"is_redirected"`
	RedirectionChain       []string   `json:"redirection_chain"`
	RedirectionChainLength int        `json:"redirection_chain_length"`
	RedirectionFinalURL    string     `json:"redirection_final_url"`
	HTTPStatus             HTTPStatus `json:"http_status"`
	SupportsHSTS           bool       `json:"is_hsts_supported"`
}

type TLDInfo struct {
	IsTrusted bool `json:"is_trusted_tld"`
	IsRisky   bool `json:"is_risky_tld"`
	IsICANN   bool `json:"is_icann"`
}

type URLChecks struct {
	IsURLShortener   bool     `json:"url_shortener"`
	UsesIP           bool     `json:"uses_ip"`
	ContainsPunycode bool     `json:"contains_punycode"`
	TooLong          bool     `json:"too_long"`
	TooDeep          bool     `json:"too_deep"`
	SubdomainCount   int      `json:"subdomain_count"`
	Keywords         Keywords `json:"keywords"`
}

// http status code related stuffs
type HTTPStatus struct {
	Code                 int    `json:"code"`
	Text                 string `json:"text"`
	Success              bool   `json:"success"`
	IsRedirectStatusCode bool   `json:"is_redirect"`
}

// performance timings
type Performance struct {
	TotalTime string            `json:"total_time"`
	Timings   map[string]string `json:"timings"`
}

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
		rankValue            int
		supportsHSTS         bool
		isIP                 bool
		ips                  []string
		hasPunycode          bool
		IsRedirected         bool
		statusCodeIsRedirect bool
		chain                []string
		trusted              bool
		risky                bool
		trusted_tld_icann    bool
		// risky_tld_icann      bool
		isShortener          bool
		statusCode           int
		statusText           string
		isSuccess            bool
		tooLong              bool
		tooDeep              bool
		keywordPresent       bool
		keywordMatches       []string
		keywordCategories    map[string][]string
		nameserversReachable bool
		mxRecordsConfigured  bool
		subdomainCount       int
		finalURL             string
		chainLength          int
		// screenshotFileName string
		// screenshotURL      string
		domainInfo *domaininfo.RegistrationData
	)

	// ---------- Concurrent Checks ----------

	// Check: Domain Rank
	wg.Add(1)
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
	wg.Add(1)
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
	wg.Add(1)
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
	wg.Add(1)
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
	wg.Add(1)
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
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			isRedir, redirChain, finURL, chainLen, err := checks.CheckRedirects(url)
			if err != nil {
				errors <- "redirect check failed"
				return
			}
			mu.Lock()
			IsRedirected = isRedir
			chain = redirChain
			finalURL = finURL
			chainLength = chainLen
			timings["redirect_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: Trusted and Risky TLD
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			t, ict := checks.IsTrustedTld(domain)
			r, _ := checks.IsRiskyTld(domain)
			mu.Lock()
			trusted = t
			trusted_tld_icann = ict
			risky = r
			// risky_tld_icann = icr
			timings["tld_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: Is URL Shortener
	wg.Add(1)
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
	wg.Add(1)
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
			statusCodeIsRedirect = redirect
			timings["status_code_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: Keywords in URL
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			mu.Lock()
			keywordPresent, keywordMatches, keywordCategories = checks.CheckURLKeywords(url)
			timings["keywords_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: NS, MX Config
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			mu.Lock()
			nameserversReachable, err = checks.CheckNSValidity(url)
			mxRecordsConfigured, err = checks.CheckMXValidity(url)
			timings["keywords_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: Subdomain
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			mu.Lock()
			subdomainCount, err = checks.GetSubdomainCount(url)
			timings["subdomain_check"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// Check: URL Structure - Length & Depth
	wg.Add(1)
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

	// Domain Info
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		select {
		case <-ctx.Done():
			return
		default:
			d, err := domaininfo.Lookup(domain)
			if err != nil {
				errors <- "whois lookup failed"
				return
			}
			mu.Lock()
			domainInfo = d
			timings["whois_lookup"] = time.Since(start).String()
			mu.Unlock()
		}
	}()

	// // Analysis: Screenshot
	// go func() {
	// 	defer wg.Done()
	// 	start := time.Now()
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	default:
	// 		fn := screenshot.TakeScreenshot(url)
	// 		mu.Lock()
	// 		screenshotFileName = fn
	// 		screenshotURL = fmt.Sprintf("localhost:8080/%s", screenshotFileName)
	// 		timings["screenshot_check"] = time.Since(start).String()
	// 		mu.Unlock()
	// 	}
	// }()
	// // fmt.Sprintf("/screenshots/%s.png", domain)

	// Wait for all concurrent checks (including WHOIS)
	wg.Wait()

	// Calculate total execution time
	totalTime := time.Since(startTime)

	response := Response{
		URL:    url,
		Domain: domain,
		Features: Features{
			Rank: rankValue,
			TLD: TLDInfo{
				IsTrusted: trusted,
				IsRisky:   risky,
				IsICANN:   trusted_tld_icann,
			},
			URL: URLChecks{
				IsURLShortener:   isShortener,
				UsesIP:           isIP,
				ContainsPunycode: hasPunycode,
				TooLong:          tooLong,
				TooDeep:          tooDeep,
				SubdomainCount:   subdomainCount,
				Keywords: Keywords{
					HasKeywords: keywordPresent,
					Found:       keywordMatches,
					Categories:  keywordCategories,
				},
			},
		},
		Infrastructure: Infrastructure{
			IPAddresses:      ips,
			NameserversValid: nameserversReachable,
			MXRecordsValid:   mxRecordsConfigured,
		},
		DomainInfo: domainInfo,
		Analysis: Analysis{
			IsRedirected:           IsRedirected,
			RedirectionChain:       chain,
			RedirectionChainLength: chainLength,
			RedirectionFinalURL:    finalURL,
			SupportsHSTS:           supportsHSTS,
			HTTPStatus: HTTPStatus{
				Code:                 statusCode,
				Text:                 statusText,
				Success:              isSuccess,
				IsRedirectStatusCode: statusCodeIsRedirect,
			},
		},
		Performance: Performance{
			TotalTime: totalTime.String(),
			Timings:   timings,
		},
	}

	// Flag incomplete result if any error occurred or timeout triggered
	if ctx.Err() != nil || len(errors) > 0 {
		response.Incomplete = true
		errList := []string{}
		for len(errors) > 0 {
			errList = append(errList, <-errors)
		}
		response.Errors = errList
	}

	c.JSON(http.StatusOK, response)
}
