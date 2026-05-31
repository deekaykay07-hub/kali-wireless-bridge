# User Guide: Kali Wireless Bridge

This guide explains exactly what the final experience will feel like and what you need to do on your laptop.

## What the Finished System Looks Like

1. You open your browser and go to your TUI:
   `http://your-vps-ip:7681`

2. You see a normal chat interface with a very smart wireless expert AI.

3. At the top or side, you see a clear indicator:
   **✅ Local Bridge Connected** (YourLaptopName)
   or
   **⚠️ No local bridge connected** (Hardware commands will be limited)

4. When you ask the AI to do something that needs your hardware (monitor mode, Bluetooth scan, etc.), it will:
   - Make a clear plan
   - Send the actual commands to your laptop
   - Show you the real output from your machine in the chat

Example conversation:
> You: Put my adapter in monitor mode and scan for networks
> AI: [Plan shown in the Planning Panel]
> AI: Running on your local machine...
> [Real output from airodump-ng appears live in the chat]

## What You Need to Do on Your Laptop (Step by Step)

### 1. Build the Bridge (One-time)

On your laptop (preferably Linux for best wireless support):

```bash
git clone https://github.com/deekaykay07-hub/kali-wireless-bridge.git
cd kali-wireless-bridge
make build
```

This creates the `kali-bridge` binary.

### 2. Get a Token from the TUI

- In the TUI, there will be a command or button like `/enable-bridge` or a settings panel.
- It will generate a one-time token for your bridge.

### 3. Start the Bridge

```bash
sudo ./kali-bridge --connect your-vps-ip:8765 --token YOUR_TOKEN_HERE
```

- Run it with `sudo` when you need monitor mode / injection.
- Leave it running in a terminal while you use the TUI.

### 4. Use the TUI Normally

The AI will automatically use your local hardware when needed. You don't have to do anything special in the chat.

## Important Notes

- The bridge only runs when **you** start it on your laptop.
- It makes an outbound connection only (safe for most firewalls).
- You can stop it anytime by pressing Ctrl+C.
- For real monitor mode work, Linux is strongly recommended on the laptop side.
