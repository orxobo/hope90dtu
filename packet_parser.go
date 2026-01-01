// packet_parser.go
package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

// MeshtasticPacket represents a decoded Meshtastic packet
type MeshtasticPacket struct {
	Timestamp   time.Time
	FromNodeID  uint32
	ToNodeID    uint32
	Channel     uint32
	DataRate    string
	Payload     []byte
	RSSI        int8
	SNR         float32
	Frequency   float64 // in MHz
	IsEncrypted bool
	MessageType string
}

// LoRaPacket represents raw LoRa PHY layer packet
type LoRaPacket struct {
	Timestamp time.Time
	Frequency float64 // MHz
	SF        int     // Spreading Factor
	BW        int     // Bandwidth in kHz
	CR        string  // Coding Rate
	Payload   []byte
	RSSI      int8
	SNR       float32
	CRC       bool
}

func parseMeshtasticPacket(data []byte) *MeshtasticPacket {
	if len(data) < 4 {
		return nil
	}

	packet := &MeshtasticPacket{
		Timestamp: time.Now(),
		RSSI:      int8(data[len(data)-1]), // Last byte is RSSI if enabled
		Payload:   data[:len(data)-1],
	}

	// Basic Meshtastic protocol parsing
	// Protocol marker is usually 0x94
	if data[0] == 0x94 {
		packet.MessageType = "MeshPacket"

		// Try to extract node IDs (positions may vary)
		if len(data) >= 8 {
			// This is simplified - real parsing needs reverse engineering
			packet.FromNodeID = binary.BigEndian.Uint32(data[1:5])
			packet.ToNodeID = binary.BigEndian.Uint32(data[5:9])
		}
	}

	log.Printf("Parsed Meshtastic packet: From=0x%X, To=0x%X, Type=%s",
		packet.FromNodeID, packet.ToNodeID, packet.MessageType)

	return packet
}

func queryChannelRSSI() {
	// Send C0 C1 C2 C3 command to query channel RSSI
	// Command format: C0 C1 C2 C3 + starting address + read length
	cmd := []byte{0xC0, 0xC1, 0xC2, 0xC3, 0x00, 0x02}

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", E90_DTU_IP, E90_DTU_PORT))
	if err != nil {
		log.Printf("Failed to query channel RSSI: %v", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(cmd)
	if err != nil {
		log.Printf("Failed to send RSSI query: %v", err)
	}

	log.Println("Channel RSSI query sent")
}

func shouldQueryChannelRSSI() bool {
	// Query every 30 seconds
	return time.Now().Unix()%30 == 0
}
