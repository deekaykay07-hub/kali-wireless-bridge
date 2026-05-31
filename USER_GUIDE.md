# User Guide – Windows 11 First

## Current Reality (as of now)

We are building this bridge with **Windows 11 support as the top priority**. Linux support will come after.

**Important Limitation on Windows:**
Real monitor mode + packet injection is much harder on Windows than on Linux. Many cheap adapters won't work well. The bridge will still let you run commands and do a lot of useful work (scanning, Bluetooth, etc.), but advanced injection may be limited.

## What You Need to Do on Your Windows 11 Laptop

### 1. Install Go (if you don't have it)

Download from: https://go.dev/dl/

### 2. Build the Bridge

Open PowerShell or Command Prompt and run:

```powershell
git clone https://github.com/deekaykay07-hub/kali-wireless-bridge.git
cd kali-wireless-bridge

go build -o kali-bridge.exe ./cmd/bridge
```

This creates `kali-bridge.exe` in the folder.

### 3. Get a Token from the TUI

In your remote TUI (browser), generate a bridge token (this feature is being added).

### 4. Run the Bridge

```powershell
# Normal run
.\kali-bridge.exe --connect YOUR_VPS_IP:8765 --token YOUR_TOKEN

# Recommended: Run as Administrator for better wireless access
# Right-click PowerShell → "Run as Administrator"
```

Leave the window open while using the TUI.

## What the Experience Will Feel Like (Once Fully Built)

- You chat normally with the AI in the browser.
- When the AI needs your laptop's hardware, it will automatically send commands to your bridge.
- You will see real output from your Windows machine appear in the chat.
- A status indicator will show "Bridge Connected" in the TUI.

## Future Improvements

- Better Windows wireless detection
- Support for Npcap (for advanced packet work)
- Linux version

## Need Help Right Now?

Run the bridge with `--help` to see current options.
