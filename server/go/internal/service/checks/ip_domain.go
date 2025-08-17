package checks

import (
	"net"
	"net/url"
)

func UsesIPInsteadOfDomain(rawURL string) (bool, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false, err
	}

	host := parsedURL.Hostname()
	ip := net.ParseIP(host)
	return ip != nil, nil
}
