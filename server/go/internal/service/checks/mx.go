package checks

import "net"

func CheckMXValidity(domain string) (bool, error) {
	mxRecords, err := net.LookupMX(domain)
	mxValid := false
	if err == nil && len(mxRecords) > 0 {
		for _, mx := range mxRecords {
			// Ensure the MX host resolves to at least one IP
			if ips, err := net.LookupIP(mx.Host); err == nil && len(ips) > 0 {
				mxValid = true
				break
			}
		}
	}

	return mxValid, nil
}
