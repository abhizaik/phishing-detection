package checks

import (
	"strings"
)

func ContainsPunycode(rawURL string) (bool, error) {
	host, err := GetDomain(rawURL)
	if err != nil {
		return false, err
	}

	labels := strings.SplitSeq(host, ".")

	for label := range labels {
		if strings.HasPrefix(label, "xn--") || strings.HasPrefix(label, ".xn--") {
			return true, nil
		}
	}
	return false, nil
}
