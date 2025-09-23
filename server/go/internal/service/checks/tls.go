package checks

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

// PageFormResult summarizes page-level findings
type TLSResult struct {
	Present          bool
	Issuer           string
	AgeDays          int
	HostnameMismatch bool
}

func GetTLSInfo(domain string) (TLSResult, error) {
	res := TLSResult{}
	// Add port if not included
	hostPort := domain
	if !strings.Contains(domain, ":") {
		hostPort = net.JoinHostPort(domain, "443")
	}

	conn, err := tls.Dial("tcp", hostPort, &tls.Config{
		InsecureSkipVerify: true, // we check hostname manually
	})
	if err != nil {
		return res, fmt.Errorf("TLS connection failed: %v", err)
	}
	defer conn.Close()

	// Grab peer certificates
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return res, nil
	}
	cert := certs[0]

	// Presence
	present := true
	var issuer string
	var hostnameMismatch bool

	// Issuer
	if len(cert.Issuer.Organization) > 0 {
		issuer = cert.Issuer.Organization[0]
	} else {
		issuer = cert.Issuer.CommonName
	}

	// Age (in days since NotBefore)
	ageDays := int(time.Since(cert.NotBefore).Hours() / 24)

	// Hostname mismatch
	if err := cert.VerifyHostname(domain); err != nil {
		hostnameMismatch = true
	} else {
		hostnameMismatch = false
	}

	res = TLSResult{
		Present:          present,
		Issuer:           issuer,
		AgeDays:          ageDays,
		HostnameMismatch: hostnameMismatch,
	}

	return res, nil
}
