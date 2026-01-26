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

// NOT COMPPLETED
// SSLCertResult holds certificate and chain checks
type SSLCertResult struct {
	Domain        string
	HasTLS        bool
	ChainValid    bool
	Issuer        string
	NotBefore     time.Time
	NotAfter      time.Time
	AgeDays       int
	Fingerprint   string
	IsSuspicious  bool
	Reasons       []string
	CTLogged      bool
	KnownBadChain bool
}

// Example blacklist of SHA256 cert fingerprints (hex, uppercase)
// Move to const dir and add more data from https://github.com/BishopFox/badcerts or abuse.ch
var knownBadFingerprints = map[string]struct{}{
	"DEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF": {},
}

// --- Core analyzer ---
func AnalyzeSSLCert(domain string) SSLCertResult {
	res := SSLCertResult{Domain: domain}

	// dial with timeout
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true, // we validate chain manually
	})
	if err != nil {
		res.HasTLS = false
		res.Reasons = append(res.Reasons, fmt.Sprintf("TLS connection failed: %v", err))
		return res
	}
	defer conn.Close()
	res.HasTLS = true

	// grab peer certs
	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		res.Reasons = append(res.Reasons, "no peer certificates")
		return res
	}

	cert := state.PeerCertificates[0]
	res.Issuer = cert.Issuer.CommonName
	res.NotBefore = cert.NotBefore
	res.NotAfter = cert.NotAfter
	res.AgeDays = int(time.Since(cert.NotBefore).Hours() / 24)

	// fingerprint
	fp := sha256.Sum256(cert.Raw)
	res.Fingerprint = strings.ToUpper(hex.EncodeToString(fp[:]))

	// validate chain using system roots
	opts := x509.VerifyOptions{
		Roots:         nil, // use system roots
		Intermediates: x509.NewCertPool(),
		DNSName:       domain,
	}
	for _, ic := range state.PeerCertificates[1:] {
		opts.Intermediates.AddCert(ic)
	}
	if _, err := cert.Verify(opts); err != nil {
		res.ChainValid = false
		res.Reasons = append(res.Reasons, fmt.Sprintf("cert chain invalid: %v", err))
	} else {
		res.ChainValid = true
	}

	// check blacklist
	if _, bad := knownBadFingerprints[res.Fingerprint]; bad {
		res.IsSuspicious = true
		res.KnownBadChain = true
		res.Reasons = append(res.Reasons, "certificate fingerprint is blacklisted")
	}

	// expiry checks
	if time.Now().After(cert.NotAfter) {
		res.IsSuspicious = true
		res.Reasons = append(res.Reasons, "certificate expired")
	}
	if time.Now().Before(cert.NotBefore) {
		res.IsSuspicious = true
		res.Reasons = append(res.Reasons, "certificate not yet valid")
	}

	// weak validity period (e.g., > 398 days is discouraged by CABF)
	if cert.NotAfter.Sub(cert.NotBefore).Hours()/24 > 398 {
		res.Reasons = append(res.Reasons, "unusually long validity period")
	}

	// Check for embedded Certificate Transparency (CT) SCTs
	// OID 1.3.6.1.4.1.11129.2.4.2 is for embedded SCTs
	// This feature is incomplete, needs to be implemented, removing from UI for now.
	res.CTLogged = false
	for _, ext := range cert.Extensions {
		if ext.Id.String() == "1.3.6.1.4.1.11129.2.4.2" {
			res.CTLogged = true
			break
		}
	}

	if !res.CTLogged {
		res.IsSuspicious = true
		res.Reasons = append(res.Reasons, "certificate does not contain embedded CT logs (SCTs)")
	}

	return res
}

