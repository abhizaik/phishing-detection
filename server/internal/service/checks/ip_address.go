package checks

import (
	"net"
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
