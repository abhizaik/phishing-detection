package domaininfo

import (
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func GetWhoisData(domain string) (*RegistrationData, error) {
	raw, err := whois.Whois(domain)
	if err != nil {
		return nil, err
	}

	whoisData, err := whoisparser.Parse(raw)
	if err != nil {
		return nil, err
	}

	// Parse WHOIS date strings to time.Time
	createdDate := parseWhoisDate(whoisData.Domain.CreatedDate)
	updatedDate := parseWhoisDate(whoisData.Domain.UpdatedDate)
	expiryDate := parseWhoisDate(whoisData.Domain.ExpirationDate)

	// Convert WHOIS data to RegistrationData format
	registrationData := &RegistrationData{
		Domain:      whoisData.Domain.Domain,
		Registrar:   whoisData.Registrar.Name,
		CreatedDate: createdDate,
		UpdatedDate: updatedDate,
		ExpiryDate:  expiryDate,
		Nameservers: whoisData.Domain.NameServers,
		Status:      whoisData.Domain.Status,
		DNSSEC:      false, // WHOIS doesn't always provide DNSSEC info
		Raw:         raw,
		Source:      "whois",
	}

	return registrationData, nil
}

// parseWhoisDate attempts to parse WHOIS date strings into time.Time
func parseWhoisDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}

	// Common WHOIS date formats
	layouts := []string{
		"2006-01-02T15:04:05Z",      // ISO8601
		"2006-01-02 15:04:05",       // Common format
		"2006-01-02",                // Date only
		"02-Jan-2006",               // WHOIS format
		"2006.01.02",                // Alternative
		"2006-01-02T15:04:05-07:00", // ISO8601 with timezone
		"2006-01-02T15:04:05",       // ISO8601 without timezone
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, dateStr); err == nil {
			return parsed
		}
	}

	// Return zero time if parsing fails
	return time.Time{}
}
