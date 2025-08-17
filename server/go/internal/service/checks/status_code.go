package checks

import "net/http"

func GetStatusCode(url string) (int, string, bool, bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, "", false, false, err
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	statusText := http.StatusText(statusCode)
	isSuccess := statusCode >= 200 && statusCode < 300
	isRedirectStatusCode := statusCode >= 300 && statusCode < 400

	return statusCode, statusText, isSuccess, isRedirectStatusCode, nil
}
