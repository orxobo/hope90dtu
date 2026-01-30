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

	// Example 1: Get simple string
	if model, err := client.GetModel(); err != nil {
		fmt.Printf("Model Failed: %v\n", err)
	} else {
		fmt.Printf("Model: %s\n", model)
	}

	// Example 2: Get struct
	if lora, err := client.GetLora(); err != nil {
		fmt.Printf("Lora Failed: %v\n", err)
	} else {
		fmt.Printf("Lora Baud Rate: %d\n", lora.Baud)
	}

	// Example 3: Set value
	if port, err := client.SetLocalPort(8088); err != nil {
		fmt.Printf("Set Port Failed: %v\n", err)
	} else {
		fmt.Printf("New Port: %d\n", port)
	}
}
