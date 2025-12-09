<!-- SEO & Social Metadata -->
<meta name="title" content="SafeSurf — Open-Source Phishing Detection & URL Analysis">
<meta name="description" content="Real-time phishing detection, URL threat analysis, domain intelligence, and risk scoring built with Go and Svelte.">
<meta property="og:title" content="SafeSurf — Phishing Detection & URL Analysis">
<meta property="og:description" content="Analyze and detect phishing URLs using DNS, TLS, redirects, entropy, WHOIS/RDAP, and page-content intelligence.">
<meta property="og:image" content="https://raw.githubusercontent.com/abhizaik/SafeSurf/main/web/static/images/safesurf-normal.png">
<meta property="og:type" content="website">
<meta name="keywords" content="phishing detection, url scanner, malicious url detection, domain risk analysis, security tools, cyber security, Go backend, Svelte frontend, SOC tools, threat intelligence, abhizaik">



<div align="center">

  <picture>
    <img src="./web/static/images/safesurf-normal.png" width="30%" style="border: none; box-shadow: none;" alt="SafeSurf Logo">
  </picture>

</div>

<div align="center">

# Open Source Phishing Detection & URL Analysis

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Svelte](https://img.shields.io/badge/Svelte-5-orange?logo=svelte)](https://svelte.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/abhizaik/SafeSurf?style=social)](https://github.com/abhizaik/SafeSurf)
![Contributors](https://img.shields.io/github/contributors/abhizaik/SafeSurf)

**SafeSurf is a fast, open-source phishing detection engine that analyzes URLs using DNS, TLS, redirects, entropy, ranking, WHOIS/RDAP, page content, and more.**



</div>

---

<div align="center">

[⚡ Quick Start](#-quick-start) ·
[📚 Docs](#-documentation) ·
[🧠 Features](#-features) ·
[🏛 Architecture](#-architecture) ·
[🧪 Testing](#-testing) ·
[🤝 Contributing](#-contributing) ·
[🌍 Community](#-community)

</div>

---

## Why SafeSurf?

- Detects phishing, malicious redirects and homoglyph attacks

- Easy to use for common people (result has a verdict, trust score and detailed report)

- Has web UI, open REST API, and browser extension

- Runs parallel analyzers (fast, real time results)

- Built with Go and Svelte (simple and fast)

- Completely open source


## ⚡ Quick Start

The full setup guide is given in `docs/setup.md`. A short version is given below:

1. Clone the repo

```bash
git clone https://github.com/abhizaik/phishing-detection.git
cd phishing-detection
```
2. Start the application

**Option 1: Docker (recommended)**
```bash
make dev-up
```

**Option 2: Local Go + Svelte**
```bash
cd server && go run ./cmd/safesurf      # backend on :8080
cd ../web/website && npm install && npm run dev   # UI on :5173
```

3. Use the app
- Open the web UI on browser: **[localhost:5173](http://localhost:5173)** 
- Or call the API:
```bash
curl "http://localhost:8080/api/v1/analyze?url=https://example.com"
```

<details>
  <summary>Example Output</summary>

```json
{
  "url": "http://example.com",
  "domain": "example.com",
  "features": {
    "rank": 164,
    "tld": {
      "tld": "com",
      "is_trusted_tld": false,
      "is_risky_tld": false,
      "is_icann": true
    },
    "url": {
      "url_shortener": false,
      "uses_ip": false,
      "contains_punycode": false,
      "too_long": false,
      "too_deep": false,
      "has_homoglyph": false,
      "subdomain_count": 0,
      "keywords": {
        "has_keywords": false,
        "found": [],
        "categories": {

        }
      }
    }
  },
  "infrastructure": {
    "ip_addresses": [
      "23.220.75.245"
    ],
    "nameservers_valid": true,
    "ns_hosts": [
      "b.iana-servers.net."
    ],
    "mx_records_valid": false,
    "mx_hosts": [
      "."
    ]
  },
  "domain_info": {
    "domain": "EXAMPLE.COM",
    "registrar": "RESERVED-Internet Assigned Numbers Authority",
    "created": "1995-08-14T04:00:00Z",
    "updated": "2025-11-25T18:49:24Z",
    "expiry": "2026-08-13T04:00:00Z",
    "nameservers": [
      "A.IANA-SERVERS.NET"
    ],
    "status": [
      "client delete prohibited"
    ],
    "dnssec": true,
    "age_human": "30 years 4 months",
    "age_days": 11075,
    "raw": "{\"ldhName\":\"EXAMPLE.COM\",\"nameservers\":[{\"ldhName\":\"A.IANA-SERVERS.NET\"},{\"ldhName\":\"B.IANA-SERVERS.NET\"}],\"events\":[{\"eventAction\":\"registration\",\"eventDate\":\"1995-08-14T04:00:00Z\"},{\"eventAction\":\"expiration\",\"eventDate\":\"2026-08-13T04:00:00Z\"},{\"eventAction\":\"last changed\",\"eventDate\":\"2025-11-25T18:49:24Z\"},{\"eventAction\":\"last update of RDAP database\",\"eventDate\":\"2025-12-09T16:08:08Z\"}],\"entities\":[{\"roles\":[\"registrar\"],\"vcardArray\":[\"vcard\",[[\"version\",{},\"text\",\"4.0\"],[\"fn\",{},\"text\",\"RESERVED-Internet Assigned Numbers Authority\"]]]}],\"status\":[\"client delete prohibited\",\"client transfer prohibited\",\"client update prohibited\"],\"secureDNS\":{\"delegationSigned\":true}}",
    "source": "RDAP"
  },
  "analysis": {
    "redirection_result": {
      "is_redirected": false,
      "chain_length": 1,
      "chain": [
        "http://example.com"
      ],
      "final_url": "http://example.com",
      "final_url_domain": "example.com",
      "has_domain_jump": false
    },
    "http_status": {
      "code": 200,
      "text": "OK",
      "success": true,
      "is_redirect": false
    },
    "is_hsts_supported": false
  },
  "result": {
    "risk_score": 5,
    "trust_score": 100,
    "final_score": 99,
    "verdict": "Safe",
    "reasons": {
      "neutral_reasons": [
        "Standard, officially recognized domain extension.",
        "No email server configured for this domain."
      ],
      "good_reasons": [
        "Global Giant: Ranked #164 worldwide.",
        "Valid DNS configuration detected.",
        "Long-standing domain history (30 years 4 months).",
        "Registered with RESERVED-Internet Assigned Numbers Authority",
        "Advanced DNS security enabled (DNSSEC)."
      ],
      "bad_reasons": null
    }
  },
  "incomplete": false,
  "errors": null
}
```
</details>


## 📚 Documentation

All detailed docs are under `docs/`:
- **Setup** — local & Docker workflows: `docs/setup.md`
- **Architecture** — components, data flow, diagrams: `docs/architecture.md`, `docs/data-flow.md`
- **Configuration** — env vars, paths, knobs: `docs/configuration.md`
- **CLI & Makefile tooling** — `docs/cli.md`
- **API reference** — endpoints, samples, OpenAPI, Postman: `docs/api.md`
- **Deployment** — Docker, K8s, CI/CD: `docs/deployment.md`
- **Testing & performance** — `docs/testing.md`, `docs/performance.md`
- **Security & operations** — `docs/security.md`, `docs/maintenance.md`
- **Design decisions, FAQ, glossary** — `docs/design-decisions.md`, `docs/faq.md`, `docs/glossary.md`

Start from the [docs index](docs/README.md) for a curated overview.



## 🧠 Features

- **Parallel signal fan‑out**  
  Ranking, DNS, TLS, redirects, entropy, homoglyph tricks, URL patterns, and keyword signals all run in parallel.

- **Deep infrastructure context**  
  Normalized WHOIS/RDAP data, MX/NS checks, IP lookups, and domain-age analysis.

- **Smart URL & content analysis**  
  Looks at URL depth, subdomain patterns, risky TLDs, shorteners, and can pull page content when needed.

- **Screenshot support**  
  Can capture full-page screenshots using headless Chrome.

- **Clear, explainable output**  
  Returns features, infra data, analysis details, timing info, and a final risk score with a simple verdict.

For a detailed look at how the analyzers work and how data moves through the system, check out `docs/architecture.md` and `docs/data-flow.md`.

## 🏛 Architecture

High‑level:
- **Go backend** (`server/`) — Gin HTTP API, analyzer orchestration, rank and domain info services, screenshot worker.
- **Web UI** (`web/website`) — Svelte + Vite frontend for manual analysis.
- **Browser extension** (`web/chrome-extension`) — Chrome MV3 helper calling the same REST API.
- **Operations & docs** (`docker/`, `docs/`, `Makefile`) — Compose stacks, deployment and security guides, testing docs.

Project structure:
```text
server/               Go backend (handlers, analyzers, services)
  cmd/safesurf        Main entry point
  internal/           Analyzer tasks, rank cache, domaininfo, screenshot
web/website           SvelteKit UI
web/chrome-extension  Chrome MV3 extension
docker/               Dev & prod Compose stacks
docs/                 Architecture, setup, API, security, testing, etc.
```

See `docs/architecture.md` for diagrams and more detail.




## 🧪 Testing

- **Backend**
  ```bash
  cd server
  go test ./...
  ```


- **Frontend**
  ```bash
  cd web/website
  npm run check       # type checks (if configured)
  npm test            # when tests are added
  ```

See `docs/testing.md` for integration tests, load testing, and coverage tips.



## 🤝 Contributing

Bug reports and feature ideas are welcome in GitHub Issues or Discussions.

1. Fork the repo and create a feature branch (for example, `feat/feature-name`).
2. Implement your change with tests (`make test-backend`) and keep Go/Svelte code formatted.
3. Update or add docs in `docs/` if behavior changes.
4. Open a PR with description 




## 🌍 Community

If SafeSurf helps you, **leave a star** so others can find it!

Feel free to start issues or discussions if you want fixes or new features added.

Thanks for helping make the web safer.



<!-- ## Community Growth Trajectory -->
<!-- [![Community Growth Trajectory](https://api.star-history.com/svg?repos=abhizaik/safesurf&type=date&legend=top-left)](https://www.star-history.com/#abhizaik/safesurf&type=date&legend=top-left) -->

<div align="center">
  <a href="https://star-history.com/#abhizaik/SafeSurf&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=abhizaik/SafeSurf&type=Date&theme=dark" />
      <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=abhizaik/SafeSurf&type=Date" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=abhizaik/SafeSurf&type=Date" style="border-radius: 15px; box-shadow: 0 0 30px rgba(0, 217, 255, 0.3);" />
    </picture>
  </a>
</div>
