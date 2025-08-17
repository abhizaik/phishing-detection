package checks

import "github.com/abhizaik/SafeSurf/internal/constants"

func IsTrustedTld(domain string) (bool, bool) {
	tld, icann := GetTld(domain)
	_, ok := constants.TrustedTLDs[tld]

	return ok, icann
}
