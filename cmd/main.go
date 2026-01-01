// cmd/main.go
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Command line flags
	configFile := flag.String("config", "", "Configuration file (JSON)")
	dtuIP := flag.String("ip", "192.168.0.125", "E90-DTU IP address")
	port := flag.Int("port", 8886, "UDP port")
	channel := flag.Int("channel", 70, "LoRa channel (0-80)")
	captureOnly := flag.Bool("capture-only", false, "Only capture, don't translate")
	outputFormat := flag.String("format", "json", "Output format: json, csv, mqtt")

	flag.Parse()

	// Update global constants
	E90_DTU_IP = *dtuIP
	E90_DTU_PORT = uint16(*port)

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Load or create configuration
	var config E90Config
	if *configFile != "" {
		data, err := os.ReadFile(*configFile)
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}
		if err := json.Unmarshal(data, &config); err != nil {
			log.Fatalf("Failed to parse config: %v", err)
		}
	} else {
		config = DefaultMeshtasticConfig
		config.Channel = uint8(*channel)
	}

	// Apply configuration to E90-DTU
	log.Println("Configuring E90-DTU...")
	if err := ConfigureE90ViaWeb(config); err != nil {
		log.Printf("Warning: Could not configure via web: %v", err)
		log.Println("Using existing configuration")
	}

	// Start packet capture
	log.Println("Starting packet capture...")
	go startUDPServer()

	// Start MQTT bridge if enabled
	if !*captureOnly && *outputFormat == "mqtt" {
		go startMQTTBridge()
	}

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down...")
}

func startUDPServer() {
	// Implementation from main.go
}

func startMQTTBridge() {
	// Bridge to MQTT (for Home Assistant, Node-RED, etc.)
	log.Println("Starting MQTT bridge...")
	// Implementation would connect to MQTT broker
	// and publish packets as JSON messages
}
