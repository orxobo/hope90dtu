package main

import (
	"fmt"
	"log"
	"net/netip"
	"strings"
)

// Device configuration
const (
	IP   = "192.168.68.113"
	PORT = 8886
)

func main() {

	deviceAddressPort := netip.AddrPortFrom(
		netip.MustParseAddr(IP),
		PORT,
	)

	fmt.Println("E90-DTU Meshtastic Control Tool")
	fmt.Println("Device: ", deviceAddressPort.String())
	fmt.Println(strings.Repeat("=", 30))

	e90, err := NewE90Device(deviceAddressPort)
	if err != nil {
		log.Fatalf("Failed to connect to device: %v", err)
	}
	defer e90.conn.Close()

	fmt.Println("✓ UDP connection established")

	//analyzeProtocol(e90)

	client := NewATClient(e90)

	// 1. Direct Method Usage
	model, err := client.GetModel()
	if err != nil {
		fmt.Println("error getting model: ", err.Error())
	}
	fmt.Println("Model: ", model)

	// 2. Enum Usage
	port, err := client.Run(CmdLPort)
	if err != nil {
		fmt.Println("error getting port: ", err.Error())
	}
	fmt.Println("Port: ", port)

	// 3. Raw String Usage (Relies on stringer generated strings)
	//client.RunRaw("LPORT", "9090")

	client.ListCommands()
}
