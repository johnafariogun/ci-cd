# CI/CD Folder

This folder demonstrates a complete **CI/CD pipeline** for a Go application using **GitHub Actions** and **Docker**.

## 📋 Contents

- **`server.go`** — A simple Go HTTP server that listens on port 8080
- **`Dockerfile`** — Multi-stage Docker build file for containerizing the Go server
- **`.github/workflows/ci-cd.yaml`** — GitHub Actions workflow for automated testing and Docker image publishing

---

## 🚀 Overview

### What It Does

1. **Builds** the Go binary from `server.go`
2. **Tests** the running server locally using curl
3. **Containerizes** the application using Docker (multi-stage build)
4. **Publishes** the Docker image to GitHub Container Registry (GHCR)

### Workflow Triggers

The CI/CD pipeline runs on:
- **Push and Pull to `main` branch**

---

## 📝 Server Details

The `server.go` file implements a basic HTTP server:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, Go Web!, Ci-Cd is a go")
}
```

**Endpoint:** `GET http://localhost:8080/` → Returns message

---

## 🐳 Docker Build

The `Dockerfile` uses a **multi-stage build** for a minimal production image:

```dockerfile
# Stage 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server server.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

**Benefits:**
- Final image only includes the binary (smaller size)
- No Go compiler in production image (smaller attack surface)

---

## ⚙️ GitHub Actions Workflow

### Job 1: Build & Test

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - Checkout code
      - Set up Go 1.22
      - Build binary
      - Start server and test with curl
      - Kill server
```

**Purpose:** Validates that the code builds and the server responds correctly.

### Job 2: Docker Build & Push

```yaml
jobs:
  docker:
    needs: build  # Runs only after build succeeds
    runs-on: ubuntu-latest
    steps:
      - Checkout code
      - Login to GHCR
      - Build Docker image
      - Push to ghcr.io
```

**Purpose:** Containerizes the application and publishes it to GitHub Container Registry.

---

## 🔑 Configuration

### Required GitHub Secret

To push to GHCR, you need a **personal access token** stored as a GitHub secret:

**Secret Name:** `SUPER_SECRET_GITHUB_TOKEN`

**Setup Steps:**
1. Go to GitHub → Settings → Developer settings → Personal access tokens
2. Generate a token with `write:packages` scope
3. Add it to your repo → Settings → Secrets → `SUPER_SECRET_GITHUB_TOKEN`

Alternatively, use the default `GITHUB_TOKEN` (available in all repos):

```yaml
password: ${{ secrets.GITHUB_TOKEN }}
```

---

## 🏃 Running Locally

### Build the server:

```bash
cd ci-cd
go build -o server server.go
```

### Run the server:

```bash
./server
```

Expected output:
```
Server starting on port 8080...
```

### Test the endpoint:

```bash
curl http://localhost:8080/
# Output: Hello, Go Web!, Ci-Cd is a go
```

---

## 🐳 Running with Docker

### Build the image:

```bash
docker build -t go-ci-cd-app .
```

### Run the container:

```bash
docker run -p 8080:8080 go-ci-cd-app
```

### Test:

```bash
curl http://localhost:8080/
```

---

## 📦 GitHub Container Registry

### Pull the published image:

```bash
docker pull ghcr.io/johnafariogun/go-ci-cd-app:latest
```

### Run it:

```bash
docker run -p 8080:8080 ghcr.io/johnafariogun/go-ci-cd-app:latest
```

---

## 🔄 CI/CD Flow

```
Code Push to main
       ↓
GitHub Actions Triggered
       ↓
Build Job (build & test Go binary)
       ├─ If failed → Stop
       └─ If successful → Continue
       ↓
Docker Job (build & push image)
       ├─ Login to GHCR
       ├─ Build Docker image
       └─ Push to ghcr.io
       ↓
Image available at ghcr.io/<repo>/go-ci-cd-app:latest
```

---

## 🛠 Troubleshooting

### Docker push fails with "unauthorized"

- Verify `SUPER_SECRET_GITHUB_TOKEN` is set correctly in repo secrets
- Ensure token has `write:packages` scope

### Build job fails

- Check Go version in `ci-cd.yaml` matches your local version
- Verify no syntax errors in `server.go`

### Test curl fails

- Ensure port 8080 is not already in use
- Check firewall settings

---

## 💡 Next Steps

- Add more routes and handlers to `server.go`
- Add unit tests and run them in the workflow
- Deploy the Docker image to a cloud platform (ECS, K8s, etc.)
- Add health checks and metrics
- Implement blue-green or canary deployments

---

## 📚 Resources

- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Docker Multi-Stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [GHCR Publishing Guide](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
