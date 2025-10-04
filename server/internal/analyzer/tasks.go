package analyzer

import (
	"github.com/abhizaik/SafeSurf/internal/service/checks"
	"github.com/abhizaik/SafeSurf/internal/service/domaininfo"
	"github.com/abhizaik/SafeSurf/internal/service/rank"
)

// Rank
type rankTask struct{}

func (rankTask) Name() string { return "domain_rank" }
func (rankTask) Run(in *Input, out *Output) error {
	r := rank.DomainRankLookup(in.Domain)
	out.mu.Lock()
	out.Rank = r
	out.mu.Unlock()
	return nil
}

// HSTS
type hstsTask struct{}

func (hstsTask) Name() string { return "hsts_check" }
func (hstsTask) Run(in *Input, out *Output) error {
	h, err := checks.SupportsHSTS(in.URL)
	if err != nil {
		return err
	}
	out.mu.Lock()
	out.SupportsHSTS = h
	out.mu.Unlock()
	return nil
}

// IP checks
type ipCheckTask struct{}

func (ipCheckTask) Name() string { return "ip_check" }
func (ipCheckTask) Run(in *Input, out *Output) error {
	b, err := checks.UsesIPInsteadOfDomain(in.URL)
	if err != nil {
		return err
	}
	out.mu.Lock()
	out.URLUsesIP = b
	out.mu.Unlock()
	return nil
}

type ipResolveTask struct{}

func (ipResolveTask) Name() string { return "ip_resolution" }
func (ipResolveTask) Run(in *Input, out *Output) error {
	ips, err := checks.GetIPAddress(in.Domain)
	if err != nil {
		return err
	}
	out.mu.Lock()
	out.IPs = ips
	out.mu.Unlock()
	return nil
}

// Punycode
type punycodeTask struct{}

func (punycodeTask) Name() string { return "punycode_check" }
func (punycodeTask) Run(in *Input, out *Output) error {
	b, err := checks.ContainsPunycode(in.URL)
	if err != nil {
		return err
	}
	out.mu.Lock()
	out.URLContainsPuny = b
	out.mu.Unlock()
	return nil
}

// Redirects
type redirectsTask struct{}

func (redirectsTask) Name() string { return "redirect_check" }
func (redirectsTask) Run(in *Input, out *Output) error {
	redir, err := checks.CheckRedirects(in.URL)
	if err != nil {
		return err
	}
	out.mu.Lock()
	out.RedirectionResult = redir
	out.mu.Unlock()
	return nil
}

// TLD
type tldTask struct{}

func (tldTask) Name() string { return "tld_check" }
func (tldTask) Run(in *Input, out *Output) error {
	t, icann, tld := checks.IsTrustedTld(in.Domain)
	r, _, _ := checks.IsRiskyTld(in.Domain)
	out.mu.Lock()
	out.TLDTrusted = t
	out.TLDICANN = icann
	out.TLDRisky = r
	out.TLD = tld
	out.mu.Unlock()
	return nil
}

// Shortener
type shortenerTask struct{}

func (shortenerTask) Name() string { return "shortener_check" }
func (shortenerTask) Run(in *Input, out *Output) error {
	s := checks.IsUrlShortener(in.Domain)
	out.mu.Lock()
	out.URLIsShortener = s
	out.mu.Unlock()
	return nil
}

// Status
type statusTask struct{}

func (statusTask) Name() string { return "status_code_check" }
func (statusTask) Run(in *Input, out *Output) error {
	code, text, success, redirect, err := checks.GetStatusCode(in.URL)
	if err != nil {
		return err
	}
	out.mu.Lock()
	out.StatusCode = code
	out.StatusText = text
	out.StatusSuccess = success
	out.StatusIsRedirect = redirect
	out.mu.Unlock()
	return nil
}

// URL structure
type structureTask struct{}

func (structureTask) Name() string { return "url_structure_check" }
func (structureTask) Run(in *Input, out *Output) error {
	out.mu.Lock()
	out.URLTooLong = checks.TooLongUrl(in.URL)
	out.URLTooDeep = checks.TooDeepUrl(in.URL)
	out.mu.Unlock()
	return nil
}

// Keywords
type keywordsTask struct{}

func (keywordsTask) Name() string { return "keywords_check" }
func (keywordsTask) Run(in *Input, out *Output) error {
	present, matches, cats := checks.CheckURLKeywords(in.URL)
	out.mu.Lock()
	out.URLKeywordsPresent = present
	out.URLKeywordMatches = matches
	out.URLKeywordCats = cats
	out.mu.Unlock()
	return nil
}

// DNS validity (NS/MX)
type dnsValidityTask struct{}

func (dnsValidityTask) Name() string { return "dns_validity_check" }
func (dnsValidityTask) Run(in *Input, out *Output) error {
	ns, nsHosts, _ := checks.CheckNSValidity(in.Domain)
	mx, mxHosts, _ := checks.CheckMXValidity(in.Domain)
	out.mu.Lock()
	out.NSValid = ns
	out.NSHosts = nsHosts
	out.MXValid = mx
	out.MXHosts = mxHosts
	out.mu.Unlock()
	return nil
}

// Subdomains
type subdomainTask struct{}

func (subdomainTask) Name() string { return "subdomain_check" }
func (subdomainTask) Run(in *Input, out *Output) error {
	count, _ := checks.GetSubdomainCount(in.URL)
	out.mu.Lock()
	out.URLSubdomainCount = count
	out.mu.Unlock()
	return nil
}

// Domain info
type whoisTask struct{}

func (whoisTask) Name() string { return "whois_lookup" }
func (whoisTask) Run(in *Input, out *Output) error {
	di, err := domaininfo.Lookup(in.Domain)
	if err != nil {
		return err
	}
	out.mu.Lock()
	out.DomainInfo = di
	out.mu.Unlock()
	return nil
}

// SSL info
type sslTask struct{}

func (sslTask) Name() string { return "ssl_check" }
func (sslTask) Run(in *Input, out *Output) error {
	sslInfo := checks.AnalyzeSSLCert(in.Domain)

	out.mu.Lock()
	out.SSLInfo = sslInfo
	out.mu.Unlock()
	return nil
}

// Entropy
type entropyTask struct{}

func (entropyTask) Name() string { return "entropy_check" }
func (entropyTask) Run(in *Input, out *Output) error {
	r := checks.AnalyzeDomainRandomness(in.Domain)

	out.mu.Lock()
	out.DomainRandomness = r
	out.mu.Unlock()
	return nil
}

// Page content
type contentTask struct{}

func (contentTask) Name() string { return "content_check" }
func (contentTask) Run(in *Input, out *Output) error {
	c, _ := checks.GetPageFormInfo(in.Domain)

	out.mu.Lock()
	out.ContentData = c
	out.mu.Unlock()
	return nil
}

// TLS
type tlsTask struct{}

func (tlsTask) Name() string { return "tls_check" }
func (tlsTask) Run(in *Input, out *Output) error {
	t, _ := checks.GetTLSInfo(in.Domain)

	out.mu.Lock()
	out.TLSInfo = t
	out.mu.Unlock()
	return nil
}

// Homoglyph
type homoglyphTask struct{}

func (homoglyphTask) Name() string { return "homoglyph_check" }
func (homoglyphTask) Run(in *Input, out *Output) error {
	h, _ := checks.HasHomoglyphs(in.Domain)

	out.mu.Lock()
	out.HomoglyphPresent = h
	out.mu.Unlock()
	return nil
}
