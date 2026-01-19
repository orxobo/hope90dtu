package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

type E90DTU struct {
	conn net.Conn
	addr string
}

func NewE90DTU(ip string, port int) (*E90DTU, error) {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	return &E90DTU{
		conn: conn,
		addr: addr,
	}, nil
}

func (e *E90DTU) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

func (e *E90DTU) SendCommand(data []byte, expectResponse bool) ([]byte, error) {
	fmt.Printf("TX: %s\n", hex.EncodeToString(data))

	// Send data
	_, err := e.conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("write error: %v", err)
	}

	if !expectResponse {
		return nil, nil
	}

	// Wait for response
	e.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	n, err := e.conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("read error: %v", err)
	}

	response := buffer[:n]
	fmt.Printf("RX: %s\n", hex.EncodeToString(response))
	return response, nil
}

// GetRSSI reads RSSI values from the E90
func (e *E90DTU) GetRSSI(addr, length byte) (int, error) {
	cmd := []byte{0xC0, 0xC1, 0xC2, 0xC3, addr, length}
	resp, err := e.SendCommand(cmd, true)
	if err != nil {
		return 0, err
	}

	// Parse response: C1 addr length value...
	if len(resp) < 4 {
		return 0, fmt.Errorf("response too short: %d bytes", len(resp))
	}

	if resp[0] != 0xC1 {
		return 0, fmt.Errorf("invalid response header: %02x", resp[0])
	}

	if resp[1] != addr || resp[2] != length {
		return 0, fmt.Errorf("address/length mismatch")
	}

	// Extract value
	var value int
	if length == 1 {
		value = int(resp[3])
	} else if length == 2 {
		if len(resp) < 5 {
			return 0, fmt.Errorf("response too short for 2-byte value")
		}
		value = int(resp[3]) | (int(resp[4]) << 8) // Little-endian
	} else {
		return 0, fmt.Errorf("unsupported length: %d", length)
	}

	return value, nil
}

// GetNoiseLevel returns current channel noise in dBm
func (e *E90DTU) GetNoiseLevel() (float64, error) {
	value, err := e.GetRSSI(0x00, 1)
	if err != nil {
		return 0, err
	}
	return -float64(value) / 2.0, nil
}

// GetLastRSSI returns RSSI of last received packet in dBm
func (e *E90DTU) GetLastRSSI() (float64, error) {
	value, err := e.GetRSSI(0x01, 1)
	if err != nil {
		return 0, err
	}
	return -float64(value) / 2.0, nil
}

// SendLoRaData sends data via LoRa
func (e *E90DTU) SendLoRaData(data []byte) error {
	_, err := e.SendCommand(data, false)
	return err
}

// SendFixedPoint sends data to specific address/channel
func (e *E90DTU) SendFixedPoint(addr uint16, channel byte, data []byte) error {
	packet := make([]byte, 3+len(data))
	packet[0] = byte(addr >> 8)
	packet[1] = byte(addr & 0xFF)
	packet[2] = channel
	copy(packet[3:], data)

	return e.SendLoRaData(packet)
}

// SendBroadcast sends to all nodes on a channel
func (e *E90DTU) SendBroadcast(channel byte, data []byte) error {
	return e.SendFixedPoint(0xFFFF, channel, data)
}

// ListenForResponses listens for incoming data
func (e *E90DTU) ListenForResponses(timeout time.Duration) ([]byte, error) {
	e.conn.SetReadDeadline(time.Now().Add(timeout))
	buffer := make([]byte, 1024)

	n, err := e.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

func main() {
	fmt.Println("=== E90-DTU Advanced Test ===")

	// Connect to E90
	dtu, err := NewE90DTU("192.168.68.113", 8886)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer dtu.Close()

	fmt.Println("✓ Connected to E90-DTU")

	// Test 1: Get noise level
	fmt.Println("\n1. Measuring channel noise...")
	noise, err := dtu.GetNoiseLevel()
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Current noise: %.1f dBm\n", noise)
	}

	// Test 2: Get last RSSI
	fmt.Println("\n2. Checking last received signal...")
	lastRSSI, err := dtu.GetLastRSSI()
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Last packet RSSI: %.1f dBm\n", lastRSSI)
	}

	// Test 3: Test different data formats
	fmt.Println("\n3. Testing various data formats...")

	// Simple text
	fmt.Println("   a) Simple text message...")
	err = dtu.SendLoRaData([]byte("Simple test from Go"))
	if err != nil {
		fmt.Printf("      Error: %v\n", err)
	} else {
		fmt.Println("      Sent")
	}
	time.Sleep(500 * time.Millisecond)

	// Binary data
	fmt.Println("   b) Binary data...")
	binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	err = dtu.SendLoRaData(binaryData)
	if err != nil {
		fmt.Printf("      Error: %v\n", err)
	} else {
		fmt.Printf("      Sent: %s\n", hex.EncodeToString(binaryData))
	}
	time.Sleep(500 * time.Millisecond)

	// Broadcast
	fmt.Println("   c) Broadcast on channel 69...")
	err = dtu.SendBroadcast(0x45, []byte("Broadcast test"))
	if err != nil {
		fmt.Printf("      Error: %v\n", err)
	} else {
		fmt.Println("      Sent")
	}

	// Test 4: Try to receive data
	fmt.Println("\n4. Listening for responses (5 seconds)...")

	start := time.Now()
	for time.Since(start) < 5*time.Second {
		dtu.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		data, err := dtu.ListenForResponses(1 * time.Second)

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue // No data yet
			}
			fmt.Printf("   Read error: %v\n", err)
			break
		}

		fmt.Printf("   Received %d bytes: %s\n", len(data), hex.EncodeToString(data))

		// Check if it's an RSSI response
		if len(data) >= 4 && data[0] == 0xC1 {
			addr := data[1]
			length := data[2]
			fmt.Printf("   RSSI response - Addr: 0x%02x, Len: %d\n", addr, length)

			if length == 1 && len(data) >= 4 {
				value := int(data[3])
				dbm := -float64(value) / 2.0
				fmt.Printf("   Value: 0x%02x (%d) = %.1f dBm\n", value, value, dbm)
			}
		} else {
			// Might be data from LoRa
			fmt.Printf("   ASCII: %s\n", string(data))
		}
	}

	fmt.Println("\n5. Testing more RSSI queries...")
	// Test reading 2 bytes at once
	fmt.Println("   Reading 2 bytes from address 0x00...")
	cmd := []byte{0xC0, 0xC1, 0xC2, 0xC3, 0x00, 0x02}
	resp, err := dtu.SendCommand(cmd, true)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Response: %s\n", hex.EncodeToString(resp))
		if len(resp) >= 5 && resp[0] == 0xC1 {
			value := int(resp[3]) | (int(resp[4]) << 8)
			fmt.Printf("   Combined value: 0x%04x (%d)\n", value, value)
		}
	}

	fmt.Println("\n=== Test Complete ===")
	fmt.Println("Summary:")
	fmt.Println("- RSSI queries work")
	fmt.Println("- Data sending works")
	fmt.Println("- No responses yet (might need T-Deck to transmit)")
	fmt.Println("\nNext step: Test with T-Deck transmitting")
}
