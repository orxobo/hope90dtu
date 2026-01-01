// config_manager.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// E90Config represents E90-DTU configuration
type E90Config struct {
	// Wireless Parameters
	ModuleAddress uint16 `json:"module_address"`
	Channel       uint8  `json:"channel"`
	NetID         uint8  `json:"net_id"`
	AirBaud       string `json:"air_baud"` // "2.4K", "62.5K", etc
	TxPower       string `json:"tx_power"` // "maximum", "middle", "low", "extremelow"
	TransMode     string `json:"trans_mode"`
	PacketLength  uint16 `json:"packet_length"`
	WORRole       string `json:"wor_role"`
	WRCycle       uint16 `json:"wor_cycle"`
	Key           uint16 `json:"key"`
	LBTEnable     bool   `json:"lbt_enable"`
	RelayEnable   bool   `json:"relay_enable"`
	DataRSSI      bool   `json:"data_rssi"`
	ChannelRSSI   bool   `json:"channel_rssi"`

	// Network Parameters
	WorkMode   string `json:"work_mode"` // "UDPS", "TCPC", etc
	LocalPort  uint16 `json:"local_port"`
	RemoteIP   string `json:"remote_ip"`
	RemotePort uint16 `json:"remote_port"`
}

// Default configuration for Meshtastic capture
var DefaultMeshtasticConfig = E90Config{
	ModuleAddress: 0,
	Channel:       70, // For 920.125 MHz
	NetID:         0,
	AirBaud:       "62.5Kbps",
	TxPower:       "maximum",
	TransMode:     "Normal",
	PacketLength:  240,
	WORRole:       "close",
	WRCycle:       2000,
	Key:           0,
	LBTEnable:     false,
	RelayEnable:   false,
	DataRSSI:      true,
	ChannelRSSI:   true,
	WorkMode:      "UDPS",
	LocalPort:     8886,
	RemoteIP:      "0.0.0.0",
	RemotePort:    0,
}

// ConfigureE90ViaWeb sends configuration via HTTP POST
func ConfigureE90ViaWeb(config E90Config) error {
	// E90-DTU web interface usually accepts JSON configuration
	configURL := fmt.Sprintf("http://%s/config", E90_DTU_IP)

	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	req, err := http.NewRequest("POST", configURL, strings.NewReader(string(configJSON)))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send config: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("configuration failed with status: %s", resp.Status)
	}

	log.Println("E90-DTU configuration applied successfully")
	return nil
}

// SendATCommand sends AT command to E90-DTU
func SendATCommand(cmd string) (string, error) {
	// Network AT command format
	atCmd := fmt.Sprintf("AT+%s\r\n", cmd)

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", E90_DTU_IP, E90_DTU_PORT))
	if err != nil {
		return "", fmt.Errorf("failed to connect for AT command: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(atCmd))
	if err != nil {
		return "", fmt.Errorf("failed to send AT command: %v", err)
	}

	// Read response
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(buffer[:n]), nil
}
