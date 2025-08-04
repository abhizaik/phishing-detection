package checks

import (
	"errors"
	"net/http"
)

func CheckRedirects(rawURL string) (bool, []string, error) {
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
		return false, nil, err
	}
	defer resp.Body.Close()

	// Initial request is not counted, so add it if any redirects happened
	if len(redirects) > 0 {
		redirects = append([]string{rawURL}, redirects...)
		return true, redirects, nil
	}

	return false, []string{rawURL}, nil
}
