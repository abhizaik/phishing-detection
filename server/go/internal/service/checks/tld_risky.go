package checks

import "github.com/abhizaik/SafeSurf/internal/constants"

func IsRiskyTld(domain string) (bool, bool) {
	tld, icann := GetTld(domain)
	_, ok := constants.RiskyTLDs[tld]

	return ok, icann
}
