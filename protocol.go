package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"
)

func analyzeProtocol(conn *net.UDPConn) {
	fmt.Println("\n=== PROTOCOL ANALYSIS ===")

	// Test known working commands
	fmt.Println("\nTest 1: RSSI Query")
	rssiCommand := []byte{0xc0, 0xc1, 0xc2, 0xc3, 0x00, 0x01}
	fmt.Printf("Sending RSSI command: %s\n", hex.EncodeToString(rssiCommand))

	err := sendUDPCommand(conn, rssiCommand)
	if err != nil {
		log.Printf("Failed to send RSSI command: %v", err)
	} else {
		fmt.Println("RSSI command sent successfully")
	}

	response, err := receiveUDPResponse(conn, 5*time.Second)
	if err != nil {
		log.Printf("Failed to receive RSSI response: %v", err)
	} else {
		fmt.Printf("Received RSSI response: %s\n", hex.EncodeToString(response))
		fmt.Printf("Response bytes: %v\n", response)
		fmt.Printf("Interpreted: %s\n", interpretRSSIResponse(response))
	}

	// Test Status command
	fmt.Println("\nTest 2: Status Query")
	statusCommand := []byte{0xc0, 0xc1, 0xc2, 0xc3, 0x00, 0x02}
	fmt.Printf("Sending Status command: %s\n", hex.EncodeToString(statusCommand))

	err = sendUDPCommand(conn, statusCommand)
	if err != nil {
		log.Printf("Failed to send Status command: %v", err)
	} else {
		fmt.Println("Status command sent successfully")
	}

	response, err = receiveUDPResponse(conn, 5*time.Second)
	if err != nil {
		log.Printf("Failed to receive Status response: %v", err)
	} else {
		fmt.Printf("Received Status response: %s\n", hex.EncodeToString(response))
		fmt.Printf("Response bytes: %v\n", response)
		fmt.Printf("Interpreted: %s\n", interpretStatusResponse(response))
	}

	// Test Forward command
	fmt.Println("\nTest 3: Forward/Enable Command")
	forwardCommand := []byte{0xc0, 0xc1, 0xc2, 0xc3, 0x01, 0x01}
	fmt.Printf("Sending Forward command: %s\n", hex.EncodeToString(forwardCommand))

	err = sendUDPCommand(conn, forwardCommand)
	if err != nil {
		log.Printf("Failed to send Forward command: %v", err)
	} else {
		fmt.Println("Forward command sent successfully")
	}

	response, err = receiveUDPResponse(conn, 5*time.Second)
	if err != nil {
		log.Printf("Failed to receive Forward response: %v", err)
	} else {
		fmt.Printf("Received Forward response: %s\n", hex.EncodeToString(response))
		fmt.Printf("Response bytes: %v\n", response)
		fmt.Printf("Interpreted: %s\n", interpretForwardResponse(response))
	}
}
