
#  Setup Guide 

This document explains all commands available via the `Makefile` and alternative Windows commands (you can use make commands in windows if you have WSL).
You can use **dev** (development environment) or **app** (production / non-dev) workflows.

---

##  1. Show Help

See all available commands:

**Linux / macOS / WSL**

```bash
make help
```

---

##  2. App (Production / Non-Dev)

These commands use the **production Docker stack**. This is what you should use for deployment or to run a stable, reproducible version locally.

| Command      | Description                                              |
| ------------ | -------------------------------------------------------- |
| `make build` | Build all production Docker images (backend + frontend). |
| `make up`    | Start production stack (detached).                       |
| `make down`  | Stop production stack.                                   |
| `make logs`  | Tail production container logs.                          |

**Windows Equivalent**

```powershell
docker compose -f docker/prod/docker-compose.prod.yml build
docker compose -f docker/prod/docker-compose.prod.yml up -d
docker compose -f docker/prod/docker-compose.prod.yml down
docker compose -f docker/prod/docker-compose.prod.yml logs -f
```

---

##  3. Dev (Development Environment)

These commands use **docker/dev/docker-compose.dev.yml**, mount source directories, and allow hot reload.

| Command            | Description                    |
| ------------------ | ------------------------------ |
| `make dev-build`   | Build dev Docker images only.  |
| `make dev-up`      | Start dev environment.         |
| `make dev-down`    | Stop dev environment.          |
| `make dev-restart` | Restart dev environment.       |
| `make dev-logs`    | Tail dev logs.                 |
| `make ps`          | Show running dev containers.   |
| `make sh-backend`  | Enter backend container shell. |

**Windows Equivalent**

```powershell
docker compose -f docker/dev/docker-compose.dev.yml build
docker compose -f docker/dev/docker-compose.dev.yml up -d
docker compose -f docker/dev/docker-compose.dev.yml down
docker compose -f docker/dev/docker-compose.dev.yml logs -f
docker compose -f docker/dev/docker-compose.dev.yml ps
docker compose -f docker/dev/docker-compose.dev.yml exec backend sh
```

---

##  4. Local / Non-Docker (Linux / macOS / WSL)

These commands allow running the backend and frontend directly, without Docker. Useful for fast iteration.

| Command                     | Description                                       |
| --------------------------- | ------------------------------------------------- |
| `make local-build-backend`  | Build backend binary locally (`server/safesurf`). |
| `make local-build-frontend` | Build frontend locally (`web/website`).           |
| `make local-run-backend`    | Run backend via Air (hot reload).                 |
| `make local-run-frontend`   | Run Svelte dev server.                            |

**Windows Equivalent**

```powershell
# Backend
cd server
go build -o safesurf ./cmd/safesurf
air

# Frontend
cd web/website
npm install
npm run build
npm run dev
```

---

##  5. Testing, Linting & CI

| Command           | Description                                         |
| ----------------- | --------------------------------------------------- |
| `make doctor`     | Verify environment: Docker, Compose, Go installed.  |
| `make fmt`        | Format Go code (`go fmt`).                          |
| `make tidy`       | Tidy Go modules (`go mod tidy`).                    |
| `make local-test` | Run backend tests locally.                          |
| `make local-lint` | Run Go vet on backend code.                         |
| `make ci`         | Run CI-equivalent checks: tidy + fmt + lint + test. |

**Windows Equivalent**

```powershell
# Lint / Test
cd server
go fmt ./...
go mod tidy
go vet ./...
go test ./...
```

---

##  6. Cleanup

Remove stopped containers, volumes, and prune unused Docker resources.

**Linux / macOS / WSL**

```bash
make clean
```

**Windows**

```powershell
docker compose -f docker/dev/docker-compose.dev.yml down -v
docker compose -f docker/prod/docker-compose.prod.yml down -v
docker system prune -f
```

---

##  7. Notes

* **Dev vs App**:

  * `dev-*` commands → development environment with source mounts, hot reload, Chrome, Air, etc.
  * `build`, `up`, `down`, `logs` → production / app stack.

* **Windows users**: `make` is not supported natively. Use the equivalent `docker compose` commands directly.

* **Timers & Feedback**: All `make` commands show build/start timers and colored messages for convenience.


