# Development image with Air
FROM golang:1.24-bullseye

# Install basic utilities
RUN apt-get update \
  && apt-get install -y wget gnupg ca-certificates unzip procps \
  && rm -rf /var/lib/apt/lists/*

# Install Air
RUN go install github.com/cosmtrek/air@v1.40.4

ENV PATH=$PATH:/go/bin

WORKDIR /app

# Copy go.mod for dependency caching during image build (fast)
COPY go.mod go.sum ./
RUN go mod download

# Copy everything (so build-time tools work); mounted at runtime by compose
COPY . .

# Expose dev port
EXPOSE 8080

# Run Air. It will use /app/.air.toml (we mount server/.air.toml from host)
CMD ["air", "-c", ".air.toml"]
