package main

import (
	"fmt"
)

func interpretRSSIResponse(response []byte) string {
	if len(response) < 4 {
		return "Invalid response"
	}

	// Response format: [response_byte] [command] [parameter] [data]
	rssiValue := int(response[3])
	return fmt.Sprintf("RSSI: %d dBm", rssiValue)
}

func interpretStatusResponse(response []byte) string {
	if len(response) < 5 {
		return "Invalid response"
	}

	// Response format: [response_byte] [command] [parameter] [data] [status]
	statusValue := int(response[4])
	return fmt.Sprintf("Status: %d", statusValue)
}

func interpretForwardResponse(response []byte) string {
	if len(response) < 4 {
		return "Invalid response"
	}

	// Response format: [response_byte] [command] [parameter] [data]
	// The data field is likely indicating success/failure or mode
	dataValue := int(response[3])
	return fmt.Sprintf("Forward mode enabled (data: %d)", dataValue)
}
