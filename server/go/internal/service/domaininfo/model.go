package domaininfo

import "time"

// RegistrationData holds normalized domain registry info
// regardless of whether it came from WHOIS or RDAP.
type RegistrationData struct {
	Domain      string
	Registrar   string
	CreatedDate time.Time
	UpdatedDate time.Time
	ExpiryDate  time.Time
	Nameservers []string
	Status      []string
	DNSSEC      bool
	Raw         string // Raw WHOIS or RDAP JSON for debugging
	Source      string // "whois" or "rdap"
}
