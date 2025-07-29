package checks

import "testing"

func TestTooLongUrl(t *testing.T) {
	url := "https://google.com"
	result := TooLongUrl(url)
	if result == true {
		t.Errorf("Error while testing TooLongUrl(%v)", url)
	}
}
