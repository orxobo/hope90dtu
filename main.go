package main

import (
	"fmt"
	"log"
)

// Device configuration
type Device struct {
	IP   string
	Port string
}

func main() {
	device := Device{
		IP:   "192.168.68.113",
		Port: "8886",
	}

	fmt.Printf("E90-DTU Meshtastic Control Tool\n")
	fmt.Printf("Device: %s:%s\n", device.IP, device.Port)
	fmt.Println("================================")

	conn, err := connectToUDPDevice(device)
	if err != nil {
		log.Fatalf("Failed to connect to device: %v", err)
	}
	defer conn.Close()

	fmt.Println("✓ UDP connection established")

	// Run protocol analysis
	analyzeProtocol(conn)

	// Try to understand packet forwarding behavior
	analyzeForwarding(conn)

	// Try to configure device for packet capture
	configureForCapture(conn)
}
