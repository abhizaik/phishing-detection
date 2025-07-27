# Docker Setup

This directory contains all Docker related configuration files.

## Files

- **Dockerfile.dev**: Used for development.
- **Dockerfile.prod**: Production ready build with optimized layers.
- **docker-compose.dev.yml**: Compose file for local development.
- **docker-compose.prod.yml**: Compose file for production deployments.

## Usage

### Development

```bash
docker-compose -f docker-compose.dev.yml up --build
```

### Production

```bash
docker-compose -f docker-compose.prod.yml up --build -d
```