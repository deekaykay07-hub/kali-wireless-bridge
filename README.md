# kali-wireless-bridge

**Local Hardware Bridge** for the [Kali Wireless Expert TUI](https://github.com/deekaykay07-hub/kali-mistral-tui).

This small agent allows the remote AI-powered TUI (running on a VPS) to securely execute wireless commands and access physical hardware on *your* local machine (WiFi adapters in monitor mode, Bluetooth, NFC readers, SDR, etc.).

## Why This Exists

Wireless security work (monitor mode, packet injection, Bluetooth LE, NFC cloning, etc.) requires direct access to local radio hardware. Running everything in a remote container makes this impossible.

This bridge solves that by letting the smart remote AI plan and orchestrate attacks, while actual execution happens on your machine where the hardware lives.

## How It Works

1. You run the bridge agent on your laptop.
2. The bridge connects outbound to your remote TUI instance.
3. The AI in the TUI can now delegate commands to your local hardware.
4. Output is streamed back in real time.

## Features (Planned)

- Secure WebSocket connection with token auth
- Automatic detection of wireless interfaces
- Full command execution with real-time streaming output
- Support for privileged operations (monitor mode, injection)
- Cross-platform (Linux, macOS, Windows)

## Quick Start (Coming Soon)

```bash
./kali-bridge --connect your-droplet-ip:8765 --token <token-from-tui>
```

## Project Structure

```
kali-wireless-bridge/
├── cmd/bridge/main.go
├── internal/
│   ├── bridge/          # Connection & protocol handling
│   └── hardware/        # Local wireless interface detection & execution
├── internal/protocol/ # Message types
└── Makefile
```

## Development

```bash
go run cmd/bridge/main.go --help
```

## License

MIT
