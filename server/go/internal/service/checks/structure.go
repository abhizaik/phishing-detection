package checks

import "strings"

func TooLongUrl(url string) bool {
	return len(url) > 75
}

func TooDeepUrl(url string) bool {
	return strings.Count(url, "/") > 5
}
