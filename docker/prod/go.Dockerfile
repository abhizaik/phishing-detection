# ============================
# 1) BUILD STAGE
# ============================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git — required for go mod download (some deps)
RUN apk add --no-cache git

# Copy go.mod + go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy full source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o safesurf ./cmd/safesurf


# ============================
# 2) RUNTIME STAGE
# ============================
FROM alpine:latest

# Install Chromium & dependencies for headless browser
RUN apk add --no-cache \
    chromium \
    chromium-chromedriver \
    nss \
    freetype \
    harfbuzz \
    ttf-dejavu \
    ca-certificates \
    tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /app/safesurf .

# Copy assets folder
COPY --from=builder /app/assets ./assets

ENV PORT=8080
EXPOSE 8080

# Chromium default path used by many libraries (Playwright/Puppeteer style)
ENV CHROME_PATH=/usr/bin/chromium-browser

CMD ["./safesurf"]
