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

	// --- 1. Popularity & Rank ---
	if resp.Features.Rank == 0 {
		// Changed from "Hardly known" to "Unranked"
		badReasons = append(badReasons, "Very low traffic volume.")
		riskScore += 20
	} else if resp.Features.Rank > 0 && resp.Features.Rank <= 10000 {
		goodReasons = append(goodReasons, fmt.Sprintf("Global Giant: Ranked #%d worldwide.", resp.Features.Rank))
		trustScore += 90
	} else if resp.Features.Rank > 50000 {
		goodReasons = append(goodReasons, fmt.Sprintf("Established website with moderate popularity (#%d).", resp.Features.Rank))
		trustScore += 50
	} else {
		goodReasons = append(goodReasons, fmt.Sprintf("Niche website with standard traffic volume (#%d).", resp.Features.Rank))
		trustScore += 20
	}

	// --- 2. TLD (Top Level Domain) ---
	if resp.Features.TLD.IsRisky {
		badReasons = append(badReasons, "High-risk domain extension detected (often associated with spam).")
		riskScore += 20
	}

	if resp.Features.TLD.IsTrusted {
		goodReasons = append(goodReasons, "High-trust official domain extension (Gov/Edu).")
		trustScore += 100
	} else if resp.Features.TLD.IsICANN && !resp.Features.TLD.IsRisky {
		// Keep it simple
		neutralReasons = append(neutralReasons, "Standard, officially recognized domain extension.")
	}

	if !resp.Features.TLD.IsICANN {
		badReasons = append(badReasons, "Unregulated or non-standard domain extension.")
		riskScore += 30
	}

	// --- 3. Security Protocols ---
	if resp.Analysis.SupportsHSTS {
		goodReasons = append(goodReasons, "Enforces strict HTTPS security (HSTS Enabled).")
		trustScore += 20
	}

	// --- 4. URL Structure / Obfuscation ---
	if resp.Features.URL.IsURLShortener {
		badReasons = append(badReasons, "URL Shortener detected (hides the true destination).")
		riskScore += 25
	}

	// Uses IP
	if resp.Features.URL.UsesIP {
		badReasons = append(badReasons, "Raw IP address usage detected (common evasion tactic).")
		riskScore += 100
	}

	// Punycode
	if resp.Features.URL.ContainsPunycode {
		badReasons = append(badReasons, "Punycode characters detected (potential phishing spoof).")
		riskScore += 100
	}

	// Too deep
	if resp.Features.URL.TooDeep {
		badReasons = append(badReasons, "Excessively deep URL path (potential request hiding).")
		riskScore += 30
	}

	// Too long
	if resp.Features.URL.TooLong {
		badReasons = append(badReasons, "URL length exceeds standard limits (potential buffer overflow/hiding).")
		riskScore += 20
	}

	// Subdomain Count
	if resp.Features.URL.SubdomainCount > 2 {
		badReasons = append(badReasons, "Suspicious number of subdomains detected.")
		riskScore += 15
	}

	// Keywords
	if resp.Features.URL.Keywords.HasKeywords {
		badReasons = append(badReasons, fmt.Sprintf("Sensitive security keywords found in URL: %s", strings.Join(resp.Features.URL.Keywords.Found, ", ")))
		riskScore += 10
	}

	// --- 5. Infrastructure Forensics ---
	if resp.Infrastructure.NameserversValid {
		goodReasons = append(goodReasons, "Valid DNS configuration detected.")
		trustScore += 10
	} else {
		badReasons = append(badReasons, "Incomplete or missing DNS configuration.")
		riskScore += 10
	}

	// MX records
	if resp.Infrastructure.MXRecordsValid {
		goodReasons = append(goodReasons, "Valid email server configuration (MX Records).")
		trustScore += 10
	} else {
		neutralReasons = append(neutralReasons, "No email server configured for this domain.")
		// Reduced risk score here, as some landing pages legitimately don't have email
		riskScore += 5
	}

	// --- 6. Domain History ---
	if resp.DomainInfo != nil {
		if resp.DomainInfo.AgeDays <= 30 {
			badReasons = append(badReasons, fmt.Sprintf("Newly created domain (%s old). High Risk.", resp.DomainInfo.AgeHuman))
			riskScore += 25
		} else if resp.DomainInfo.AgeDays <= 365 {
			badReasons = append(badReasons, fmt.Sprintf("Young domain (%s old). Use caution.", resp.DomainInfo.AgeHuman))
			riskScore += 15
		} else if resp.DomainInfo.AgeDays <= 1825 {
			neutralReasons = append(neutralReasons, fmt.Sprintf("Operational for %s.", resp.DomainInfo.AgeHuman))
			trustScore += 5
		} else {
			goodReasons = append(goodReasons, fmt.Sprintf("Long-standing domain history (%s).", resp.DomainInfo.AgeHuman))
			trustScore += 15
		}

		if resp.DomainInfo.Registrar != "" {
			goodReasons = append(goodReasons, fmt.Sprintf("Registered with %s", resp.DomainInfo.Registrar))
			trustScore += 5
		}

		// DNSSEC Logic Updated
		if resp.DomainInfo.DNSSEC {
			goodReasons = append(goodReasons, "Advanced DNS security enabled (DNSSEC).")
			trustScore += 10
		} else {
			// Moved to Neutral. Not having DNSSEC is NOT a sign of phishing for .coms
			neutralReasons = append(neutralReasons, "Standard DNS security (DNSSEC not enabled).")
			// Removed riskScore penalty
		}
	}

	// --- 7. Redirection Analysis ---
	if resp.Analysis.RedirectionResult.IsRedirected {
		if resp.Analysis.RedirectionResult.ChainLength > 3 {
			badReasons = append(badReasons, fmt.Sprintf("Excessive redirection chain detected (%d hops).", resp.Analysis.RedirectionResult.ChainLength))
			riskScore += 40
		}

		if resp.Analysis.RedirectionResult.HasDomainJump {
			badReasons = append(badReasons, "Cross-domain redirection detected (destination differs from source).")
			// Add the destination as a neutral fact so they can see where they are going
			badReasons = append(badReasons, fmt.Sprintf("Final Destination: %s", resp.Analysis.RedirectionResult.FinalURLHost))
			riskScore += 50
		}
	}

	// --- 8. Homoglyphs ---
	if resp.Features.URL.HasHomoglyph {
		badReasons = append(badReasons, "Homoglyph attack detected (deceptive visual characters).")
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
