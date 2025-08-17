package checks

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

func GetTld(domain string) (string, bool) {
	tld, icann := publicsuffix.PublicSuffix(domain)
	return tld, icann
}

func GetDomain(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()
	domain, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return "", err
	}
	return domain, nil
}

func GetHost(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()
	return host, nil
}

func IsValidURL(rawURL string) (*url.URL, bool, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL // assume https by default
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, false, err
	}
	// parsed url, is valid url, err
	return parsed, parsed.Scheme != "" && parsed.Host != "", nil
}

func GetDomainAge(created time.Time) (string, int, error) {
	now := time.Now()
	if created.After(now) {
		return "not yet registered", 0, nil
	}

	years := now.Year() - created.Year()
	months := int(now.Month()) - int(created.Month())
	days := int(now.Sub(created).Hours() / 24)

	if months < 0 {
		years--
		months += 12
	}

	if years <= 0 && months <= 0 {
		switch {
		case days == 0:
			return "registered today", days, nil
		case days == 1:
			return "1 day old", days, nil
		case days < 30:
			return fmt.Sprintf("%d days old", days), days, nil
		default:
			return "less than a month old", days, nil
		}
	}

	parts := []string{}
	if years > 0 {
		if years == 1 {
			parts = append(parts, "1 year")
		} else {
			parts = append(parts, fmt.Sprintf("%d years", years))
		}
	}
	if months > 0 {
		if months == 1 {
			parts = append(parts, "1 month")
		} else {
			parts = append(parts, fmt.Sprintf("%d months", months))
		}
	}
	return strings.Join(parts, " "), days, nil
}
