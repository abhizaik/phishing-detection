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

	// stub CT check (replace with real CT API lookup)
	if fakeCheckCTLogs(cert) {
		res.CTLogged = true
	} else {
		res.CTLogged = false
		res.IsSuspicious = true
		res.Reasons = append(res.Reasons, "certificate not found in CT logs (stubbed)")
	}

	return res
}

// fakeCheckCTLogs simulates a CT log lookup
func fakeCheckCTLogs(cert *x509.Certificate) bool {
	// Replace with real CT query against Google/Cloudflare logs
	// e.g., using crt.sh API or CertSpotter API
	return false // default: not found
}

// Example runner
func SSLMain() {
	tests := []string{
		"google.com",
		"expired.badssl.com",
	}

	for _, d := range tests {
		r := AnalyzeSSLCert(d)
		fmt.Printf("Domain: %s\n", r.Domain)
		fmt.Printf(" HasTLS: %v ChainValid: %v Issuer: %s\n", r.HasTLS, r.ChainValid, r.Issuer)
		fmt.Printf(" NotBefore: %v NotAfter: %v AgeDays: %d\n", r.NotBefore, r.NotAfter, r.AgeDays)
		fmt.Printf(" Fingerprint: %s\n", r.Fingerprint)
		fmt.Printf(" CTLogged: %v KnownBadChain: %v\n", r.CTLogged, r.KnownBadChain)
		fmt.Printf(" Suspicious: %v\n", r.IsSuspicious)
		if len(r.Reasons) > 0 {
			fmt.Printf(" Reasons: %v\n", r.Reasons)
		}
		fmt.Println(strings.Repeat("-", 60))
	}
}
