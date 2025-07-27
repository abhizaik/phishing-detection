# Contributing to SafeSurf

Thank you for your interest in contributing! 
We welcome bug reports, feature requests, code, documentation, and testing help.


## Project Structure

- `server/go/` – Backend (Golang)
- `web/` – Frontend (HTML/CSS/JS or framework)
- `docker/` – Dockerfiles & Compose configs
- `docs/` – Project documentation


##  How to Contribute

1. Fork the repo
2. Clone your fork and create a branch:
   ```bash
   git checkout -b your-feature-name
    ```

3. Make your changes
4. Run tests:

   ```bash
   make test
   ```
5. Push to your fork and open a Pull Request


##  Development Setup

### Backend (Go)

```bash
cd server/go
go mod tidy
go run cmd/safesurf/safesurf.go
```

### Frontend

Serve `web/` using your preferred method (e.g., `live-server`, Python HTTP server, or static hosting).

### Docker (Full Stack)

```bash
docker-compose -f docker/docker-compose.dev.yml up --build
```


##  Code Style & Tools

* Format Go code:

  ```bash
  gofmt -w .
  golint ./...
  ```
* Write clear commit messages:

  ```
  feat(auth): add token expiration check
  fix(api): correct 404 response logic
  ```


## Reporting Bugs

1. Search existing issues first
2. If not found, open a new issue
3. Include:

   * Steps to reproduce
   * Logs or screenshots
   * Your environment (OS, browser, etc.)


<!-- ## Feature Requests

Use the [Feature Request template](../../issues/new?template=feature_request.yml)
Describe the use case and expected behavior. -->


## Code of Conduct

Please follow our [Code of Conduct](CODE_OF_CONDUCT.md).


##  Thanks

Your contributions make this project better — whether it’s code, feedback, or documentation.

