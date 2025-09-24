package analyzer

import (
	"fmt"
	"strings"
)

func GenerateResult(resp Response) Result {
	var neutralReasons []string
	var goodReasons []string
	var badReasons []string
	trustScore := 0
	riskScore := 0

	// Rank
	if resp.Features.Rank == 0 {
		badReasons = append(badReasons, "Domain rank is 10,00,000+ (very low reputation)")
		riskScore += 20
	} else if resp.Features.Rank > 0 && resp.Features.Rank <= 10000 {
		goodReasons = append(goodReasons, fmt.Sprintf("Domain rank is #%d globally (high reputation)", resp.Features.Rank))
		trustScore += 20
	} else if resp.Features.Rank > 50000 {
		goodReasons = append(goodReasons, fmt.Sprintf("Domain rank is #%d (medium reputation)", resp.Features.Rank))
		trustScore += 10
	} else {
		neutralReasons = append(neutralReasons, fmt.Sprintf("Domain rank is #%d (low reputation)", resp.Features.Rank))
	}

	// TLD
	if resp.Features.TLD.IsRisky {
		badReasons = append(badReasons, "Domain uses a risky TLD, often misused by attackers")
		riskScore += 20
	}

	if resp.Features.TLD.IsTrusted {
		goodReasons = append(goodReasons, "Domain uses a trusted TLD, verified by registry")
		trustScore += 15
	} else if resp.Features.TLD.IsICANN && !resp.Features.TLD.IsRisky {
		neutralReasons = append(neutralReasons, "Domain uses a normal TLD")
	}

	if !resp.Features.TLD.IsICANN {
		badReasons = append(badReasons, "Domain uses a TLD which is not under ICANN")
		riskScore += 10
	}

	// HSTS
	if resp.Analysis.SupportsHSTS {
		goodReasons = append(goodReasons, "Domain supports HSTS, common in legit sites")
		trustScore += 10
	}

	// URL Shortener
	if resp.Features.URL.IsURLShortener {
		badReasons = append(badReasons, "Domain is a URL shortener, often used to hide real destination")
		riskScore += 15
	}

	// Uses IP
	if resp.Features.URL.UsesIP {
		badReasons = append(badReasons, "IP instead of domain, not used by trusted entities")
		riskScore += 25
	}

	// Punycode
	if resp.Features.URL.ContainsPunycode {
		badReasons = append(badReasons, "URL contains punycode, might spoof original entities")
		riskScore += 20
	}

	// Too deep
	if resp.Features.URL.TooDeep {
		badReasons = append(badReasons, "URL has many slashes (too deep), suspicious")
		riskScore += 10
	}

	// Too long
	if resp.Features.URL.TooLong {
		badReasons = append(badReasons, "URL is too long, may be hiding malicious parts")
		riskScore += 10
	}

	// Subdomain Count
	if resp.Features.URL.SubdomainCount > 2 {
		badReasons = append(badReasons, "URL has many subdomains, may be spoofing")
		riskScore += 15
	}

	// Keywords
	if resp.Features.URL.Keywords.HasKeywords {
		badReasons = append(badReasons, fmt.Sprintf("URL has sensitive keywords like %s", strings.Join(resp.Features.URL.Keywords.Found, ", ")))
		riskScore += 20
	}

	// Nameservers
	if resp.Infrastructure.NameserversValid {
		goodReasons = append(goodReasons, "Nameservers are valid")
		trustScore += 10
	} else {
		badReasons = append(badReasons, "Nameservers are invalid/unverified")
		riskScore += 10
	}

	// MX records
	if resp.Infrastructure.MXRecordsValid {
		goodReasons = append(goodReasons, "MX records are valid (can receive email)")
		trustScore += 10
	} else {
		badReasons = append(badReasons, "MX records are missing/invalid")
		riskScore += 10
	}

	// Domain age
	if resp.DomainInfo != nil {
		if resp.DomainInfo.AgeDays <= 30 {
			badReasons = append(badReasons, fmt.Sprintf("Domain age is %s (very new), high risk", resp.DomainInfo.AgeHuman))
			riskScore += 25
		} else if resp.DomainInfo.AgeDays <= 365 {
			badReasons = append(badReasons, fmt.Sprintf("Domain age is %s (new), medium risk", resp.DomainInfo.AgeHuman))
			riskScore += 15
		} else if resp.DomainInfo.AgeDays <= 1825 {
			neutralReasons = append(neutralReasons, fmt.Sprintf("Domain age is %s (established)", resp.DomainInfo.AgeHuman))
			trustScore += 5
		} else if resp.DomainInfo.AgeDays <= 3650 {
			neutralReasons = append(neutralReasons, fmt.Sprintf("Domain age is %s (mature)", resp.DomainInfo.AgeHuman))
			trustScore += 10
		} else {
			goodReasons = append(goodReasons, fmt.Sprintf("Domain age is %s (very old)", resp.DomainInfo.AgeHuman))
			trustScore += 15
		}

		if resp.DomainInfo.Registrar != "" {
			goodReasons = append(goodReasons, fmt.Sprintf("Registrar is %s", resp.DomainInfo.Registrar))
			trustScore += 5
		}

		if resp.DomainInfo.DNSSEC {
			goodReasons = append(goodReasons, "DNSSEC is enabled (extra security)")
			trustScore += 10
		} else {
			badReasons = append(badReasons, "DNSSEC is not enabled")
			riskScore += 5
		}
	}

	// Redirect chain
	if resp.Analysis.RedirectionResult.IsRedirected {
		if resp.Analysis.RedirectionResult.ChainLength > 3 {
			badReasons = append(badReasons, fmt.Sprintf("Redirect chain is long (%d hops)", resp.Analysis.RedirectionResult.ChainLength))
			riskScore += 10
		} else {
			goodReasons = append(goodReasons, fmt.Sprintf("Redirect chain length is %d (normal)", resp.Analysis.RedirectionResult.ChainLength))
			trustScore += 5
		}

		if resp.Analysis.RedirectionResult.HasDomainJump {
			badReasons = append(badReasons, "Website jumps to a different domain (very risky)")
			riskScore += 20
		}
	}

	// Homoglyph
	if resp.Features.URL.HasHomoglyph {
		badReasons = append(badReasons, "Has homoglyphs, may spoof legit websites")
		riskScore += 20
	}

	// --- Normalize / cap scores ---
	if riskScore > 100 {
		riskScore = 100
	}
	if trustScore > 100 {
		trustScore = 100
	}

	var verdict string
	switch {
	// Very risky: high risk, low trust
	case riskScore >= 70 && trustScore < 40:
		verdict = "Risky"
	// Suspicious: moderate risk OR conflicting signals
	case riskScore >= 40 && riskScore < 70:
		if trustScore >= 40 {
			verdict = "Suspicious (Trust signals present)"
		}
		verdict = "Suspicious"
	// Trustworthy: low risk, high trust
	case riskScore < 40 && trustScore >= 60:
		verdict = "Trustworthy"
	// Unclear / low trust but also low risk
	default:
		verdict = "Unclear"
	}

	combined := int(riskScore - trustScore) // -100..100
	finalScore := (combined + 100) / 2

	res := Result{
		RiskScore:  riskScore,
		TrustScore: trustScore,
		FinalScore: finalScore,
		Verdict:    verdict,
		Reasons: Reasons{
			NeutralReasons: neutralReasons,
			GoodReasons:    goodReasons,
			BadReasons:     badReasons,
		},
	}

	return res
}

