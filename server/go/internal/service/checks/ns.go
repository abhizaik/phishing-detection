package checks

import "net"

func CheckNSValidity(domain string) (bool, error) {
	nsRecords, err := net.LookupNS(domain)
	nsValid := false
	if err == nil && len(nsRecords) > 0 {
		for _, ns := range nsRecords {
			// Ensure the nameserver host resolves to at least one IP
			if ips, err := net.LookupIP(ns.Host); err == nil && len(ips) > 0 {
				nsValid = true
				break
			}
		}
	}

	return nsValid, nil
}
