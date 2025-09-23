package analyzer

import (
	"fmt"
	"strings"
)

func GenerateResult(resp Response) ([]string, []string, []string) {
	var neutralReasons []string
	var goodReasons []string
	var badReasons []string

	// Rank
	if resp.Features.Rank == 0 {
		badReasons = append(badReasons, fmt.Sprintf("Domain rank is 10,00,000+ (very low reputation)"))
	} else if resp.Features.Rank > 0 && resp.Features.Rank <= 10000 {
		goodReasons = append(goodReasons, fmt.Sprintf("Domain rank is #%d globally (high reputation)", resp.Features.Rank))
	} else if resp.Features.Rank > 50000 {
		goodReasons = append(goodReasons, fmt.Sprintf("Domain rank is #%d (medium reputation)", resp.Features.Rank))
	} else {
		neutralReasons = append(neutralReasons, fmt.Sprintf("Domain rank is #%d (low reputation)", resp.Features.Rank))
	}

	// TLD
	if resp.Features.TLD.IsRisky {
		badReasons = append(badReasons, "Domain uses a risky TLD, which is often  misused by attckers.")
	}

	if resp.Features.TLD.IsTrusted {
		goodReasons = append(goodReasons, "Domain uses a trusted TLD, which belongs to trusted entities with verification")
	} else if resp.Features.TLD.IsICANN && !resp.Features.TLD.IsRisky {
		neutralReasons = append(neutralReasons, "Domain uses a normal TLD")
	}

	if !resp.Features.TLD.IsICANN {
		badReasons = append(badReasons, "Domain uses a TLD which is not under ICANN")
	}

	// HSTS
	if resp.Analysis.SupportsHSTS {
		goodReasons = append(goodReasons, "Domain supports HSTS, common practice among legit entites.")
	}

	// URL Shortener
	if resp.Features.URL.IsURLShortener {
		badReasons = append(badReasons, "Domain is of a URL shortener, might be used to hide the actual URL")
	}

	// Uses IP
	if resp.Features.URL.UsesIP {
		badReasons = append(badReasons, "IP instead of domain, not done by trusted enitities, high risk")
	}

	// Punycode
	if resp.Features.URL.ContainsPunycode {
		badReasons = append(badReasons, "URL contains punycode, might be used to fake original entities, high risk")
	}

	// Too deep
	if resp.Features.URL.TooDeep {
		badReasons = append(badReasons, "URL has many slashes, not usually done by trusted entities, high risk")
	}

	// Too long
	if resp.Features.URL.TooLong {
		badReasons = append(badReasons, "URL is too long, might be used to hide sketchy part, high risk")
	}

	// Subdomain Count
	if resp.Features.URL.SubdomainCount > 2 {
		badReasons = append(badReasons, "URL has many subdomains, might be used to fake original entities, high risk")
	}

	// Keywords
	if resp.Features.URL.Keywords.HasKeywords {
		badReasons = append(badReasons, fmt.Sprintf("URL has sensitive keywords like %s", strings.Join(resp.Features.URL.Keywords.Found, ", ")))
	}

	// Nameservers
	if resp.Infrastructure.NameserversValid {
		goodReasons = append(goodReasons, "Nameservers are valid")
	} else {
		badReasons = append(badReasons, "Nameservers are invalid or could not be verified")
	}

	// MX records
	if resp.Infrastructure.MXRecordsValid {
		goodReasons = append(goodReasons, "MX records are valid (can receive email)")
	} else {
		badReasons = append(badReasons, "MX records are missing/invalid")
	}

	// Domain age
	if resp.DomainInfo != nil {
		if resp.DomainInfo.AgeDays <= 30 { // 1 month
			badReasons = append(badReasons, fmt.Sprintf("Domain age is %s (very new), high risk", resp.DomainInfo.AgeHuman))
		} else if resp.DomainInfo.AgeDays <= 365 { // 1 year
			badReasons = append(badReasons, fmt.Sprintf("Domain age is %s (new), medium risk", resp.DomainInfo.AgeHuman))
		} else if resp.DomainInfo.AgeDays <= 1825 { // 5 years
			badReasons = append(badReasons, fmt.Sprintf("Domain age is %s (established), low risk", resp.DomainInfo.AgeHuman))
		} else if resp.DomainInfo.AgeDays <= 3650 { // 10 years
			neutralReasons = append(neutralReasons, fmt.Sprintf("Domain age is %s (mature), very low risk", resp.DomainInfo.AgeHuman))
		} else { // 10+ years
			goodReasons = append(goodReasons, fmt.Sprintf("Domain age is %s (very old), minimal risk", resp.DomainInfo.AgeHuman))
		}
		// Registrar
		if resp.DomainInfo.Registrar != "" {
			goodReasons = append(goodReasons, fmt.Sprintf("Registrar is %s", resp.DomainInfo.Registrar))
		}
		// DNSSEC
		if resp.DomainInfo.DNSSEC {
			goodReasons = append(goodReasons, "DNSSEC is enabled (extra security)")
		} else {
			badReasons = append(badReasons, "DNSSEC is not enabled")
		}
	}

	// Redirect chain
	if resp.Analysis.RedirectionResult.IsRedirected {
		if resp.Analysis.RedirectionResult.ChainLength > 3 {
			badReasons = append(badReasons, fmt.Sprintf("Redirect chain is long (%d hops)", resp.Analysis.RedirectionResult.ChainLength))
		} else {
			goodReasons = append(goodReasons, fmt.Sprintf("Redirect chain length is %d → normal", resp.Analysis.RedirectionResult.ChainLength))
		}

		if resp.Analysis.RedirectionResult.HasDomainJump {
			badReasons = append(badReasons, fmt.Sprintf("Website jumps to different domain than the original one, very risky."))
		}
	}

	// Homoglyph
	if resp.Features.URL.HasHomoglyph {
		badReasons = append(badReasons, fmt.Sprintf("Has homoglyphs, special characters used to spoof legit websites, very risky."))
	}

	return neutralReasons, goodReasons, badReasons
}
