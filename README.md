# kali-wireless-bridge

**Local Hardware Bridge** for the [Kali Wireless Expert TUI](https://github.com/deekaykay07-hub/kali-mistral-tui).

This small agent allows the remote AI-powered TUI (running on a VPS) to securely execute wireless commands and access physical hardware on *your* local machine (WiFi adapters in monitor mode, Bluetooth, NFC readers, SDR, etc.).

## Why This Exists

Wireless security work (monitor mode, packet injection, Bluetooth LE, NFC cloning, etc.) requires direct access to local radio hardware. Running everything in a remote container makes this impossible.

This bridge solves that by letting the smart remote AI plan and orchestrate attacks, while actual execution happens on your machine where the hardware lives.

## System Architecture & Flow

### High-Level Overview

- **AI Brain (Remote)**: Runs on your VPS in the `kali-mistral-tui` container. This is where the LLM (Mistral / Qwen2.5 etc.) lives and does all the thinking and planning.
- **Bridge Agent (Local)**: A small program you run on your laptop. It has direct access to your physical wireless hardware.
- **Communication**: The local bridge connects outbound to the remote TUI over WebSocket.

You do **NOT** need to run any AI/Ollama on your local machine.

### Detailed Flow (When You Give a Command)

1. You open the TUI in your browser (`http://your-vps-ip:7681`).
2. You type a request, e.g.:
   > "Put my WiFi card into monitor mode and scan for networks"
3. The remote AI (in the container) receives your message and creates a plan.
4. When the plan requires real hardware (monitor mode, packet injection, Bluetooth scan, etc.), the AI sends a command over the WebSocket to your local bridge.
5. The bridge on your laptop receives the command, executes it locally using your actual hardware, and streams the output back in real time.
6. The remote AI sees the results and continues the plan (or asks you for next steps).

This gives you the best of both worlds:
- Powerful AI + good models on the VPS (where resources are better)
- Real hardware access on your laptop

## How to Use (Planned Flow)

```bash
# On your laptop
./kali-bridge --connect your-vps-ip:8765 --token <token-from-tui>
```

The TUI will show when a bridge is connected and will automatically route hardware-dependent commands to it.

## Current Status

This project is in early development. The core connection skeleton exists, but full command execution, hardware detection, and integration with the TUI are still being built.

See the `internal/` folders for the current structure.

## Project Structure

```
kali-wireless-bridge/
├── cmd/bridge/main.go
├── internal/
│   ├── bridge/          # WebSocket connection & protocol
│   └── hardware/        # Local wireless detection & execution
├── internal/protocol/ # Message types
└── Makefile
```

## License

MIT
