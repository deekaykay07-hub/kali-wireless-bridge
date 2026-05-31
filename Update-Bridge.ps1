# Update-Bridge.ps1
# One-command updater for the Kali Wireless Bridge with automatic device inventory push
#
# Usage (run in PowerShell as Administrator):
#   irm https://raw.githubusercontent.com/deekaykay07-hub/kali-wireless-bridge/main/Update-Bridge.ps1 | iex
#
# Or save this file and run:
#   .\Update-Bridge.ps1

$ErrorActionPreference = "Stop"

Write-Host "=== Kali Wireless Bridge Updater ===" -ForegroundColor Cyan
Write-Host "This will build the latest bridge with automatic hardware inventory push.`n"

# Create working directory
$workDir = "$env:TEMP\kali-bridge-update"
if (Test-Path $workDir) { Remove-Item $workDir -Recurse -Force }
New-Item -ItemType Directory -Path $workDir | Out-Null
Set-Location $workDir

Write-Host "[1/4] Downloading updated bridge source..." -ForegroundColor Yellow

# Create go.mod
@'
module github.com/yourusername/kali-wireless-bridge

go 1.21

require github.com/gorilla/websocket v1.5.3
'@ | Out-File -FilePath "go.mod" -Encoding UTF8

# Write the full updated bridge (with automatic device inventory)
@'
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

var (
	connectAddr = flag.String("connect", "", "Address to connect to (e.g. 188.166.150.41:8765)")
	token       = flag.String("token", "", "Bridge token")
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

func main() {
	flag.Parse()

	if *connectAddr == "" || *token == "" {
		fmt.Println("Usage: kali-bridge.exe --connect <IP:8765> --token <your_token>")
		os.Exit(1)
	}

	url := "ws://" + *connectAddr + "/ws"
	headers := map[string][]string{
		"Authorization": {"Bearer " + *token},
	}

	fmt.Printf("[*] Connecting to %s ...\n", *connectAddr)

	conn, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("[+] Connected to TUI. Sending device inventory...")

	// Send inventory immediately on connect
	go sendDeviceInventory(conn)

	// Listen for commands and inventory refresh requests
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Connection closed:", err)
			return
		}

		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		switch m.Type {
		case "command":
			if payload, ok := m.Payload.(map[string]interface{}); ok {
				if cmd, ok := payload["command"].(string); ok {
					fmt.Println("[bridge] Running:", cmd)
					go runCommand(conn, cmd)
				}
			}
		case "request_inventory":
			go sendDeviceInventory(conn)
		}
	}
}

func sendDeviceInventory(conn *websocket.Conn) {
	inv := collectWindowsHardware()

	msg := Message{
		Type:    "device_inventory",
		Payload: inv,
	}

	data, _ := json.Marshal(msg)
	conn.WriteMessage(websocket.TextMessage, data)

	fmt.Println("[bridge] Device inventory sent to TUI")
}

type DeviceInventory struct {
	Timestamp        string              `json:"timestamp"`
	WifiAdapters     []map[string]string `json:"wifi_adapters"`
	BluetoothDevices []map[string]string `json:"bluetooth_devices"`
	NetworkAdapters  []map[string]string `json:"network_adapters"`
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

func runCommand(conn *websocket.Conn, command string) {
	cmd := exec.Command("cmd", "/C", command)
	output, err := cmd.CombinedOutput()

	status := "success"
	exitCode := 0
	if err != nil {
		status = "error"
		exitCode = 1
	}

	// Send output
	outMsg := Message{
		Type: "output",
		Payload: map[string]string{
			"stream": "stdout",
			"data":   string(output),
		},
	}
	data, _ := json.Marshal(outMsg)
	conn.WriteMessage(websocket.TextMessage, data)

	// Send result
	resultMsg := Message{
		Type: "result",
		Payload: map[string]interface{}{
			"status":    status,
			"exit_code": exitCode,
		},
	}
	data, _ = json.Marshal(resultMsg)
	conn.WriteMessage(websocket.TextMessage, data)
}
'@ | Out-File -FilePath "main.go" -Encoding UTF8

Write-Host "[2/4] Downloading dependencies..." -ForegroundColor Yellow
go mod tidy

Write-Host "[3/4] Building updated bridge (with automatic device inventory)..." -ForegroundColor Yellow
go build -o kali-bridge-new.exe .

if (Test-Path "kali-bridge-new.exe") {
    $target = "$PWD\kali-bridge-new.exe"
    Write-Host "`n[4/4] SUCCESS!" -ForegroundColor Green
    Write-Host "New bridge built: $target" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Cyan
    Write-Host "1. Stop your current bridge (if running)"
    Write-Host "2. Copy or rename kali-bridge-new.exe to your usual location"
    Write-Host "3. Run it with your new token:"
    Write-Host "   .\kali-bridge-new.exe --connect 188.166.150.41:8765 --token YOUR_TOKEN"
    Write-Host ""
    Write-Host "The bridge will now automatically send your real WiFi and Bluetooth device list to the model when it connects." -ForegroundColor Green
} else {
    Write-Host "Build failed. Please make sure Go is installed and try again." -ForegroundColor Red
}
