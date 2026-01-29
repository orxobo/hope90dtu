package main

import (
	"fmt"
	"log"
	"net/netip"
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
	fmt.Println("================================")

	e90, err := NewE90Device(deviceAddressPort)
	if err != nil {
		log.Fatalf("Failed to connect to device: %v", err)
	}
	defer e90.conn.Close()

	fmt.Println("✓ UDP connection established")

	// Run protocol analysis
	//analyzeProtocol(e90)

	e90.ATInitialise()
}
