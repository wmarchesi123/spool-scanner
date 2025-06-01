# Quick Start Guide

## 1. Prerequisites

- OctoPrint with TWO plugins:
  - [OctoPrint-Spoolman-API](https://github.com/wmarchesi123/octoprint-spoolman-api) - **Manual install required!**

    ```bash
    ~/oprint/bin/pip install https://github.com/wmarchesi123/octoprint-spoolman-api/archive/main.zip
    ```

  - [Spoolman plugin](https://plugins.octoprint.org/plugins/Spoolman/) - Install from plugin repository
- [Spoolman server](https://github.com/Donkie/Spoolman) running
- Docker installed (or Go 1.22+ for manual setup)

## 2. Get OctoPrint API Key

1. Open OctoPrint → Settings → Application Keys
2. Generate new key
3. Copy the key

## 3. Run with Docker

```bash
docker run -d \
  --name spool-scanner \
  -p 8080:8080 \
  -e SPOOLMAN_URL=http://your-spoolman:7912 \
  -e PRINTER_1_NAME="Your Printer" \
  -e PRINTER_1_URL=http://your-octoprint \
  -e PRINTER_1_KEY=your-api-key \
  ghcr.io/wmarchesi123/spool-scanner:latest
```

## 4. Create QR Codes

For each spool in Spoolman (check IDs in Spoolman UI):

- Generate QR code containing: `http://spool-scanner:8080/select/[SPOOL_ID]`
- Print on labels
- Attach to spools

## 5. Start Scanning

1. 📱 Scan QR code with phone
2. 🖨️ Select printer
3. ✅ Confirm
4. 🎯 Spool is now assigned!

## Example QR Code Content

```txt
http://192.168.1.100:8080/select/1   # For spool ID 1
http://192.168.1.100:8080/select/2   # For spool ID 2
```

Replace IP with your server's address.

## Troubleshooting

- **Can't connect**: Check SPOOLMAN_URL and PRINTER_1_URL are accessible
- **Assignment fails**: Verify API key permissions include "Plugin:SpoolMan"
- **"Plugin not found" error**: The OctoPrint-Spoolman-API plugin is not installed - see step 1
- **Spool not found**: Ensure spool ID exists in Spoolman

## Need More Help?

See the full [README](README.md) for detailed configuration options.
