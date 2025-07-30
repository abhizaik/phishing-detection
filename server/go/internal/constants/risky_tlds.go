package constants

// Risky TLDs which are often used by attackers as they are cheap and easy to get
var RiskyTLDs = map[string]struct{}{
	"xyz":   {},
	"top":   {},
	"tk":    {},
	"ml":    {},
	"ga":    {},
	"cf":    {},
	"gq":    {},
	"click": {},
	"zip":   {},
}
