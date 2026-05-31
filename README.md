# kali-wireless-bridge

**Local Hardware Bridge** for the [Kali Wireless Expert TUI](https://github.com/deekaykay07-hub/kali-mistral-tui).

This small agent lets the remote AI (running on Colab T4 or a VPS) execute real commands and access wireless hardware on **your local Windows machine**.

## Why This Exists

The smart AI lives remotely. Real wireless attacks (monitor mode, Bluetooth, etc.) require physical hardware on *your* computer. This bridge solves that.

## Key Features (Current)

- One-command updater for Windows users (no manual Go work needed)
- Automatic hardware inventory push (WiFi + Bluetooth devices) to the AI on connect
- The model now knows your exact adapter names instead of guessing
- `/refresh-devices` command in the TUI to force a fresh scan
- Token reuse support (you don't need a new token every time)

## Quick Start (Windows)

### Option 1: One-Command Update (Recommended)

Run this in **PowerShell as Administrator**:

```powershell
irm https://raw.githubusercontent.com/deekaykay07-hub/kali-wireless-bridge/main/Update-Bridge.ps1 | iex
```

Then run the bridge with your token:

```powershell
.\kali-bridge-new.exe --connect 188.166.150.41:8765 --token YOUR_TOKEN
```

### Option 2: Manual Build

```powershell
git clone https://github.com/deekaykay07-hub/kali-wireless-bridge.git
cd kali-wireless-bridge
go build -o kali-bridge.exe ./cmd/bridge
```

## Getting a Token

In the TUI, type `/bridge-token` (it will reuse your existing token if one was previously saved).

Use `/show-token` to display your current token anytime.

## How It Works

1. You run the bridge on Windows.
2. It connects to the remote TUI.
3. On connect, it automatically sends your current WiFi + Bluetooth devices.
4. The AI can now see real device names and request commands to be run on your hardware.
5. Output streams back live to the TUI.

## Repository Structure

- `Update-Bridge.ps1` — One-command updater for Windows users
- `cmd/bridge/main.go` — The actual bridge (includes automatic device inventory)
- `INVENTORY_UPDATE.md` — Technical details about the hardware push feature
- `README-ONE-COMMAND.md` — Detailed one-command instructions

## Related Projects

- [Kali Wireless TUI](https://github.com/deekaykay07-hub/kali-mistral-tui) — The main interface

## License

MIT
