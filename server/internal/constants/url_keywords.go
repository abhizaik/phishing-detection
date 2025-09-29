package constants

// UrlKeywords maps keyword -> risk category
var UrlKeywords = map[string]string{
	//  Login & Authentication
	"login":        "auth",
	"signin":       "auth",
	"signon":       "auth",
	"authenticate": "auth",
	"securelogin":  "auth",
	"sign":         "auth",
	"user":         "auth",
	"session":      "auth",
	"passcode":     "auth",
	"access":       "auth",
	"token":        "auth",

	//  Account Management
	"account":     "account",
	"verify":      "account",
	"update":      "account",
	"reset":       "account",
	"recovery":    "account",
	"unlock":      "account",
	"validate":    "account",
	"approval":    "account",
	"credentials": "account",

	//  Security & Alerts
	"secure":   "security",
	"security": "security",
	"alert":    "security",
	"warning":  "security",
	"notice":   "security",

	//  Financial & Payment
	"paypal":      "finance",
	"banking":     "finance",
	"wallet":      "finance",
	"billing":     "finance",
	"invoice":     "finance",
	"transaction": "finance",
	"payment":     "finance",
	"webscr":      "finance",

	//  Generic Support / Social Engineering
	"support":  "support",
	"helpdesk": "support",
	"contact":  "support",
	"service":  "support",
	"customer": "support",
}
