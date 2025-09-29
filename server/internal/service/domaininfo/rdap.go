package domaininfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type rdapResponse struct {
	LDHName     string `json:"ldhName"`
	Nameservers []struct {
		LDHName string `json:"ldhName"`
	} `json:"nameservers"`
	Events []struct {
		Action string    `json:"eventAction"`
		Date   time.Time `json:"eventDate"`
	} `json:"events"`
	Entities []struct {
		Roles      []string `json:"roles"`
		VCardArray []any    `json:"vcardArray"`
	} `json:"entities"`
	Status    []string `json:"status"`
	SecureDNS struct {
		DelegationSigned bool `json:"delegationSigned"`
	} `json:"secureDNS"`
}

// fetchRDAP queries RDAP and returns normalized RegistrationData.
func fetchRDAP(domain string) (*RegistrationData, error) {
	tld := strings.Split(domain, ".")
	if len(tld) < 2 {
		return nil, fmt.Errorf("invalid domain")
	}

	// Get the TLD (last part of domain)
	tldPart := tld[len(tld)-1]

	// Try to find appropriate RDAP server for the TLD
	rdapURL, err := getRDAPServer(tldPart)
	if err != nil {
		return nil, fmt.Errorf("RDAP not supported for TLD .%s: %w", tldPart, err)
	}

	url := fmt.Sprintf("%s/domain/%s", rdapURL, domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("RDAP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RDAP query failed: %s", resp.Status)
	}

	var rd rdapResponse
	if err := json.NewDecoder(resp.Body).Decode(&rd); err != nil {
		return nil, fmt.Errorf("RDAP response parsing failed: %w", err)
	}

	var created, updated, expiry time.Time
	for _, e := range rd.Events {
		switch e.Action {
		case "registration":
			created = e.Date
		case "last changed":
			updated = e.Date
		case "expiration":
			expiry = e.Date
		}
	}

	var ns []string
	for _, n := range rd.Nameservers {
		ns = append(ns, n.LDHName)
	}

	raw, _ := json.Marshal(rd)

	return &RegistrationData{
		Domain:      rd.LDHName,
		CreatedDate: created,
		UpdatedDate: updated,
		ExpiryDate:  expiry,
		Nameservers: ns,
		Status:      rd.Status,
		DNSSEC:      rd.SecureDNS.DelegationSigned,
		Raw:         string(raw),
		Source:      "rdap",
	}, nil
}

// getRDAPServer returns the RDAP server URL for a given TLD
func getRDAPServer(tld string) (string, error) {
	// Common RDAP servers for popular TLDs
	rdapServers := map[string]string{
		"com":  "https://rdap.verisign.com/com/v1",
		"net":  "https://rdap.verisign.com/net/v1",
		"org":  "https://rdap.pir.org/rdap/org/v1",
		"info": "https://rdap.afilias.net/rdap/info/v1",
		"biz":  "https://rdap.afilias.net/rdap/biz/v1",
		"io":   "https://rdap.nic.io/v1",
		"co":   "https://rdap.nic.co/v1",
		"me":   "https://rdap.nic.me/v1",
		"tv":   "https://rdap.nic.tv/v1",
		"cc":   "https://rdap.nic.cc/v1",
	}

	if server, exists := rdapServers[tld]; exists {
		return server, nil
	}

	// For TLDs not in our list, we could implement IANA bootstrap file lookup
	// For now, return an error indicating RDAP is not supported
	return "", fmt.Errorf("no RDAP server configured for .%s", tld)
}
