package device

import (
	"fmt"
)

func (e *E90Device) GetBackgroundNoise() string {
	err := e.SendUDPCommand([]byte{0xc0, 0xc1, 0xc2, 0xc3, 0x00, 0x01})
	if err != nil {
		return ""
	}

	response, err := e.ReceiveUDPResponse()
	if err != nil {
		return ""
	}

	return interpretRSSIResponse(response)
}

func (e *E90Device) GetLastResponseNoise() string {
	err := e.SendUDPCommand([]byte{0xc0, 0xc1, 0xc2, 0xc3, 0x01, 0x01})
	if err != nil {
		return ""
	}

	response, err := e.ReceiveUDPResponse()
	if err != nil {
		return ""
	}

	return interpretRSSIResponse(response)
}

// return is 4 bytes: C1 + [Address] + [Read Length] + [RSSI Values]
func interpretRSSIResponse(response []byte) string {
	if len(response) < 4 {
		return "Invalid response"
	}

	rssiValue := int(response[3])
	return fmt.Sprintf("RSSI: %d dBm", rssiValue/2)
}
