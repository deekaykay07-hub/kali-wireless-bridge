# Architecture Overview

This document explains how the **Kali Wireless Bridge** system is designed to work.

## Core Philosophy

- The **intelligence** (LLM + planning) lives remotely on your VPS.
- The **execution** (especially anything requiring real wireless hardware) happens locally on your machine.
- The bridge is a thin, secure agent that only does what the remote AI tells it to.

## Components

### 1. Remote TUI (kali-mistral-tui)

- Runs in Docker on your VPS.
- Contains the LLM (via Ollama) + the Textual TUI.
- This is where you interact with the AI.
- It decides when a task needs local hardware.
- It will (in the future) expose a WebSocket endpoint for bridges to connect to.

### 2. Local Bridge (this repo)

- A small Go binary you run on your laptop.
- Connects **outbound** to the remote TUI (firewall friendly).
- Advertises what hardware/capabilities it has.
- Receives commands from the remote AI.
- Executes them locally and streams results back.

### 3. Communication Protocol

- WebSocket (currently planned on port 8765 on the VPS side).
- JSON messages.
- Main message types:
  - `capabilities` (bridge tells TUI what it can do)
  - `command` (TUI tells bridge to run something)
  - `output` (streaming stdout/stderr from bridge)
  - `result` (final exit code)

## User Flow (End-to-End)

1. You start the bridge on your laptop:
   ```bash
   ./kali-bridge --connect 188.166.150.41:8765 --token abc123
   ```

2. The bridge connects to the remote TUI and registers itself.

3. In the TUI (browser), you type a command that needs hardware:
   > "Scan for Bluetooth devices using my local adapter"

4. The remote AI creates a plan and realizes it needs local execution.

5. The AI sends a `command` message over the WebSocket to your bridge.

6. Your bridge runs the appropriate local command (e.g. `hcitool scan` or bettercap equivalent).

7. Output is streamed back in real time to the TUI.

8. The AI can now continue with the next step using the real data from your hardware.

## Security Model

- Bridge only makes outbound connections.
- Token-based authentication.
- All commands from the AI are visible to you in the TUI before execution (future feature).
- The bridge never initiates connections back to your VPS except for the control channel.

## Current State (as of May 2026)

- Basic connection + capabilities skeleton exists.
- Full command execution, hardware modules, and TUI integration are still in development.

## Future Enhancements

- Automatic hardware discovery on the bridge side.
- Secure command allow-list on the bridge.
- Better UI in the TUI showing "Local hardware bridge connected" status.
- Support for multiple bridges (e.g. one on laptop, one on a dedicated WiFi rig).
