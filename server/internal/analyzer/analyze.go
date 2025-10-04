package analyzer

import (
	"context"
	"time"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
)

// Analyze runs all tasks and builds the final response
func Analyze(ctx context.Context, rawURL string) (Response, []error) {
	start := time.Now()

	// Validate and extract domain
	_, isValid, _ := checks.IsValidURL(rawURL)
	if !isValid {
		return Response{}, []error{ErrInvalidURL}
	}
	domain, err := checks.GetDomain(rawURL)
	if err != nil {
		return Response{}, []error{err}
	}

	in := &Input{URL: rawURL, Domain: domain}

	tasks := []Task{
		rankTask{},
		hstsTask{},
		ipCheckTask{},
		ipResolveTask{},
		punycodeTask{},
		redirectsTask{},
		tldTask{},
		shortenerTask{},
		statusTask{},
		structureTask{},
		keywordsTask{},
		dnsValidityTask{},
		subdomainTask{},
		whoisTask{},
		sslTask{},
		entropyTask{},
		tlsTask{},
		contentTask{},
		homoglyphTask{},
	}

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
