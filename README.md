# kali-wireless-bridge

**Local Hardware Bridge** for the [Kali Wireless Expert TUI](https://github.com/deekaykay07-hub/kali-mistral-tui).

This agent allows the remote AI (running on a VPS) to execute commands and access wireless hardware on **your local machine**.

**Current Priority: Windows 11 support first**, Linux support coming after.

## Why This Exists

The smart AI lives on a remote VPS. But real wireless work (monitor mode, Bluetooth, NFC, etc.) requires hardware that is physically plugged into *your* computer. This bridge bridges that gap.

## Current Status (Windows Focus)

- The bridge runs on Windows 11.
- It can connect to the remote TUI and execute commands.
- Basic hardware detection for Windows is in progress.
- Full monitor mode + injection on Windows is limited by the OS and drivers (we'll add best-effort support + clear warnings).

## How It Works (High Level)

1. You run the small `kali-bridge.exe` on your Windows laptop.
2. It connects outbound to your remote TUI.
3. When the AI needs to use your local WiFi/Bluetooth hardware, it sends commands to the bridge.
4. The bridge runs them on your machine and streams the results back.

You do **not** need Ollama or any AI running locally.

## Building on Windows (Easy)

### Requirements
- Go 1.22+ installed (https://go.dev/dl/)
- Git (optional but recommended)

### Build Steps (Command Prompt or PowerShell)

```powershell
# 1. Clone the repo
git clone https://github.com/deekaykay07-hub/kali-wireless-bridge.git
cd kali-wireless-bridge

# 2. Build for Windows
go build -o kali-bridge.exe ./cmd/bridge

# Or cross-compile if needed
$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -o kali-bridge.exe ./cmd/bridge
```

This will create `kali-bridge.exe`.

## Running on Windows

```powershell
# Basic run (replace with your VPS IP and token)
.\kali-bridge.exe --connect 188.166.150.41:8765 --token YOUR_TOKEN

# Run as Administrator if you need to do privileged wireless operations
# Right-click PowerShell → "Run as Administrator" then run the command above
```

## Getting a Token

In the remote TUI (browser), there will be a way to generate a connection token for your bridge (planned feature: `/bridge-token` or a button in settings).

## Next Steps (Being Built)

- Better Windows wireless interface detection
- Streaming command execution (already partially implemented)
- Status indicator in the TUI showing when your bridge is connected
- Linux version (after solid Windows support)

## License

MIT
