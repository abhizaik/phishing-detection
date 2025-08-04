package checks

import (
	"net/url"
	"strings"
)

func ContainsPunycode(rawURL string) (bool, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false, err
	}

	host := parsedURL.Hostname()
	labels := strings.Split(host, ".")

	for _, label := range labels {
		if strings.HasPrefix(label, "xn--") {
			return true, nil
		}
	}
	return false, nil
}
