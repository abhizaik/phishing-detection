package checks

import (
	"strings"
)

type BrandResult struct {
	BrandFound    string `json:"brand_found"`
	IsMismatch    bool   `json:"is_mismatch"`
	DetectedNames []string `json:"detected_names"`
}

var highValueBrands = map[string][]string{
	"Google":    {"google", "gmail", "youtube"},
	"Microsoft": {"microsoft", "outlook", "office365", "azure", "windows"},
	"Apple":     {"apple", "icloud", "itunes", "iphone"},
	"Amazon":    {"amazon", "aws"},
	"Facebook":  {"facebook", "meta", "instagram", "whatsapp"},
	"PayPal":    {"paypal"},
	"Netflix":   {"netflix"},
	"Adobe":     {"adobe"},
	"Bank":      {"sbi", "hdfc", "icici", "chase", "bank of america", "wells fargo", "hsbc", "citibank"},
	"Binance":   {"binance", "coinbase", "kraken"},
}

// CheckBrandMismatch looks for brand keywords in the page content/title
// and compares them against the domain.
func CheckBrandMismatch(domain string, pageTitle string, bodyText string) BrandResult {
	domain = strings.ToLower(domain)
	pageTitle = strings.ToLower(pageTitle)
	bodyText = strings.ToLower(bodyText)

	res := BrandResult{
		DetectedNames: []string{},
	}

	for brand, keywords := range highValueBrands {
		brandFound := false
		for _, kw := range keywords {
			if strings.Contains(pageTitle, kw) || strings.Contains(bodyText, " "+kw+" ") {
				brandFound = true
				res.DetectedNames = append(res.DetectedNames, brand)
				break
			}
		}

		if brandFound {
			// Check if the domain actually belongs to this brand
			isOfficial := false
			brandLower := strings.ToLower(brand)
			if strings.Contains(domain, brandLower) {
				isOfficial = true
			}
			for _, kw := range keywords {
				if strings.Contains(domain, kw) {
					isOfficial = true
					break
				}
			}

			// If brand found but domain doesn't match known brand domains, it's a mismatch
			if !isOfficial {
				res.BrandFound = brand
				res.IsMismatch = true
			}
		}
	}

	return res
}
