package checks

import (
	"regexp"
	"strings"

	"github.com/abhizaik/SafeSurf/internal/constants"
)

func CheckURLKeywords(url string) (bool, []string, []string) {
	lowerURL := strings.ToLower(url)

	// Split by non-alphanumeric characters
	re := regexp.MustCompile(`[^a-z0-9]+`)
	words := re.Split(lowerURL, -1)

	matches := []string{}
	categories := []string{}
	keywordPresent := false

	for _, word := range words {
		if word == "" {
			continue
		}
		if category, exists := constants.UrlKeywords[word]; exists {
			keywordPresent = true
			matches = append(matches, word)
			categories = append(categories, category)
		}
	}

	return keywordPresent, matches, categories
}
