package checks

import (
	"errors"
	"net/http"
	"net/url"
)

type RedirectionResult struct {
	IsRedirected  bool     `json:"is_redirected"`
	ChainLength   int      `json:"chain_length"`
	Chain         []string `json:"chain"`
	FinalURL      string   `json:"final_url"`
	HasDomainJump bool     `json:"has_domain_jump"`
}

func CheckRedirects(rawURL string) (RedirectionResult, error) {
	var redirects []string

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirects = append(redirects, req.URL.String())
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}

	resp, err := client.Get(rawURL)
	if err != nil {
		return RedirectionResult{}, err
	}
	defer resp.Body.Close()

	chain := append([]string{rawURL}, redirects...)
	finalURL := chain[len(chain)-1]

	// Detect domain jumps
	parsedStart, err := url.Parse(rawURL)
	if err != nil {
		return RedirectionResult{}, err
	}
	origDomain := parsedStart.Host
	hasJump := false

	for _, u := range chain[1:] {
		parsed, err := url.Parse(u)
		if err != nil {
			continue
		}
		if parsed.Host != origDomain {
			hasJump = true
			break
		}
	}

	return RedirectionResult{
		IsRedirected:  len(redirects) > 0,
		Chain:         chain,
		FinalURL:      finalURL,
		ChainLength:   len(chain),
		HasDomainJump: hasJump,
	}, nil
}
