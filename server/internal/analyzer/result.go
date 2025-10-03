package analyzer

import (
	"fmt"
	"math"
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
		badReasons = append(badReasons, "This website is hardly known, so it may not be very trustworthy.")
		riskScore += 20
	} else if resp.Features.Rank > 0 && resp.Features.Rank <= 10000 {
		goodReasons = append(goodReasons, fmt.Sprintf("This website is very popular worldwide, which is a good sign (#%d globally).", resp.Features.Rank))
		trustScore += 90
	} else if resp.Features.Rank > 50000 {
		goodReasons = append(goodReasons, fmt.Sprintf("This website has moderate popularity (#%d), generally considered okay.", resp.Features.Rank))
		trustScore += 50
	} else {
		neutralReasons = append(neutralReasons, fmt.Sprintf("This website is not very popular (#%d), which is neither good nor bad.", resp.Features.Rank))
	}

	// TLD
	if resp.Features.TLD.IsRisky {
		badReasons = append(badReasons, "The website uses a domain type that is sometimes used by untrustworthy sites.")
		riskScore += 20
	}

	if resp.Features.TLD.IsTrusted {
		goodReasons = append(goodReasons, "The website uses a trusted domain type, verified by official authorities.")
		trustScore += 100
	} else if resp.Features.TLD.IsICANN && !resp.Features.TLD.IsRisky {
		neutralReasons = append(neutralReasons, "The website uses a standard domain type, which is normal.")
	}

	if !resp.Features.TLD.IsICANN {
		badReasons = append(badReasons, "The domain type is uncommon and not officially regulated, which may be risky.")
		riskScore += 30
	}

	// HSTS
	if resp.Analysis.SupportsHSTS {
		goodReasons = append(goodReasons, "The website uses extra security to protect your connection (HSTS).")
		trustScore += 20
	}

	// URL Shortener
	if resp.Features.URL.IsURLShortener {
		badReasons = append(badReasons, "This is a shortened URL, which may hide the real website address.")
		riskScore += 25
	}

	// Uses IP
	if resp.Features.URL.UsesIP {
		badReasons = append(badReasons, "The website uses an IP address instead of a domain name, which is unusual for trusted sites.")
		riskScore += 100
	}

	// Punycode
	if resp.Features.URL.ContainsPunycode {
		badReasons = append(badReasons, "The URL uses special characters to imitate a real website, which may be misleading.")
		riskScore += 100
	}

	// Too deep
	if resp.Features.URL.TooDeep {
		badReasons = append(badReasons, "The URL has many slashes, which can sometimes indicate a suspicious page.")
		riskScore += 30
	}

	// Too long
	if resp.Features.URL.TooLong {
		badReasons = append(badReasons, "The URL is unusually long, which may hide malicious content.")
		riskScore += 20
	}

	// Subdomain Count
	if resp.Features.URL.SubdomainCount > 2 {
		badReasons = append(badReasons, "The website has many subdomains, which can be a trick to mimic trusted sites.")
		riskScore += 15
	}

	// Keywords
	if resp.Features.URL.Keywords.HasKeywords {
		badReasons = append(badReasons, fmt.Sprintf("The URL contains sensitive keywords like %s, which may indicate phishing or scams.", strings.Join(resp.Features.URL.Keywords.Found, ", ")))
		riskScore += 10
	}

	// Nameservers
	if resp.Infrastructure.NameserversValid {
		goodReasons = append(goodReasons, "The website’s server information is properly verified.")
		trustScore += 10
	} else {
		badReasons = append(badReasons, "The website’s server information couldn’t be fully verified.")
		riskScore += 10
	}

	// MX records
	if resp.Infrastructure.MXRecordsValid {
		goodReasons = append(goodReasons, "The website can properly receive emails, which is a good sign.")
		trustScore += 10
	} else {
		badReasons = append(badReasons, "The website may not have proper email setup.")
		riskScore += 10
	}

	// Domain age
	if resp.DomainInfo != nil {
		if resp.DomainInfo.AgeDays <= 30 {
			badReasons = append(badReasons, fmt.Sprintf("The website is very new (%s), which may be risky.", resp.DomainInfo.AgeHuman))
			riskScore += 25
		} else if resp.DomainInfo.AgeDays <= 365 {
			badReasons = append(badReasons, fmt.Sprintf("The website is new (%s), use caution.", resp.DomainInfo.AgeHuman))
			riskScore += 15
		} else if resp.DomainInfo.AgeDays <= 1825 {
			neutralReasons = append(neutralReasons, fmt.Sprintf("The website has been around for a while (%s).", resp.DomainInfo.AgeHuman))
			trustScore += 5
		} else if resp.DomainInfo.AgeDays <= 3650 {
			neutralReasons = append(neutralReasons, fmt.Sprintf("The website is mature (%s), generally trustworthy.", resp.DomainInfo.AgeHuman))
			trustScore += 10
		} else {
			goodReasons = append(goodReasons, fmt.Sprintf("The website is very old (%s), which is a good sign of reliability.", resp.DomainInfo.AgeHuman))
			trustScore += 15
		}

		if resp.DomainInfo.Registrar != "" {
			goodReasons = append(goodReasons, fmt.Sprintf("The website is registered with %s.", resp.DomainInfo.Registrar))
			trustScore += 5
		}

		if resp.DomainInfo.DNSSEC {
			goodReasons = append(goodReasons, "The website has extra security checks enabled (DNSSEC).")
			trustScore += 10
		} else {
			badReasons = append(badReasons, "The website does not have extra security checks enabled (DNSSEC).")
			riskScore += 5
		}
	}

	// Redirect chain
	if resp.Analysis.RedirectionResult.IsRedirected {
		if resp.Analysis.RedirectionResult.ChainLength > 3 {
			badReasons = append(badReasons, fmt.Sprintf("The website redirects multiple times (%d hops), which can be suspicious.", resp.Analysis.RedirectionResult.ChainLength))
			riskScore += 40
		}

		if resp.Analysis.RedirectionResult.HasDomainJump {
			badReasons = append(badReasons, "The website redirects to a different domain, which is risky.")
			badReasons = append(badReasons, fmt.Sprintf("The final website you are sent to is: %s", resp.Analysis.RedirectionResult.FinalURLHost))
			riskScore += 50
		}
	}

	// Homoglyph
	if resp.Features.URL.HasHomoglyph {
		badReasons = append(badReasons, "The website uses letters or characters that look like real ones to trick users, which is risky.")
		riskScore += 60
	}

	// --- Normalize / cap scores ---
	riskScore = clamp(riskScore)
	trustScore = clamp(trustScore)

	// 1
	// combined := int(riskScore - trustScore) // -100..100
	// finalScore := (combined + 100) / 2

	// 2
	// trustContribution := 100 - trustScore
	// finalScore := int(float64(riskScore)*0.7 + float64(trustContribution)*0.3)

	// 3
	finalScore := int(float64(trustScore) - float64(riskScore)*0.2)

	finalScore = clamp(finalScore)
	var verdict string
	switch {
	// Very risky: high risk, low trust
	case finalScore < 50:
		verdict = "Risky"
	// Suspicious: moderate risk OR conflicting signals
	case finalScore < 80:
		verdict = "Suspicious"
	// Safe: low risk, high trust
	case finalScore >= 80 && finalScore <= 100:
		verdict = "Safe"
	// Unclear / low trust but also low risk
	default:
		verdict = "Unclear"
	}

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

func clamp(score int) int {
	return int(math.Max(0, math.Min(100, float64(score))))
}
