package checks

import (
	"net"
	"net/url"
)

func GetIPAddress(domain string) ([]string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, ip := range ips {
		result = append(result, ip.String())
	}
	return result, nil
}

func UsesIPInsteadOfDomain(rawURL string) (bool, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false, err
	}

	host := parsedURL.Hostname()
	ip := net.ParseIP(host)
	return ip != nil, nil
}
