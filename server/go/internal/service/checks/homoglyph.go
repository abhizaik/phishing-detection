package checks

import "unicode"

func HasHomoglyphs(rawURL string) (bool, error) {
	host, err := GetDomain(rawURL)
	if err != nil {
		return false, err
	}

	// If any non-ASCII letter appears, flag it
	for _, r := range host {
		if r > unicode.MaxASCII && unicode.IsLetter(r) {
			return true, nil
		}
	}
	return false, nil
}
