package constants

// Not currently in use, can add while adding trust/risk score
// UrlKeywordsRiskWeights maps risk category -> risk level

var UrlKeywordsRiskWeights = map[string]int{
	"auth":     3,
	"account":  2,
	"security": 3,
	"finance":  4,
	"support":  1,
}
