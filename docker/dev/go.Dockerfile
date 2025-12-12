# Development image with Air + Chrome for screenshots
FROM golang:1.24-bullseye

# Install chrome deps and small utilities
RUN apt-get update \
  && apt-get install -y wget gnupg ca-certificates unzip \
     libnss3 libatk1.0-0 libcups2 libxdamage1 libxcomposite1 \
     libxrandr2 libgbm1 libasound2 libxshmfence1 libxkbcommon0 \
     libgtk-3-0 libdrm2 xvfb procps \
  && rm -rf /var/lib/apt/lists/*

# Install Google Chrome stable
RUN wget -q -O /tmp/chrome.deb https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
  && apt-get update && apt-get install -y /tmp/chrome.deb || true \
  && rm -f /tmp/chrome.deb \
  && rm -rf /var/lib/apt/lists/*

# Install Air
# RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/cosmtrek/air@v1.40.4


ENV PATH=$PATH:/go/bin

WORKDIR /app/server

# Copy go.mod for dependency caching during image build (fast)
COPY go.mod go.sum ./
RUN go mod download

# Copy everything (so build-time tools work); mounted at runtime by compose
COPY . .

# Expose dev port
EXPOSE 8080

# Run Air. It will use /app/.air.toml (we mount server/.air.toml from host)
CMD ["air", "-c", ".air.toml"]
