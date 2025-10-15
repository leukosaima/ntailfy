# ntailfy

A Go service that monitors your Tailscale tailnet and sends notifications to ntfy when machines connect or disconnect.

## Features

- 🔍 Monitors Tailscale device state changes
- 📬 Sends authenticated notifications to ntfy
- ⚙️ Configurable polling interval
- 🛡️ Graceful shutdown support
- 🔐 Environment-based configuration

## Setup

### Prerequisites

- Go 1.23 or later
- Tailscale API key ([create one here](https://login.tailscale.com/admin/settings/keys))
- Self-hosted ntfy instance with authentication

### Configuration

Set the following environment variables:

| Variable | Description | Required | Example |
|----------|-------------|----------|---------|
| `TAILSCALE_API_KEY` | Your Tailscale API key | Yes | `tskey-api-xxx` |
| `TAILSCALE_TAILNET` | Your tailnet name | Yes | `example.com` or `user@example.com` |
| `NTFY_URL` | Your ntfy instance URL | Yes | `https://ntfy.example.com` |
| `NTFY_AUTH_TOKEN` | ntfy auth token | No | `tk_xxxxxxxxxxxxx` |
| `NTFY_TOPIC` | ntfy topic to publish to | Yes | `tailscale-alerts` |
| `POLL_INTERVAL` | How often to check (min 10s) | No | `60s` (default) |

### Running

```bash
# Set environment variables
$env:TAILSCALE_API_KEY="your-api-key"
$env:TAILSCALE_TAILNET="your-tailnet"
$env:NTFY_URL="https://ntfy.example.com"
$env:NTFY_AUTH_TOKEN="your-auth-token"
$env:NTFY_TOPIC="tailscale"

# Build and run
go run .
```

Or build a binary:

```bash
go build -o ntailfy.exe
./ntailfy.exe
```

### Docker

Build and run with Docker:

```bash
# Build the image
docker build -t ntailfy .

# Run the container
docker run -d --name ntailfy \
  -e TAILSCALE_API_KEY="your-api-key" \
  -e TAILSCALE_TAILNET="your-tailnet" \
  -e NTFY_URL="https://ntfy.example.com" \
  -e NTFY_AUTH_TOKEN="your-auth-token" \
  -e NTFY_TOPIC="tailscale" \
  -e POLL_INTERVAL="60s" \
  ntailfy
```

Or use Docker Compose:

```bash
# Create .env file with your values (see .env.example)
# Then run:
docker-compose up -d
```

## How it works

1. Polls the Tailscale API at the configured interval
2. Compares current device states with previous poll
3. Sends a notification to ntfy when a device state changes
4. Continues monitoring until interrupted (Ctrl+C)

## Note on Tailscale API

Tailscale supports webhooks for device events, but they require a paid plan. This tool uses polling which works on all plan tiers. The API is rate-limited, so keep your polling interval reasonable (60s recommended).
