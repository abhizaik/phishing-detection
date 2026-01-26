package threatfeeds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	PhishTankAPIURL = "http://checkurl.phishtank.com/checkurl/"
)

type PhishTankResponse struct {
	Meta struct {
		Timestamp string `json:"timestamp"`
		ServerID  string `json:"server_id"`
		Status    string `json:"status"`
	} `json:"meta"`
	Results struct {
		URL          string `json:"url"`
		InDatabase   bool   `json:"in_database"`
		PhishID      string `json:"phish_id"`
		PhishDetailURL string `json:"phish_detail_url"`
		Verified     bool   `json:"verified"`
		VerifiedAt   string `json:"verified_at"`
		IsOnline     bool   `json:"online"`
		Target       string `json:"target"`
	} `json:"results"`
}

type PhishTankResult struct {
	InDatabase bool   `json:"in_database"`
	Verified   bool   `json:"verified"`
	IsOnline   bool   `json:"is_online"`
	Target     string `json:"target"`
}

func CheckPhishTank(targetURL string) (*PhishTankResult, error) {
	apiKey := os.Getenv("PHISHTANK_API_KEY")
	userAgent := os.Getenv("PHISHTANK_USER_AGENT")
	if userAgent == "" {
		userAgent = "phishtank/SafeSurf"
	}

	data := url.Values{}
	data.Set("url", targetURL)
	data.Set("format", "json")
	if apiKey != "" {
		data.Set("app_key", apiKey)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("POST", PhishTankAPIURL, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = data.Encode()
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("phishtank api returned status: %d", resp.StatusCode)
	}

	var ptResp PhishTankResponse
	if err := json.NewDecoder(resp.Body).Decode(&ptResp); err != nil {
		return nil, err
	}

	return &PhishTankResult{
		InDatabase: ptResp.Results.InDatabase,
		Verified:   ptResp.Results.Verified,
		IsOnline:   ptResp.Results.IsOnline,
		Target:     ptResp.Results.Target,
	}, nil
}
