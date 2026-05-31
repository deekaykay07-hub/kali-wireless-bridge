# Hardware Inventory Push Update (for Bridge)

This update makes the bridge automatically send the list of real WiFi and Bluetooth devices to the TUI when it connects. 
The TUI then injects this list into the AI's context so the model knows the *exact* device names instead of guessing.

## Changes Needed on Windows Bridge Side

In your bridge code, after successful connection and auth (right after you set up the websocket), add this call:

```go
go sendDeviceInventory(wsConn)   // run in goroutine so it doesn't block
```

Add the following function (you can put it in internal/bridge or in main.go):

```go
package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

type DeviceInventory struct {
	Timestamp        string              `json:"timestamp"`
	WifiAdapters     []map[string]string `json:"wifi_adapters"`
	BluetoothDevices []map[string]string `json:"bluetooth_devices"`
	NetworkAdapters  []map[string]string `json:"network_adapters"`
}

func sendDeviceInventory(conn *websocket.Conn) {
	inv := collectWindowsHardware()

	msg := map[string]interface{}{
		"type":    "device_inventory",
		"payload": inv,
	}

	data, _ := json.Marshal(msg)
	conn.WriteMessage(websocket.TextMessage, data)

	fmt.Println("[bridge] Sent device inventory to TUI")
}

func collectWindowsHardware() DeviceInventory {
	inv := DeviceInventory{
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// WiFi adapters
	wifiCmd := `Get-NetAdapter | Where-Object {$_.Name -like '*Wi-Fi*' -or $_.InterfaceDescription -like '*WiFi*' -or $_.Name -like '*Wireless*'} | Select-Object Name, InterfaceDescription, Status, MacAddress | ConvertTo-Json -Compress`
	wifiOut, _ := exec.Command("powershell", "-Command", wifiCmd).Output()
	var wifi []map[string]string
	json.Unmarshal(wifiOut, &wifi)
	inv.WifiAdapters = wifi

	// Bluetooth devices
	btCmd := `Get-PnpDevice -Class Bluetooth | Select-Object Name, Status | ConvertTo-Json -Compress`
	btOut, _ := exec.Command("powershell", "-Command", btCmd).Output()
	var bt []map[string]string
	json.Unmarshal(btOut, &bt)
	inv.BluetoothDevices = bt

	// All network adapters (limited)
	netCmd := `Get-NetAdapter | Select-Object Name, InterfaceDescription, Status | ConvertTo-Json -Compress`
	netOut, _ := exec.Command("powershell", "-Command", netCmd).Output()
	var nets []map[string]string
	json.Unmarshal(netOut, &nets)
	if len(nets) > 0 {
		if len(nets) > 8 {
			nets = nets[:8]
		}
		inv.NetworkAdapters = nets
	}

	return inv
}
```

## Optional: Support Refresh

In your message handler (where you receive commands from TUI), add:

```go
if msg.Type == "request_inventory" {
    go sendDeviceInventory(conn)
}
```

This way the AI (or you) can ask for a fresh device list later.

## TUI Side Commands (already added)

- `/refresh-devices` or `/devices` → asks the bridge for a fresh inventory
- When inventory arrives, the model automatically gets the exact device names in its context.

## After Updating

1. Add the code above.
2. Rebuild: `go build -o kali-bridge.exe .\cmd\bridge`
3. Restart the bridge with the new token.

The TUI will now automatically receive and tell the model the real device names.