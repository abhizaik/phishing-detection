<div align="center">

# Phishing Detection 
**A fast phishing detection engine with a web UI, API, and browser extension.**


[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Svelte](https://img.shields.io/badge/Svelte-5-orange?logo=svelte)](https://svelte.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/abhizaik/phishing-detection?style=social)](https://github.com/abhizaik/phishing-detection)
![Contributors](https://img.shields.io/github/contributors/abhizaik/phishing-detection)



[⚡ Quick Start](#-quick-start) ·
[📚 Docs](#-documentation) ·
[🏛 Architecture](#-architecture) ·
[🤝 Contributing](#-contributing) ·
[🌍 Community](#-community)


</div>


## Demo


> Paste a URL → get a verdict, trust score, and explanation in under a second.

![Phishing Analysis Demo](assets/demo.gif)




## Why this project?
Phishing is still one of the most effective attack vectors.  
Most tools today are either **half-baked, opaque, slow, or locked behind expensive APIs**.

This project gives you:
- **Transparent** detection logic
- **Fast**, parallel analysis
- Multiple ways to consume results **(UI, API, extension)**
- Full open source **control**

## What it does

- Analyzes URLs for **phishing indicators** and **malicious redirects**
- Runs multiple analyzers in parallel for **low latency** results
- Produces a **clear verdict**, **trust score**, and **detailed report**
- Designed for both **non-technical** users and **developers**
- Built with **Go** and **Svelte**


## Quick Start

Full setup: [docs/setup.md](docs/setup.md) 

1. Clone the repo

```bash
git clone https://github.com/abhizaik/phishing-detection.git
cd phishing-detection
```
2. Start the application

**Option 1: Docker (recommended)**
```bash
make build
make up
```
Web UI: **[localhost:3000](http://localhost:3000)** 

**Option 2: Local Go + Svelte**

Requires Go and Node.js.
```bash
cd server && go run ./cmd/safesurf      # backend on :8080
cd ../web/website && npm install && npm run dev   # UI on :5173
```
Web UI: **[localhost:5173](http://localhost:5173)** 



##  Documentation

All documentation lives under `docs/`. Start here [docs/README.md](docs/README.md) 




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
