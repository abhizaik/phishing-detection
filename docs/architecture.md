# Architecture

## Overview
Brief description of what the system does, its main components, and what problems it solves.

## High-Level Diagram
(Embed or link to an image like system_architecture.png showing services, DB, external APIs, etc.)

## Components

### 1. Web Client (`/web`)
- HTML/CSS/JS frontend rendered with Go templates
- Handles user input and shows results

### 2. Backend Server (`/server/go`)
- Handles routing, validation, business logic
- Exposes REST endpoints
- Interacts with PostgreSQL and external APIs

### 3. Python Service (`/server/py_files`)
- Runs ML inference on phishing URLs
- Exposed via internal HTTP or CLI call

### 4. Database
- PostgreSQL
- Tables: `users`, `urls`, `scans`, etc.

## Data Flow
1. User submits URL on web
2. Frontend sends request to Go backend
3. Backend calls Python script for prediction
4. Result saved to DB and returned to user

## Deployment
- Dockerized via `docker-compose`
- Dev and Prod configs separated
- Deployed behind Nginx (reverse proxy)


