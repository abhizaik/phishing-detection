package checks

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// FormInfo describes a discovered form
type FormInfo struct {
	Action           string   // resolved absolute action URL (may be empty -> current URL)
	Method           string   // GET/POST etc.
	Inputs           []string // raw input names/types/placeholder snippets
	ContainsPassword bool
	ContainsUserLike bool // username/email-like inputs
	SubmitTexts      []string
	ExternalAction   bool // action host != page host (true if external)
}

// PageFormResult summarizes page-level findings
type PageFormResult struct {
	URL           string
	HasForms      bool
	HasLoginForm  bool
	FormCount     int
	Forms         []FormInfo
	FetchDuration time.Duration
}

// GetPageFormInfo fetches the page (with timeout), parses HTML and returns form info.
// It does not execute JavaScript. For JS-rendered pages use a headless browser.
func GetPageFormInfo(pageURL string) (*PageFormResult, error) {
	start := time.Now()

	// HTTP client with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", pageURL, nil)
	if err != nil {
		return nil, err
	}
	// polite UA
	req.Header.Set("User-Agent", "PhishScanner/1.0 (+https://example.local)")

	client := &http.Client{
		Timeout: 12 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read body (limited)
	body, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024)) // 5MB cap
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	uParsed, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}
	pageHost := uParsed.Hostname()

	var results []FormInfo

	// traverse nodes to find <form>
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "form" {
			f := extractFormInfo(n, uParsed, pageHost)
			results = append(results, f)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	res := &PageFormResult{
		URL:           pageURL,
		HasForms:      len(results) > 0,
		HasLoginForm:  false,
		FormCount:     len(results),
		Forms:         results,
		FetchDuration: time.Since(start),
	}

	// mark if any form looks like a login form
	for _, f := range results {
		if f.ContainsPassword || f.ContainsUserLike {
			res.HasLoginForm = true
			break
		}
	}
	return res, nil
}

// extractFormInfo inspects a <form> node and returns a FormInfo
func extractFormInfo(form *html.Node, base *url.URL, pageHost string) FormInfo {
	info := FormInfo{
		Method: strings.ToUpper(getAttr(form, "method")),
	}
	if info.Method == "" {
		info.Method = "GET"
	}
	rawAction := getAttr(form, "action")
	actionResolved := resolveAction(base, rawAction)
	info.Action = actionResolved

	// determine if action host differs
	if actionResolved != "" {
		if parsed, err := url.Parse(actionResolved); err == nil {
			info.ExternalAction = !sameHost(parsed.Hostname(), pageHost)
		}
	}

	// collect inputs and buttons inside the form
	var ftr func(*html.Node)
	ftr = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "input":
				it := strings.ToLower(getAttr(n, "type"))
				name := getAttr(n, "name")
				placeholder := getAttr(n, "placeholder")
				aria := getAttr(n, "aria-label")
				id := getAttr(n, "id")
				info.Inputs = append(info.Inputs, fmtInputSummary(it, name, placeholder, aria, id))

				// password detection
				if it == "password" || strings.Contains(strings.ToLower(name), "pass") || strings.Contains(strings.ToLower(id), "pass") || strings.Contains(strings.ToLower(placeholder), "pass") {
					info.ContainsPassword = true
				}
				// username/email-ish detection
				if it == "email" ||
					strings.Contains(strings.ToLower(name), "user") ||
					strings.Contains(strings.ToLower(name), "login") ||
					strings.Contains(strings.ToLower(name), "email") ||
					strings.Contains(strings.ToLower(id), "user") ||
					strings.Contains(strings.ToLower(placeholder), "user") ||
					strings.Contains(strings.ToLower(aria), "user") ||
					strings.Contains(strings.ToLower(placeholder), "email") {
					info.ContainsUserLike = true
				}
			case "button":
				// capture text content of button
				txt := strings.TrimSpace(nodeText(n))
				if txt != "" {
					info.SubmitTexts = append(info.SubmitTexts, txt)
					l := strings.ToLower(txt)
					if looksLikeLoginText(l) {
						info.ContainsUserLike = info.ContainsUserLike || false // don't override password flag
					}
				}
			case "a":
				// sometimes login is an <a> styled as button
				txt := strings.TrimSpace(nodeText(n))
				if txt != "" && looksLikeLoginText(strings.ToLower(txt)) {
					info.SubmitTexts = append(info.SubmitTexts, txt)
				}
			case "label":
				// label text may indicate username/password fields
				txt := strings.ToLower(strings.TrimSpace(nodeText(n)))
				if strings.Contains(txt, "password") {
					info.ContainsPassword = true
				}
				if strings.Contains(txt, "username") || strings.Contains(txt, "email") || strings.Contains(txt, "sign in") {
					info.ContainsUserLike = true
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			ftr(c)
		}
	}
	ftr(form)

	// additional heuristic: if form has a password OR (text inputs + "sign in" in action or form text) -> login-like
	if !info.ContainsUserLike {
		// check action path / query for login keywords
		if strings.Contains(strings.ToLower(info.Action), "login") || strings.Contains(strings.ToLower(info.Action), "signin") || strings.Contains(strings.ToLower(info.Action), "auth") {
			info.ContainsUserLike = true
		}
	}
	return info
}

// helpers

func sameHost(a, b string) bool {
	return strings.EqualFold(a, b)
}

func resolveAction(base *url.URL, raw string) string {
	if raw == "" || raw == "#" {
		// empty action means submit to same URL
		return base.String()
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	return base.ResolveReference(parsed).String()
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val
		}
	}
	return ""
}

// nodeText returns concatenated text of a node subtree
func nodeText(n *html.Node) string {
	var b strings.Builder
	var walker func(*html.Node)
	walker = func(nn *html.Node) {
		if nn.Type == html.TextNode {
			b.WriteString(nn.Data)
		}
		for c := nn.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(n)
	return strings.TrimSpace(b.String())
}

func fmtInputSummary(typ, name, placeholder, aria, id string) string {
	parts := []string{}
	if typ != "" {
		parts = append(parts, "type="+typ)
	}
	if name != "" {
		parts = append(parts, "name="+name)
	}
	if id != "" {
		parts = append(parts, "id="+id)
	}
	if placeholder != "" {
		parts = append(parts, "ph="+placeholder)
	}
	if aria != "" {
		parts = append(parts, "aria="+aria)
	}
	return strings.Join(parts, "|")
}

func looksLikeLoginText(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	keywords := []string{"login", "log in", "sign in", "signin", "submit", "sign-on", "signon"}
	for _, k := range keywords {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}
