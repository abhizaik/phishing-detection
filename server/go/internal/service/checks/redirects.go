package checks

import (
	"errors"
	"net/http"
)

func CheckRedirects(rawURL string) (isRedirected bool, chain []string, finalURL string, chainLength int, err error) {
	var redirects []string

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirects = append(redirects, req.URL.String())
			// Continue following
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}

	resp, err := client.Get(rawURL)
	if err != nil {
		return false, nil, "", 0, err
	}
	defer resp.Body.Close()

	// Initial request is not counted, so add it if any redirects happened
	if len(redirects) > 0 {
		chain = append([]string{rawURL}, redirects...)
		isRedirected = true
	} else {
		chain = []string{rawURL}
		isRedirected = false
	}

	chainLength = len(chain)
	finalURL = chain[chainLength-1]

	return isRedirected, chain, finalURL, chainLength, nil
}
