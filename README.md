# ntailfy

A Go service that monitors your Tailscale tailnet and sends notifications to ntfy when machines connect or disconnect.

## Features

- 🔍 Monitors Tailscale device state changes via API polling
- 📬 Sends authenticated notifications to ntfy (token-based auth)
- ⚙️ Configurable polling interval (default: 60s, minimum: 10s)
- 🛡️ Graceful shutdown support
- 🔐 Environment-based configuration
- 💓 Heartbeat logging every 10 polls to confirm operation

## Setup

### Prerequisites

- Go 1.23 or later (or Docker)
- Tailscale API key ([create one here](https://login.tailscale.com/admin/settings/keys))
- Self-hosted ntfy instance
- ntfy auth token ([create in ntfy web UI](https://docs.ntfy.sh/publish/#access-tokens))

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
| `DEVICE_FILTER` | Comma-separated device hostnames to monitor | No | `server1,server2,server3` |

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

1. Polls the Tailscale API at the configured interval (default: 60s)
2. Determines device online status based on `clientConnectivity.connectedToControl` field
3. Compares current device states with previous poll
4. Sends a notification to ntfy when a device state changes (connected/disconnected)
5. Logs a heartbeat every 10 polls showing device counts
6. Continues monitoring until interrupted (Ctrl+C)

**Note:** On first run, all devices are discovered but no notifications are sent. Notifications only occur on subsequent state changes.

## Notes

### Tailscale API
Tailscale supports webhooks for device events on paid plans. This tool uses polling which works on all plan tiers (free included). The API is rate-limited, so keep your polling interval reasonable (60s recommended).

### Device Online Status
Devices are considered "online" if they are actively connected to the Tailscale control plane (`clientConnectivity.connectedToControl` is true). This directly reflects the device's real-time connection status.

### Authentication
This service uses ntfy auth tokens (not username/password). Create a token in your ntfy instance's web UI under Account → Access Tokens.
