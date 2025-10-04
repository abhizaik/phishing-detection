package checks

import "net"

func CheckNSValidity(domain string) (bool, []string, error) {
	nsRecords, err := net.LookupNS(domain)
	nsValid := false
	var nsHosts []string

	if err == nil && len(nsRecords) > 0 {
		for _, ns := range nsRecords {
			nsHosts = append(nsHosts, ns.Host)
			// Only check if we haven't found a valid NS yet
			if !nsValid {
				// Ensure the nameserver host resolves to at least one IP
				if ips, err := net.LookupIP(ns.Host); err == nil && len(ips) > 0 {
					nsValid = true
					break
				}
			}
		}
	}

	return nsValid, nsHosts, nil
}
