package screenshot

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func TakeScreenshot(url string) string {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// timeout to avoid hanging
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var buf []byte

	// run task
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second), // let page load
		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		log.Fatal(err)
	}

	timestamp := time.Now().Format("20060102-150405")
	urlStr := sanitizeURL(url)
	filename := "screenshot-" + timestamp + "-" + urlStr + ".png"
	dir := filepath.Join(".", "server", "tmp", "screenshots") // ./server/tmp/screenshots

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	fullPath := filepath.Join(dir, filename)

	if err := os.WriteFile(fullPath, buf, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Screenshot saved")
	return fullPath
}

func sanitizeURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "sanitizeURL : invalid-url"
	}

	combined := parsed.Host + parsed.Path
	if parsed.RawQuery != "" {
		combined += "?" + parsed.RawQuery
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	safe := re.ReplaceAllString(combined, "-")

	safe = strings.Trim(safe, "-")
	if len(safe) > 80 {
		safe = safe[:80]
	}
	return safe
}
