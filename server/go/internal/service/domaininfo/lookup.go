package domaininfo

import (
	"github.com/abhizaik/SafeSurf/internal/service/checks"
)

// Lookup tries RDAP first, falls back to WHOIS if RDAP fails.
// Returns: (domaininfo, age as string, error)
func Lookup(domain string) (*RegistrationData, string, error) {
	// Try RDAP first
	rdapData, err := fetchRDAP(domain)
	if err == nil && rdapData != nil {
		age, err := checks.GetDomainAge(rdapData.CreatedDate)
		if err != nil {
			return rdapData, "Age calculation failed (RDAP)", err
		}
		return rdapData, age, nil
	}

	// RDAP failed, fall back to WHOIS
	whoisData, err := GetWhoisData(domain)
	if err != nil {
		return nil, "", err
	}

	age, err := checks.GetDomainAge(whoisData.CreatedDate)
	if err != nil {
		return whoisData, "Age calculation failed (WHOIS)", err
	}

	return whoisData, age, nil
}
