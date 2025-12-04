<div align="center">

  <picture>
    <img src="./web/static/images/safesurf-normal.png" width="22%" style="border: none; box-shadow: none;" alt="SafeSurf Logo">
  </picture>

</div>

<div align="center">

# SafeSurf · Phishing Intelligence

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Svelte](https://img.shields.io/badge/Svelte-5-orange?logo=svelte)](https://svelte.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/abhizaik/SafeSurf?style=social)](https://github.com/abhizaik/SafeSurf)
![Contributors](https://img.shields.io/github/contributors/abhizaik/SafeSurf)

**Real‑time phishing and domain risk intelligence for security teams, SOCs, and browser clients.**

</div>

---

<div align="center">

[🚀 Quick Start](#-quick-start) ·
[📚 Docs](#-documentation) ·
[🧠 Features](#-core-capabilities) ·
[🏗️ Architecture](#-architecture--project-layout) ·
[🧪 Testing](#-testing--quality) ·
[🤝 Contributing](#-contributing)  

</div>

---

## What is SafeSurf?

SafeSurf is a phishing detection and URL intelligence engine. It fans out multiple analyzers in parallel (DNS, TLS, redirects, entropy, homoglyphs, keywords, rank, content, domain info) and returns both:
- **Machine‑readable signals** that are easy to integrate into SIEM, SOAR, and browser extensions.
- **Human‑friendly context** for common people, security analysts and incident responders.

SafeSurf powers:
- A **REST API** you can drop behind an API gateway.
- A **web tool** for manual investigations.
- A **chrome extension** that flags suspicious URLs directly in the browser.

## 🧠 Core Capabilities

- **Parallel signal fan‑out**  
  Rank, DNS, TLS, redirects, entropy, homoglyphs, URL structure, and keyword-based heuristics run concurrently with per‑task timing.

- **Deep infrastructure context**  
  WHOIS/RDAP normalization, MX/NS health checks, IP resolution, and domain age analysis via `domaininfo` services.

- **Lexical & content analysis**  
  URL length/depth, subdomain patterns, risky/trusted TLD sets, URL shortener detection, and page content extraction hooks.

- **Evidence generation**  
  Full‑page screenshots via headless Chrome (`chromedp`), stored under `server/tmp/screenshots` for later review.

- **Explainable results**  
  Responses include features, infrastructure, analysis, performance timings, and a synthesized `result` section with a risk score.

For a deeper walkthrough of analyzers and data flow, see `docs/architecture.md` and `docs/data-flow.md`.

## 🏗️ Architecture & Project Layout

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

## 🚀 Quick Start

The full setup guide lives in `docs/setup.md`. The ultra‑short version:

```bash
git clone https://github.com/abhizaik/SafeSurf.git
cd SafeSurf

# Option 1: Docker (recommended)
make dev-up

# Option 2: Local Go + Svelte
cd server && go run ./cmd/safesurf      # backend on :8080
cd ../web/website && npm install && npm run dev   # UI on :5173
```

Then:
```bash
curl "http://localhost:8080/api/v1/analyze?url=https://example.com"
```

Make sure `server/assets/top-1m.csv` contains a recent rank dataset; the backend loads it automatically on startup.

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

## 🧪 Testing & Quality

- **Backend unit tests**
  ```bash
  cd server
  go test ./...
  ```

- **Static analysis**
  ```bash
  make lint-backend   # go vet ./...
  ```

- **Frontend checks**
  ```bash
  cd web/website
  npm run check       # type checks (if configured)
  npm test            # when tests are added
  ```

See `docs/testing.md` for integration tests, load testing, and coverage tips.



## 🤝 Contributing

Bug reports and design discussions are welcome via GitHub Issues and Discussions.

1. Fork and create a feature branch (for example, `feat/tls-strength-signal`).
2. Implement your change with tests (`make test-backend`) and keep Go/Svelte code formatted.
3. Update or add docs in `docs/` so users understand the new behavior.
4. Open a PR with:
   - A short description of what changed.
   - Sample analyzer outputs (before/after) where relevant.
   - Any operational or security implications.




## Community

If SafeSurf is useful to you, **please consider starring the repository**,  it helps others discover the project.

You share analyzer outputs or use‑cases as GitHub issues to get them fixed or added as a feature.


Thank you for helping make the web a little safer. 



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
