package checks

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

// CombinedTLSResult contains both TLS and SSL certificate information
type CombinedTLSResult struct {
	// TLS info
	TLSInfo TLSResult
	// SSL cert info
	SSLInfo SSLCertResult
}

// CheckTLSCombined performs a single TLS connection to extract both TLS info
// and SSL certificate details
func CheckTLSCombined(domain string) (CombinedTLSResult, error) {
	result := CombinedTLSResult{}
	result.SSLInfo.Domain = domain

	// Add port if not included
	hostPort := domain
	if !strings.Contains(domain, ":") {
		hostPort = net.JoinHostPort(domain, "443")
	}

	// Create dialer with timeout
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	// Single TLS connection
	conn, err := tls.DialWithDialer(dialer, "tcp", hostPort, &tls.Config{
		InsecureSkipVerify: true, // we validate chain manually
	})
	if err != nil {
		result.SSLInfo.HasTLS = false
		result.SSLInfo.Reasons = append(result.SSLInfo.Reasons, fmt.Sprintf("TLS connection failed: %v", err))
		return result, fmt.Errorf("TLS connection failed: %v", err)
	}
	defer conn.Close()

	// Get connection state
	state := conn.ConnectionState()
	certs := state.PeerCertificates

	if len(certs) == 0 {
		result.SSLInfo.Reasons = append(result.SSLInfo.Reasons, "no peer certificates")
		return result, nil
	}

	cert := certs[0]

	// Extract TLS info
	var issuer string
	if len(cert.Issuer.Organization) > 0 {
		issuer = cert.Issuer.Organization[0]
	} else {
		issuer = cert.Issuer.CommonName
	}

	ageDays := int(time.Since(cert.NotBefore).Hours() / 24)

	var hostnameMismatch bool
	if err := cert.VerifyHostname(domain); err != nil {
		hostnameMismatch = true
	}

	result.TLSInfo = TLSResult{
		Present:          true,
		Issuer:           issuer,
		AgeDays:          ageDays,
		HostnameMismatch: hostnameMismatch,
	}

	// Extract SSL cert info
	result.SSLInfo.HasTLS = true
	result.SSLInfo.Issuer = cert.Issuer.CommonName
	result.SSLInfo.NotBefore = cert.NotBefore
	result.SSLInfo.NotAfter = cert.NotAfter
	result.SSLInfo.AgeDays = ageDays

	// Fingerprint
	fp := sha256.Sum256(cert.Raw)
	result.SSLInfo.Fingerprint = strings.ToUpper(hex.EncodeToString(fp[:]))

	// Validate chain using system roots
	opts := x509.VerifyOptions{
		Roots:         nil, // use system roots
		Intermediates: x509.NewCertPool(),
		DNSName:       domain,
	}
	for _, ic := range certs[1:] {
		opts.Intermediates.AddCert(ic)
	}
	if _, err := cert.Verify(opts); err != nil {
		result.SSLInfo.ChainValid = false
		result.SSLInfo.Reasons = append(result.SSLInfo.Reasons, fmt.Sprintf("cert chain invalid: %v", err))
	} else {
		result.SSLInfo.ChainValid = true
	}

	// Check blacklist
	if _, bad := knownBadFingerprints[result.SSLInfo.Fingerprint]; bad {
		result.SSLInfo.IsSuspicious = true
		result.SSLInfo.KnownBadChain = true
		result.SSLInfo.Reasons = append(result.SSLInfo.Reasons, "certificate fingerprint is blacklisted")
	}

	// Expiry checks
	if time.Now().After(cert.NotAfter) {
		result.SSLInfo.IsSuspicious = true
		result.SSLInfo.Reasons = append(result.SSLInfo.Reasons, "certificate expired")
	}
	if time.Now().Before(cert.NotBefore) {
		result.SSLInfo.IsSuspicious = true
		result.SSLInfo.Reasons = append(result.SSLInfo.Reasons, "certificate not yet valid")
	}

	// Weak validity period (e.g., > 398 days is discouraged by CABF)
	if cert.NotAfter.Sub(cert.NotBefore).Hours()/24 > 398 {
		result.SSLInfo.Reasons = append(result.SSLInfo.Reasons, "unusually long validity period")
	}

	// Check for embedded Certificate Transparency (CT) SCTs
	// OID 1.3.6.1.4.1.11129.2.4.2 is for embedded SCTs
	result.SSLInfo.CTLogged = false
	for _, ext := range cert.Extensions {
		if ext.Id.String() == "1.3.6.1.4.1.11129.2.4.2" {
			result.SSLInfo.CTLogged = true
			break
		}
	}

	if !result.SSLInfo.CTLogged {
		result.SSLInfo.IsSuspicious = true
		result.SSLInfo.Reasons = append(result.SSLInfo.Reasons, "certificate does not contain embedded CT logs (SCTs)")
	}

	return result, nil
}
