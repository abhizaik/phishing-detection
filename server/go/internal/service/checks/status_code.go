package checks

import "net/http"

func GetStatusCode(url string) (int, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	return resp.StatusCode, nil
}
