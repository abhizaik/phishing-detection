package checks

import (
	"time"

	"github.com/likexian/whois"
	whois_parser "github.com/likexian/whois-parser"
)

func GetWhoisData(domain string) (whois_parser.WhoisInfo, string, error) {
	raw, err := whois.Whois(domain)
	if err != nil {
		return whois_parser.WhoisInfo{}, "Whois fetch error", err
	}

	whoisData, err := whois_parser.Parse(raw)
	if err != nil {
		return whois_parser.WhoisInfo{}, "Whois parsing error", err
	}

	createdStr := whoisData.Domain.CreatedDate
	if createdStr == "" {
		return whoisData, "Invalid created date", nil
	}

	layouts := []string{
		"2006-01-02T15:04:05Z", // ISO8601
		"2006-01-02 15:04:05",  // Common format
		"2006-01-02",           // Date only
		"02-Jan-2006",          // WHOIS format
		"2006.01.02",           // Alternative
	}

	var created time.Time
	for _, layout := range layouts {
		created, err = time.Parse(layout, createdStr)
		if err == nil {
			break
		}
	}

	if created.IsZero() {
		return whoisData, "Invalid created date", nil
	}

	age, err := GetDomainAge(created)

	if err != nil {
		return whoisData, "Couldn't get age", err
	}

	return whoisData, age, nil
}
