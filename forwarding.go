package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"
)

func analyzeForwarding(conn *net.UDPConn) {
	fmt.Println("\n=== FORWARDING ANALYSIS ===")

	// Try various parameter combinations for command 0x01 (forward)
	parameterTests := [][]byte{
		{0xc0, 0xc1, 0xc2, 0xc3, 0x01, 0x00}, // Forward with 0x00 parameter
		{0xc0, 0xc1, 0xc2, 0xc3, 0x01, 0x01}, // Forward with 0x01 parameter (already known to work)
		{0xc0, 0xc1, 0xc2, 0xc3, 0x01, 0x02}, // Forward with 0x02 parameter
		{0xc0, 0xc1, 0xc2, 0xc3, 0x01, 0x03}, // Forward with 0x03 parameter
		{0xc0, 0xc1, 0xc2, 0xc3, 0x01, 0x04}, // Forward with 0x04 parameter
		{0xc0, 0xc1, 0xc2, 0xc3, 0x01, 0x05}, // Forward with 0x05 parameter
	}

	for i, cmd := range parameterTests {
		fmt.Printf("\nSending parameter test command %d: %s\n", i+1, hex.EncodeToString(cmd))
		err := sendUDPCommand(conn, cmd)
		if err != nil {
			log.Printf("Failed to send parameter test command %d: %v", i+1, err)
			continue
		}

		response, err := receiveUDPResponse(conn, 3*time.Second)
		if err != nil {
			log.Printf("No response to parameter test command %d", i+1)
			continue
		}

		fmt.Printf("Response to parameter test command %d: %s\n", i+1, hex.EncodeToString(response))
		fmt.Printf("Response bytes: %v\n", response)
	}
}

func configureForCapture(conn *net.UDPConn) {
	fmt.Println("\n=== PACKET CAPTURE CONFIGURATION ===")

	// Try different command/parameter combinations to see if we can configure
	// capture behavior or UDP forwarding
	fmt.Println("Trying to understand capture configuration...")

	// Commands 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08
	// These might be for configuration or status
	for cmd := 0x02; cmd <= 0x08; cmd++ {
		fmt.Printf("\nTrying command 0x%02x with parameter 0x00\n", cmd)
		command := []byte{0xc0, 0xc1, 0xc2, 0xc3, byte(cmd), 0x00}

		err := sendUDPCommand(conn, command)
		if err != nil {
			log.Printf("Failed to send command 0x%02x: %v", cmd, err)
			continue
		}

		response, err := receiveUDPResponse(conn, 3*time.Second)
		if err != nil {
			log.Printf("No response to command 0x%02x", cmd)
			continue
		}

		fmt.Printf("Response to command 0x%02x: %s\n", cmd, hex.EncodeToString(response))
		fmt.Printf("Response bytes: %v\n", response)
	}
}
