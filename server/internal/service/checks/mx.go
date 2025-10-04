package checks

import "net"

func CheckMXValidity(domain string) (bool, []string, error) {
	mxRecords, err := net.LookupMX(domain)
	mxValid := false
	var mxHosts []string

	if err == nil && len(mxRecords) > 0 {
		for _, mx := range mxRecords {
			mxHosts = append(mxHosts, mx.Host)
			// Only check if we haven't found a valid MX yet
			if !mxValid {
				// Ensure the MX host resolves to at least one IP
				if ips, err := net.LookupIP(mx.Host); err == nil && len(ips) > 0 {
					mxValid = true
					break
				}
			}
		}
	}

	return mxValid, mxHosts, nil
}
