package domaininfo

import "time"

// RegistrationData holds normalized domain registry info
// regardless of whether it came from WHOIS or RDAP.
type RegistrationData struct {
	Domain      string    `json:"domain"`
	Registrar   string    `json:"registrar"`
	CreatedDate time.Time `json:"created"`
	UpdatedDate time.Time `json:"updated"`
	ExpiryDate  time.Time `json:"expiry"`
	Nameservers []string  `json:"nameservers"`
	Status      []string  `json:"status"`
	DNSSEC      bool      `json:"dnssec"`
	AgeHuman    string    `json:"age_human"` // e.g. "2 years 3 months"
	AgeDays     int       `json:"age_days"`  // total days since registration
	Raw         string    `json:"raw"`       // Raw WHOIS or RDAP JSON for debugging
	Source      string    `json:"source"`    // "whois" or "rdap"
}
