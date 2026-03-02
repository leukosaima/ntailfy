---
name: ntailfy-develop-image
description: Build and push a multi-arch Docker image tagged :develop for the ntailfy repo (Docker Hub or GHCR). Use when asked to publish a develop/test image so another machine can docker pull it.
---

# Build + push `:develop`

## Inputs to confirm
- Registry + image name to push (defaults for this repo: `leukosaima/ntailfy:develop`).
- Whether multi-arch is needed (default: `linux/amd64,linux/arm64`).

## Steps
1. Check Docker + buildx are available:
   - `docker version`
   - `docker buildx version`
2. Ensure login to the target registry:
   - Docker Hub: `docker login`
   - GHCR: `docker login ghcr.io` (requires credentials with permission to write packages)
3. Build and push (multi-arch):

```bash
docker buildx build --platform linux/amd64,linux/arm64 --push -t leukosaima/ntailfy:develop .
```

## GHCR variant
```bash
docker buildx build --platform linux/amd64,linux/arm64 --push -t ghcr.io/leukosaima/ntailfy:develop .
```

## Verification (on the target machine)
```bash
docker pull leukosaima/ntailfy:develop
```
