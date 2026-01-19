package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

// E90DTUCommunicator handles communication with E90-DTU
type E90DTUCommunicator struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

// NewE90DTUCommunicator creates a new communicator
func NewE90DTUCommunicator(ip string, port int) (*E90DTUCommunicator, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &E90DTUCommunicator{
		conn: conn,
		addr: addr,
	}, nil
}

// Close closes the connection
func (e *E90DTUCommunicator) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

// SendCommand sends a command and returns the response
func (e *E90DTUCommunicator) SendCommand(command []byte) ([]byte, error) {
	_, err := e.conn.Write(command)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	e.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := e.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

// GetRSSI gets the current RSSI value
func (e *E90DTUCommunicator) GetRSSI() (int, error) {
	cmd := []byte{0xc0, 0xc1, 0xc2, 0xc3, 0x00, 0x01}
	response, err := e.SendCommand(cmd)
	if err != nil {
		return 0, err
	}

	if len(response) >= 4 && response[0] == 0xc1 && response[1] == 0x00 && response[2] == 0x01 {
		return int(response[3]), nil
	}

	return 0, fmt.Errorf("unexpected RSSI response: %s", hex.EncodeToString(response))
}

// GetDeviceInfo gets device information
func (e *E90DTUCommunicator) GetDeviceInfo() error {
	cmd := []byte{0xc0, 0xc1, 0xc2, 0xc3, 0x00, 0x02}
	response, err := e.SendCommand(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("Device info response: %s\n", hex.EncodeToString(response))
	return nil
}

// SendMeshPacket sends a mesh packet to the E90-DTU
func (e *E90DTUCommunicator) SendMeshPacket(packet []byte) error {
	// Send command to indicate we're sending a packet
	cmd := []byte{0xc0, 0xc1, 0xc2, 0xc3, 0x01, 0x00}
	_, err := e.SendCommand(cmd)
	if err != nil {
		return err
	}

	// Send the actual packet data
	_, err = e.conn.Write(packet)
	return err
}

// CreateNodePacket creates a basic node packet
func (e *E90DTUCommunicator) CreateNodePacket() []byte {
	// Create a minimal node packet
	// This is a simplified version - real mesh protocol would be more complex

	// Node packet structure (example):
	// [Header][Node ID][Packet Type][Data]

	packet := make([]byte, 20)

	// Header
	packet[0] = 0x01 // Packet type
	packet[1] = 0x02 // Version
	packet[2] = 0x03 // Channel
	packet[3] = 0x04 // Flags

	// Node ID (4 bytes)
	packet[4] = 0x01
	packet[5] = 0x02
	packet[6] = 0x03
	packet[7] = 0x04

	// Payload (simplified)
	packet[8] = 0x00  // Payload type
	packet[9] = 0x01  // Payload length
	packet[10] = 0x48 // 'H'
	packet[11] = 0x65 // 'e'
	packet[12] = 0x6c // 'l'
	packet[13] = 0x6c // 'l'
	packet[14] = 0x6f // 'o'

	// CRC or other fields
	packet[15] = 0x00
	packet[16] = 0x00
	packet[17] = 0x00
	packet[18] = 0x00
	packet[19] = 0x00

	return packet
}

// Main function
func main() {
	fmt.Println("E90-DTU Mesh Node Creation Tool")
	fmt.Println("===============================")

	ip := "192.168.68.113"
	port := 8886

	// Connect to E90-DTU
	comm, err := NewE90DTUCommunicator(ip, port)
	if err != nil {
		fmt.Printf("Failed to connect to E90-DTU: %v\n", err)
		return
	}
	defer comm.Close()

	fmt.Printf("Connected to E90-DTU at %s:%d\n", ip, port)

	// Test basic functionality
	fmt.Println("\n--- Testing Basic Functionality ---")

	// Get RSSI
	rssi, err := comm.GetRSSI()
	if err != nil {
		fmt.Printf("Failed to get RSSI: %v\n", err)
	} else {
		fmt.Printf("RSSI: %d\n", rssi)
	}

	// Get device info
	fmt.Println("\nGetting device info...")
	if err := comm.GetDeviceInfo(); err != nil {
		fmt.Printf("Failed to get device info: %v\n", err)
	}

	// Create and send a node packet
	fmt.Println("\n--- Creating Node Packet ---")

	packet := comm.CreateNodePacket()
	fmt.Printf("Created node packet: %s\n", hex.EncodeToString(packet))

	// Try to send the packet
	fmt.Println("Sending node packet...")
	if err := comm.SendMeshPacket(packet); err != nil {
		fmt.Printf("Failed to send packet: %v\n", err)
	} else {
		fmt.Println("Packet sent successfully")
	}

	fmt.Println("\n--- Next Steps ---")
	fmt.Println("1. The E90-DTU likely needs properly formatted mesh packets")
	fmt.Println("2. You may need to study the actual Meshtastic protocol")
	fmt.Println("3. Consider using existing tools like meshtastic-cli to generate valid packets")
	fmt.Println("4. The E90-DTU may have its own specific packet format")
	fmt.Println("5. For T-Deck to recognize a node, it needs a valid node announcement packet")
}
