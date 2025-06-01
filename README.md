# Spool Scanner

A web-based system for managing 3D printer filament spools using QR codes or NFC tags. Scan a spool's tag, select which printer to load it into, and the system automatically updates OctoPrint with the correct spool information via the Spoolman integration.

> ⚠️ **Important**: This project requires the [OctoPrint-Spoolman-API](https://github.com/wmarchesi123/octoprint-spoolman-api) plugin to be manually installed on each OctoPrint instance. This plugin is not yet available in the OctoPrint plugin repository.

## Features

- 📱 **Mobile-friendly** web interface optimized for phones/tablets
- 🏷️ **QR Code/NFC Support** - Scan tags to instantly identify spools
- 🖨️ **Multi-printer Support** - Manage up to 10 printers
- 🔄 **Real-time Status** - See which spools are loaded in which printers
- 🎯 **Direct Integration** - Works with OctoPrint's Spoolman plugin
- 🚀 **Fast Assignment** - Load spools with just two taps
- 🐳 **Docker Ready** - Easy deployment with included Dockerfile

## How It Works

1. **Tag Your Spools**: Create QR codes or program NFC tags with URLs like:

   ```txt
   http://spool-scanner/select/[SPOOL_ID]
   ```

   Where `[SPOOL_ID]` matches the spool ID in Spoolman

2. **Scan & Select**: When loading filament:
   - Scan the spool's QR code/NFC tag with your phone
   - Select which printer you're loading it into
   - Confirm the assignment

3. **Automatic Updates**: The system updates OctoPrint immediately, ensuring accurate filament tracking

## System Architecture

```txt
┌─────────────┐     ┌──────────────┐     ┌─────────────────────┐
│   Phone     │────▶│ Spool Scanner│────▶│ OctoPrint Instance  │
│ (QR/NFC)    │     │  Web Server  │     │ + Spoolman-API      │
└─────────────┘     └──────┬───────┘     │   Plugin            │
                           │             └──────────┬──────────┘
                           ▼                        │
                    ┌─────────────┐                 │
                    │  Spoolman   │◀────────────────┘
                    │   Server    │
                    └─────────────┘
```

- **Spool Scanner** acts as the user interface for spool assignment
- **OctoPrint-Spoolman-API** plugin provides the API endpoints for spool management
- **Spoolman** maintains the central database of all spools

## Prerequisites

### Required OctoPrint Plugins

Each OctoPrint instance needs **TWO** plugins:

#### 1. OctoPrint-Spoolman-API (REQUIRED - Manual Install)

- **Repository**: [OctoPrint-Spoolman-API](https://github.com/wmarchesi123/octoprint-spoolman-api)
- **Status**: Not yet in OctoPrint plugin repository (manual install required)
- **Purpose**: Provides API endpoints for spool-scanner to update spool assignments
- **Installation**: See detailed instructions below
- **Future**: Will be submitted to the official plugin repository

#### 2. Spoolman Plugin

- **Repository**: [Spoolman](https://github.com/Donkie/Spoolman-Octoprint)
- **Status**: Available in OctoPrint plugin repository
- **Purpose**: Integrates OctoPrint with Spoolman server
- **Installation**: Via Plugin Manager

### Other Requirements

- **Spoolman** server for filament tracking
- **QR codes or NFC tags** on your spools (can use a label printer)

## Installation

### Prerequisites Setup

#### 1. Install OctoPrint-Spoolman-API Plugin (Required)

This plugin must be installed **manually** on each OctoPrint instance:

```bash
# SSH into your OctoPrint instance, then:
~/oprint/bin/pip install https://github.com/wmarchesi123/octoprint-spoolman-api/archive/main.zip
```

Or via OctoPrint web interface:

1. Settings → Plugin Manager → Get More
2. Enter URL: `https://github.com/wmarchesi123/octoprint-spoolman-api/archive/main.zip`
3. Click Install

#### 2. Install Spoolman Plugin

Install from OctoPrint plugin repository:

1. Settings → Plugin Manager → Get More
2. Search for "Spoolman"
3. Install

#### 3. Restart OctoPrint

After installing both plugins, restart OctoPrint for changes to take effect.

#### 4. Verify Installation

To verify the OctoPrint-Spoolman-API plugin is installed correctly:

1. Go to OctoPrint Settings → Plugin Manager
2. Look for "Spoolman API" in the installed plugins list
3. Or test with curl:

   ```bash
   curl -H "X-Api-Key: YOUR_API_KEY" \
        -H "Content-Type: application/json" \
        -d '{"command":"get_current_spool","tool":0}' \
        http://your-octoprint/api/plugin/spoolman_api
   ```

   You should get a JSON response (not a 404 error)

### What Happens Without OctoPrint-Spoolman-API?

Without this plugin installed:

- ❌ Spool assignments will fail with "API error" or 404 errors
- ❌ The test script will show the plugin is missing
- ❌ You'll see errors in the logs about `/api/plugin/spoolman_api` not found

This plugin is **essential** for spool-scanner to function - it cannot work without it.

### Quick Start

See [QUICKSTART.md](QUICKSTART.md) for a simplified setup guide.

### Using Docker (Recommended)

```bash
docker build -t spool-scanner .
docker run -d \
  --name spool-scanner \
  -p 8080:8080 \
  -e SPOOLMAN_URL=http://spoolman:7912 \
  -e PRINTER_1_NAME="Prusa MK4" \
  -e PRINTER_1_URL=http://octoprint1 \
  -e PRINTER_1_KEY=your-api-key \
  spool-scanner
```

### Manual Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/wmarchesi123/spool-scanner.git
   cd spool-scanner
   ```

2. **Install Go 1.22+** if not already installed

3. **Build and run**:

   ```bash
   go mod download
   go build -o spool-scanner cmd/server/main.go
   ./spool-scanner
   ```

## Configuration

Configure via environment variables. See [.env.example](.env.example) for a template.

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `SPOOLMAN_URL` | URL of your Spoolman server | `http://spoolman:7912` |
| `PRINTER_1_NAME` | Display name for printer 1 | `Prusa MK4` |
| `PRINTER_1_URL` | OctoPrint URL for printer 1 | `http://octoprint1` |
| `PRINTER_1_KEY` | OctoPrint API key for printer 1 | `ABCD1234...` |

### Additional Printers

Add more printers by incrementing the number (up to 10):

```bash
PRINTER_2_NAME="Ender 3 V2"
PRINTER_2_URL=http://octoprint2
PRINTER_2_KEY=your-api-key-2

PRINTER_3_NAME="Voron 2.4"
PRINTER_3_URL=http://octoprint3
PRINTER_3_KEY=your-api-key-3
```

### Getting OctoPrint API Keys

1. Open OctoPrint web interface
2. Go to Settings → Application Keys
3. Generate a new application key
4. Copy the key for use in configuration

## Creating QR Codes/NFC Tags

### QR Codes

1. **Get Spool IDs**: Check Spoolman web UI for spool IDs
2. **Generate QR Codes**: Use any QR code generator with URLs like:

   ```txt
   http://spool-scanner/select/1
   http://spool-scanner/select/2
   ```

3. **Print Labels**: Use a label printer to create physical tags

### NFC Tags

1. **Get NFC Tags**: NTAG213/215/216 tags work well
2. **Program Tags**: Use apps like NFC Tools (iPhone) or TagWriter (Android)
3. **Write URL**: Same format as QR codes
4. **Attach to Spools**: Stick tags on spool sides

## Usage Workflow

### Loading a Spool

1. **Remove old spool** from printer (if any)
2. **Scan new spool** with phone camera (QR) or tap (NFC)
3. **Select printer** from the list
4. **Confirm** assignment
5. **Load filament** physically into printer

### Viewing Current Assignments

Status page coming soon.

## Kubernetes Deployment

Example deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: spool-scanner
spec:
  replicas: 1
  selector:
    matchLabels:
      app: spool-scanner
  template:
    metadata:
      labels:
        app: spool-scanner
    spec:
      containers:
      - name: spool-scanner
        image: spool-scanner:latest
        ports:
        - containerPort: 8080
        env:
        - name: SPOOLMAN_URL
          value: "http://spoolman:7912"
        - name: PRINTER_1_NAME
          value: "Prusa MK4"
        - name: PRINTER_1_URL
          value: "http://octoprint-prusa"
        - name: PRINTER_1_KEY
          valueFrom:
            secretKeyRef:
              name: octoprint-keys
              key: prusa-key
---
apiVersion: v1
kind: Service
metadata:
  name: spool-scanner
spec:
  selector:
    app: spool-scanner
  ports:
  - port: 80
    targetPort: 8080
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Home page with instructions |
| `/select/{spool_id}` | GET | Spool selection page (for QR/NFC) |
| `/api/printers` | GET | List all printers and their status |
| `/api/spool/{id}` | GET | Get spool details from Spoolman |
| `/api/assign` | POST | Assign spool to printer (via OctoPrint-Spoolman-API plugin) |

Note: The `/api/assign` endpoint communicates with OctoPrint through the OctoPrint-Spoolman-API plugin's `/api/plugin/spoolman_api` endpoint.

## Troubleshooting

### Testing Your Setup

Use the included test script to verify your configuration:

```bash
chmod +x tools/test-setup.sh
./tools/test-setup.sh
```

This will check:

- Environment variables are set correctly
- Spoolman server is accessible
- OctoPrint instances are reachable
- Spoolman plugin is installed in OctoPrint
- API keys have proper permissions

### Common Issues

#### Spool Not Found

- Verify spool ID exists in Spoolman
- Check SPOOLMAN_URL is correct
- Ensure Spoolman server is accessible

#### Can't Connect to Printer

- Verify OctoPrint URL is correct
- Check API key has proper permissions
- Ensure OctoPrint has Spoolman plugin installed

#### Assignment Fails

- Check printer has Spoolman plugin enabled
- Verify API key has "Plugin:SpoolMan" permission
- Look at container logs for detailed errors

## License

Apache License Version 2.0, see [LICENSE.md](LICENSE.md) for details.

## FAQ

### Why do I need two OctoPrint plugins?

1. **Spoolman Plugin**: Handles the integration between OctoPrint and Spoolman server for tracking filament usage
2. **OctoPrint-Spoolman-API Plugin**: Provides the API endpoints that spool-scanner uses to update which spool is loaded

Without both plugins, the system cannot function properly.

### When will OctoPrint-Spoolman-API be in the plugin repository?

The plugin is planned to be submitted to the official OctoPrint plugin repository. Once approved, installation will be much simpler through the Plugin Manager.

### Can I use this without Spoolman?

No, Spoolman is the central database that tracks all your spools. Both OctoPrint and spool-scanner rely on it.

### Does this work with Klipper/Moonraker?

Currently, this only supports OctoPrint. Klipper support could be added in the future.

## Security Considerations

- **API Keys**: Store OctoPrint API keys securely (use Kubernetes secrets, environment files, or secret management systems)
- **Network**: Run on a trusted network - the system has no built-in authentication
- **HTTPS**: Use a reverse proxy with HTTPS for production deployments
- **Permissions**: OctoPrint API keys only need "Plugin:SpoolMan" permission

## Related Projects

- [OctoPrint-Spoolman-API](https://github.com/wmarchesi123/octoprint-spoolman-api) - **Required OctoPrint plugin** for spool-scanner functionality
- [OctoDash](https://github.com/wmarchesi123/octodash) - Real-time dashboard for monitoring multiple 3D printers
- [Spoolman](https://github.com/Donkie/Spoolman) - Filament spool management server
- [OctoPrint](https://octoprint.org/) - The snappy web interface for your 3D printer
