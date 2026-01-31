# tsctl - Tailscale Controller

Self-contained HTTP API for Tailscale with embedded daemon using tsnet.

## Features

- **Single Binary** - No external dependencies (68MB)
- **Embedded Daemon** - Built-in Tailscale daemon via tsnet
- **REST API** - Control Tailscale via HTTP endpoints
- **Serve & Funnel** - Easy setup for exposing services
- **SSH Control** - Enable/disable SSH access
- **Simple Deployment** - Just copy and run

## Architecture

tsctl embeds the entire Tailscale daemon using the tsnet package, eliminating the need for separate `tailscaled` binary.

```
┌──────────────────────────────────────┐
│         Single Binary: tsctl         │
│  ┌────────────┐    ┌──────────────┐  │
│  │ HTTP API   │───▶│ tsnet.Server │  │
│  └────────────┘    └──────────────┘  │
└──────────────────────────────────────┘
```

## Installation

### Build from Source

```bash
git clone <repo>
cd tsctl
go mod tidy
go build -o bin/tsctl ./cmd/server
```

### Binary Size

- **tsctl**: ~68MB (includes full Tailscale stack)

## Configuration

Create `config/app.yaml`:

```yaml
tailscale:
  state_dir: ./tsnet-state
  hostname: tsctl
  auth_key: ""
  ephemeral: false

logging:
  app_log_path: ./logs/app.log
  level: info
  format: json

server:
  port: 8080
  host: localhost
  mode: release
```

### Config Fields

**tailscale:**
- `state_dir`: Directory for tsnet state storage
- `hostname`: Node hostname on tailnet
- `auth_key`: Optional auth key for auto-login
- `ephemeral`: If true, node is removed on exit

**logging:**
- `app_log_path`: Log file path
- `level`: Log level (debug, info, warn, error)
- `format`: Log format (json, text)

**server:**
- `port`: HTTP server port
- `host`: HTTP server host
- `mode`: Gin mode (debug, release)

## Usage

### Start Server

```bash
./bin/tsctl -c config/app.yaml
```

**First Run (No Auth Key):**
```
tsctl starting (embedded mode)
starting embedded tailscale daemon
To authenticate, visit: https://login.tailscale.com/a/xxxxx
```

**With Auth Key:**
Auto-connects without browser.

### API Endpoints

Base URL: `http://localhost:8080`

#### Status

```bash
GET /api/v1/status
```

Response:
```json
{
  "success": true,
  "data": "..." 
}
```

#### Authentication Status

```bash
GET /api/v1/auth/status
```

Response (Authenticated):
```json
{
  "success": true,
  "data": {
    "authenticated": true,
    "backend_state": "Running",
    "node_key": "nodekey:...",
    "hostname": "tsctl",
    "ips": ["100.x.y.z"]
  }
}
```

Response (Not Authenticated):
```json
{
  "success": true,
  "data": {
    "authenticated": false,
    "backend_state": "NeedsLogin",
    "auth_url": "https://login.tailscale.com/a/xxxxx"
  }
}
```

#### Logout

```bash
POST /api/v1/auth/logout
```

Response:
```json
{
  "success": true,
  "message": "logged out from tailnet - restart required to reconnect"
}
```

#### Start Serve

```bash
POST /api/v1/serve
Content-Type: application/json

{
  "port": 8080,
  "background": false
}
```

#### Start Funnel

```bash
POST /api/v1/funnel
Content-Type: application/json

{
  "port": 443,
  "background": false
}
```

#### Get Serve Status

```bash
GET /api/v1/serve/status
```

#### Reset Serve

```bash
DELETE /api/v1/serve
```

#### Enable SSH

```bash
POST /api/v1/ssh/enable
```

#### Get Logs

```bash
GET /api/v1/logs/app?lines=100
```

### Swagger Documentation

API documentation available at:
```
http://localhost:8080/docs/index.html
```

## Development

### Run Tests

```bash
go test ./...
```

### Build

```bash
make build
```

### Run in Debug Mode

Update `config/app.yaml`:
```yaml
server:
  mode: debug
```

## Deployment

### Single Binary Deployment

```bash
scp bin/tsctl server:/usr/local/bin/
scp config/app.yaml server:/etc/tsctl/config.yaml

ssh server
/usr/local/bin/tsctl -c /etc/tsctl/config.yaml
```

### Systemd Service

Create `/etc/systemd/system/tsctl.service`:

```ini
[Unit]
Description=Tailscale Controller
After=network.target

[Service]
Type=simple
User=tsctl
ExecStart=/usr/local/bin/tsctl -c /etc/tsctl/config.yaml
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable tsctl
sudo systemctl start tsctl
```

## Differences from v1

### v1 (External Daemon)
- Required `tailscaled` binary
- Required `tailscale` CLI
- Multiple processes
- Complex setup

### v2 (Embedded tsnet)
- Single binary
- No external dependencies
- One process
- Simple deployment

### Daemon Endpoints Behavior

- `POST /daemon/start` - Returns success (already running)
- `POST /daemon/stop` - Returns error (can't stop embedded)
- `POST /daemon/restart` - Returns error (restart app instead)
- `GET /daemon/status` - Returns embedded mode info

## Troubleshooting

### State Directory

tsnet stores state in configured `state_dir`. This contains:
- Machine key
- Node preferences
- Login state

**Backup this directory to preserve node identity.**

### Logs

Check application logs:
```bash
tail -f logs/app.log
```

### Port Already in Use

Change port in config:
```yaml
server:
  port: 9090
```

## License

BSD-3-Clause

## Credits

Built with [Tailscale](https://tailscale.com/) tsnet package.
