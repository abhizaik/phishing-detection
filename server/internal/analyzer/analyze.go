package analyzer

import (
	"context"
	"log"
	"time"

	"github.com/abhizaik/SafeSurf/internal/service/cache"
	"github.com/abhizaik/SafeSurf/internal/service/checks"
)

// Analyze runs all tasks and builds the final response
func Analyze(ctx context.Context, rawURL string) (Response, []error) {
	// Validate and extract domain
	_, isValid, _ := checks.IsValidURL(rawURL)
	if !isValid {
		return Response{}, []error{ErrInvalidURL}
	}
	domain, err := checks.GetDomain(rawURL)
	if err != nil {
		return Response{}, []error{err}
	}

	// Initialize cache (non-blocking - if cache fails, continue without it)
	var cacheInstance CacheInterface
	cacheConn, err := cache.New()
	if err != nil {
		log.Printf("Warning: Failed to initialize cache: %v. Continuing without cache.", err)
	} else {
		cacheInstance = cacheConn
		defer cacheConn.Close()
	}

	in := &Input{URL: rawURL, Domain: domain, Cache: cacheInstance}

	tasks := []Task{
		rankTask{},
		httpCombinedTask{}, // Optimized: combines redirects, HSTS, and status code
		ipCheckTask{},
		ipResolveTask{},
		punycodeTask{},
		tldTask{},
		shortenerTask{},
		structureTask{},
		keywordsTask{},
		dnsValidityTask{},
		subdomainTask{},
		whoisTask{},
		tlsCombinedTask{}, // Optimized: combines TLS and SSL checks
		entropyTask{},
		contentTask{},
		homoglyphTask{},
	}

	// Start timing right before tasks run
	start := time.Now()
	out, errs := runTasks(ctx, in, tasks)

	resp := Response{
		URL:    rawURL,
		Domain: domain,
		Features: Features{
			Rank: out.Rank,
			TLD: TLDInfo{
				TLD:       out.TLD,
				IsTrusted: out.TLDTrusted,
				IsRisky:   out.TLDRisky,
				IsICANN:   out.TLDICANN,
			},
			URL: URLChecks{
				IsURLShortener:   out.URLIsShortener,
				UsesIP:           out.URLUsesIP,
				ContainsPunycode: out.URLContainsPuny,
				TooLong:          out.URLTooLong,
				TooDeep:          out.URLTooDeep,
				SubdomainCount:   out.URLSubdomainCount,
				HasHomoglyph:     out.HomoglyphPresent,
				Keywords: Keywords{
					HasKeywords: out.URLKeywordsPresent,
					Found:       out.URLKeywordMatches,
					Categories:  out.URLKeywordCats,
				},
			},
		},
		Infrastructure: Infrastructure{
			IPAddresses:      out.IPs,
			NameserversValid: out.NSValid,
			NSHosts:          out.NSHosts,
			MXRecordsValid:   out.MXValid,
			MXHosts:          out.MXHosts,
		},
		DomainInfo: out.DomainInfo,
		Analysis: Analysis{
			RedirectionResult: out.RedirectionResult,
			SupportsHSTS:      out.SupportsHSTS,
			HTTPStatus: HTTPStatus{
				Code:                 out.StatusCode,
				Text:                 out.StatusText,
				Success:              out.StatusSuccess,
				IsRedirectStatusCode: out.StatusIsRedirect,
			},
		},
		Performance: Performance{
			TotalTime: time.Since(start).String(),
			Timings:   ConvertTimings(out.Timings),
		},
	}

	result := GenerateResult(resp)
	resp.Result = result

	if len(errs) > 0 {
		resp.Incomplete = true
		for _, e := range errs {
			resp.Errors = append(resp.Errors, e.Error())
		}
	}

	return resp, errs
}
