// main.go
package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	E90_DTU_IP   = "192.168.0.125" // Default E90-DTU IP
	E90_DTU_PORT = 8886            // Default UDP port
)

func main() {
	// Create UDP connection
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", E90_DTU_IP, E90_DTU_PORT))
	if err != nil {
		log.Fatalf("Failed to resolve address: %v", err)
	}

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	log.Printf("Listening for E90-DTU packets on %s:%d", E90_DTU_IP, E90_DTU_PORT)

	// First, send an initialization packet to tell E90-DTU where to send data
	go initializeE90DTU()

	// Buffer for receiving packets
	buffer := make([]byte, 1024) // Max LoRa packet size is 240 bytes + headers

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		if n > 0 {
			packet := buffer[:n]
			go processRawPacket(packet, addr)
		}
	}
}

func initializeE90DTU() {
	// Wait a moment for the UDP server to start
	time.Sleep(2 * time.Second)

	// Send a packet to E90-DTU to set it as the destination
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", E90_DTU_IP, E90_DTU_PORT))
	if err != nil {
		log.Printf("Failed to send init packet: %v", err)
		return
	}
	defer conn.Close()

	// Send initialization packet (can be any data)
	_, err = conn.Write([]byte("INIT"))
	if err != nil {
		log.Printf("Failed to send init data: %v", err)
	}
	log.Println("Initialization packet sent to E90-DTU")
}

func processRawPacket(data []byte, source *net.UDPAddr) {
	log.Printf("[%s] Received %d bytes", source.IP, len(data))

	// Print hex dump
	fmt.Println(hex.Dump(data))

	// Parse packet based on length
	if len(data) >= 1 {
		// Check if this looks like a Meshtastic packet
		// Meshtastic often starts with 0x94 (Protocol marker)
		if data[0] == 0x94 {
			log.Println("  → Possible Meshtastic packet detected (0x94 header)")
			parseMeshtasticPacket(data)
		}

		// If Data RSSI is enabled, last byte is RSSI
		rssi := int8(data[len(data)-1])
		log.Printf("  → RSSI: %d dBm", rssi)

		// Extract payload (everything except RSSI byte)
		payload := data[:len(data)-1]
		log.Printf("  → Payload: %s", hex.EncodeToString(payload))
	}

	// You could also query channel RSSI using C0 C1 C2 C3 command
	if shouldQueryChannelRSSI() {
		queryChannelRSSI()
	}
}
