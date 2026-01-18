<div align="center">

# Phishing Detection Engine (Open-Source, Real-Time URL Scanner)

An **open-source, production-grade phishing detection engine** for fast, explainable, and auditable URL analysis.
Detect phishing and malicious URLs with **transparent scoring, clear verdicts, and detailed reports** — all in under a second.
Designed for **self-hosting, auditing and easy integration** into web apps, APIs, and browser extensions.

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go\&logoColor=white)](https://go.dev)
[![Svelte](https://img.shields.io/badge/Svelte-5-orange?logo=svelte)](https://svelte.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/abhizaik/phishing-detection?style=social)](https://github.com/abhizaik/phishing-detection)

[⚡ Quick Start](#quick-start) · [🏛 Architecture](#architecture) · [📚 Docs](#documentation) · [🤝 Contributing](#contributing) · [🌍 Community](#community)

</div>

---

## Phishing Detection Demo

> Paste a URL → get a **trust score, verdict, and detailed report** in under a second.

![Phishing Detection Demo](assets/demo.gif)

---

## Why Use This Tool?

Most phishing detection solutions are either **closed commercial APIs** or **academic ML demos**:

* **Commercial tools**: expensive, opaque, and impossible to audit
* **ML demos**: slow, fragile, and not built for real-world deployment

**Phishing remains a top cyber threat** because defenders lack **fast, explainable, and controllable detection systems**.
This engine fills that gap by providing:

* **Transparent, explainable analysis** — every verdict is backed by concrete signals
* **Fast, real-time scanning** — multiple analyzers run in parallel
* **Flexible integration** — web UI, HTTP API, browser extension
* **Full open-source control** — audit, modify, self-host, and scale

---

## What This Tool Does

* Scans URLs for **phishing, malicious behavior, and unsafe redirects**
* Produces a **trust score, clear verdict, and detailed report**
* Supports **developers and non-technical users** via UI, API, and extension
* Uses **multiple independent heuristic analyzers** for accurate detection
* Built with **Go (backend)** and **Svelte (frontend)** for production use

---

## Who This Is For

**Everyday users**

* Quickly check suspicious URLs in a website or browser extension

**Developers**

* Integrate phishing detection into applications or backend services
* Replace or supplement commercial phishing APIs

**Security engineers & SOC teams**

* Build explainable phishing detection pipelines
* Audit URLs with transparent, actionable signals

**Students & researchers**

* Use a real-world, production-grade reference for **academic or security projects**




## Use Cases

- Detect phishing links before users click them
- Scan URLs for malicious behavior
- Build anti-phishing browser extensions
- Integrate phishing detection into backend services
- Replace or supplement commercial phishing APIs


## Detection Techniques Used

The engine evaluates URLs using multiple independent analyzers, including:

- Domain reputation & age checks
- Suspicious URL patterns and homoglyphs
- Redirect chain analysis
- HTTPS / certificate anomalies
- Known phishing indicators and heuristics
- Content-based signals (HTML, scripts, forms)

Each analyzer contributes to a final trust score and verdict.


## Architecture
High level repository layout:

```text
server/               Go backend 
  cmd/safesurf        Backend entry point
  internal/           Analyzers, domaininfo, screenshot
web/website           SvelteKit UI
web/chrome-extension  Chrome extension
docker/               Dev & prod
docs/                 Setup, architecture, API, security, testing etc.
Makefile
```



## Quick Start

Full setup: [docs/setup.md](docs/setup.md) 

1. Clone the repo

```bash
git clone https://github.com/abhizaik/phishing-detection.git
cd phishing-detection
```
2. Start the application

Prerequisite: Docker must be installed and running. <br>
Windows: Use WSL or install make.
```bash
make build
make up
```
Web UI: **[localhost:3000](http://localhost:3000)** 

## API Example

The phishing detection engine exposes a simple HTTP API for real-time URL analysis.
Scan a URL using the API:

```bash
curl -X GET http://localhost:8080/api/v1/analyze?url=https://example.com
```
<details>
<summary>Example API response</summary>
```json
{
    "url": "http://google.com/abhi",
    "domain": "google.com",
    "features": {
        "rank": 1,
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
                "categories": {}
            }
        }
    },
    "infrastructure": {
        "ip_addresses": [
            "142.250.182.78",
            "2404:6800:4007:810::200e"
        ],
        "nameservers_valid": true,
        "ns_hosts": [
            "ns2.google.com."
        ],
        "mx_records_valid": true,
        "mx_hosts": [
            "smtp.google.com."
        ]
    },
    "domain_info": {
        "domain": "GOOGLE.COM",
        "registrar": "MarkMonitor Inc.",
        "created": "1997-09-15T04:00:00Z",
        "updated": "2019-09-09T15:39:04Z",
        "expiry": "2028-09-14T04:00:00Z",
        "nameservers": [
            "NS1.GOOGLE.COM"
        ],
        "status": [
            "client delete prohibited"
        ],
        "dnssec": false,
        "age_human": "28 years 4 months",
        "age_days": 10350,
        "raw": "{\"ldhName\":\"GOOGLE.COM\ etc."}",
        "source": "RDAP"
    },
    "analysis": {
        "redirection_result": {
            "is_redirected": false,
            "chain_length": 1,
            "chain": [
                "http://google.com/abhi"
            ],
            "final_url": "http://google.com/abhi",
            "final_url_domain": "google.com",
            "has_domain_jump": false
        },
        "http_status": {
            "code": 404,
            "text": "Not Found",
            "success": false,
            "is_redirect": false
        },
        "is_hsts_supported": false
    },
    "result": {
        "risk_score": 0,
        "trust_score": 100,
        "final_score": 100,
        "verdict": "Safe",
        "reasons": {
            "neutral_reasons": [
                "Standard, officially recognized domain extension.",
                "Standard DNS security (DNSSEC not enabled)."
            ],
            "good_reasons": [
                "Global Giant: Ranked #1 worldwide.",
                "Valid DNS configuration detected.",
                "Valid email server configuration (MX Records).",
                "Long-standing domain history (28 years 4 months).",
                "Registered with MarkMonitor Inc."
            ],
            "bad_reasons": null
        }
    },
    "incomplete": false,
    "errors": null
}
```
</details>

## Real-Time Performance

- Typical scan time: **~300–700 ms per URL**
- Designed to handle multiple concurrent scans efficiently
- Optimized for **real-time phishing detection at scale**

Exact performance depends on enabled analyzers and network conditions.


## Security & Privacy

- No URL data is sent to third-party services by default
- All analysis runs locally or in your own infrastructure
- Designed for auditability, privacy, and controlled environments

## Limitations
- Heuristic false positives

##  Documentation

All documentation is under `docs/`. Start here [docs/README.md](docs/README.md) 



## Contributing

Bug reports, feature requests, and pull requests are welcome.

Use [GitHub Issues](https://github.com/abhizaik/phishing-detection/issues) to report bugs or suggest features. For code contributions, see [CONTRIBUTING.md](.github/CONTRIBUTING.md).




## Community

**If you found this project helpful, consider giving it a star.** It directly helps visibility and continued development.

Have bugs, ideas, or feature requests?
Open an [issue](https://github.com/abhizaik/phishing-detection/issues) or start a [discussion](https://github.com/abhizaik/phishing-detection/discussions). Contributions and feedback are welcome.

Thanks for helping make the web safer.




<div align="center">
  <a href="https://star-history.com/#abhizaik/phishing-detection&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=Date&theme=dark" />
      <source media="(prefers-color-scheme: light)" srcset="https://api/star-history.com/svg?repos=abhizaik/phishing-detection&type=Date" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=Date" />
    </picture>
  </a>
</div>

<!-- [![Community Growth Trajectory](https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=date&legend=top-left)](https://www.star-history.com/#abhizaik/phishing-detection&type=date&legend=top-left) -->
