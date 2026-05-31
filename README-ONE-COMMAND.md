# One-Command Bridge Update (Device Inventory)

This gives you the updated bridge with automatic hardware inventory push in **minimal steps**.

## Single PowerShell Command (Run as Administrator)

Open **PowerShell as Administrator** and run this one command:

```powershell
irm https://raw.githubusercontent.com/deekaykay07-hub/kali-wireless-bridge/main/Update-Bridge.ps1 | iex
```

## What This Gives You

- Bridge automatically sends your real WiFi adapters + Bluetooth devices to the AI when it connects.
- The model will now see the exact device names (no more guessing "Bluetooth Device 1").
- You can still use `/refresh-devices` in the TUI to force a fresh scan anytime.

## After Running

Just use the new `kali-bridge-new.exe` with your latest token:

```powershell
.\kali-bridge-new.exe --connect 188.166.150.41:8765 --token <your-new-token>
```

That's it. Minimal steps, one build, full hardware awareness for the model.