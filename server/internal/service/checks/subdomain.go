package checks

import (
	"fmt"
	"net"
	"strings"

	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

func GetSubdomainCount(rawURL string) (int, error) {

	host, err := GetHost(rawURL)
	if err != nil {
		return 0, err
	}
	if host == "" {
		return 0, fmt.Errorf("empty host")
	}

	host = strings.TrimSuffix(host, ".")

	// If host is an IP address, there are no subdomains
	if net.ParseIP(host) != nil {
		return 0, nil
	}

	// If no dot present (e.g., "localhost")
	if !strings.Contains(host, ".") {
		return 0, nil
	}

	// Normalize IDN -> ASCII (punycode). This is important for international domains.
	asciiHost, err := idna.ToASCII(host)
	if err != nil {
		return 0, err
	}
	asciiHost = strings.ToLower(asciiHost)

	// Get the registrable domain (eTLD+1). This handles multi-part TLDs like "co.uk".
	etldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(asciiHost)
	if err != nil {
		// If publicsuffix can't determine eTLD+1, return an error as result is ambiguous.
		return 0, err
	}

	// If host == etldPlusOne then there are no subdomains
	if asciiHost == etldPlusOne {
		return 0, nil
	}

	// Remove the registrable domain part to isolate subdomain labels.
	subdomainPart := strings.TrimSuffix(asciiHost, "."+etldPlusOne)
	if subdomainPart == "" {
		return 0, nil
	}

	labels := strings.Split(subdomainPart, ".")
	return len(labels), nil
}
