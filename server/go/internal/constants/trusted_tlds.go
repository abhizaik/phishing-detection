package constants

// Highly trustable TLDs which are only given to trusted entities after proper verification
var TrustedTLDs = map[string]struct{}{
	"gov":       {},
	"mil":       {},
	"edu":       {},
	"int":       {},
	"bank":      {},
	"insurance": {},
	"pharmacy":  {},
	"post":      {},
	"museum":    {},
	"aero":      {},
	"ac.in":     {},
	"edu.in":    {},
	"gov.in":    {},
	"ac.uk":     {},
	"gov.uk":    {},
	"edu.au":    {},
	"gov.au":    {},
	"ac.jp":     {},
	"go.jp":     {},
	"edu.cn":    {},
	"gov.cn":    {},
	"ac.za":     {},
}
