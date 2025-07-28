package checks

import (
	"log"
	"net/http"
	"strings"
)

func IncludeProtocol(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "https://" + url
}

func TooLongUrl(url string) bool {
	return len(url) > 75
}

func TooDeepUrl(url string) bool {
	return strings.Count(url, "/") > 5
}

func GetStatusCode(url string) int {
	resp, err := http.Head(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	return resp.StatusCode
}
