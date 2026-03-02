# Agent Notes: Build + Push a `:develop` Docker image (ntailfy)

## Goal
Build and push a Docker image tagged `:develop` so it can be pulled on another machine for testing.

## Preconditions
- Docker Desktop / Docker Engine installed.
- `docker buildx` available.
- You are logged in to the target registry (`docker login`).

Docs:
- Docker buildx build: https://docs.docker.com/engine/reference/commandline/buildx_build/
- Multi-platform builds: https://docs.docker.com/build/building/multi-platform/

## Image names used for this repo
- Git remote: `https://github.com/leukosaima/ntailfy.git`
- Working registry in practice: Docker Hub
  - `leukosaima/ntailfy:develop`
- Attempted registry (failed due to permissions on this machine): GHCR
  - `ghcr.io/leukosaima/ntailfy:develop` failed with `permission_denied: write_package`

## Recommended command (multi-arch + push)
From the repo root:

PowerShell:
```powershell
docker buildx build --platform linux/amd64,linux/arm64 --push -t leukosaima/ntailfy:develop .
```

Bash:
```bash
docker buildx build --platform linux/amd64,linux/arm64 --push -t leukosaima/ntailfy:develop .
```

## Verify on the other machine
```bash
docker pull leukosaima/ntailfy:develop
```

## If pushing to GHCR instead
1. Ensure you can authenticate to GHCR (`docker login ghcr.io`).
2. Ensure the credential you use has permission to write packages (typically a PAT with `write:packages`).
3. Then run:

```bash
docker buildx build --platform linux/amd64,linux/arm64 --push -t ghcr.io/leukosaima/ntailfy:develop .
```
