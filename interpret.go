package main

import (
	"fmt"
)

// return is 4 bytes: C1 + [Address] + [Read Length] + [RSSI Values]
func interpretRSSIResponse(response []byte) string {
	if len(response) < 4 {
		return "Invalid response"
	}

	rssiValue := int(response[3])
	return fmt.Sprintf("RSSI: %d dBm", rssiValue/2)
}

func genericStringResponse(response []byte) string {
	return fmt.Sprintf("Returned value: %s", string(response))
}
