package device

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

type E90Command struct {
	CommandBytes      []byte
	CommandName       string
	InterpretFunction func(response []byte) string
	SecondsTimeout    int
	Description       string
}

func analyzeProtocol(e90 *E90Device) {
	fmt.Println("\n=== PROTOCOL ANALYSIS ===")

	// can read both registers at once with []byte{0xc0, 0xc1, 0xc2, 0xc3, 0x00, 0x02}
	// ie. C0 C1 C2 C3 + [Start Address] + [Read Length], returns 5 bytes instead of 4
	commands := []E90Command{
		{
			CommandBytes:      []byte{0xc0, 0xc1, 0xc2, 0xc3, 0x00, 0x01},
			CommandName:       "RSSI background noise",
			InterpretFunction: interpretRSSIResponse,
			Description:       "Read first RSSI register.\nMeasures the background noise level on the current channel (no signal present).",
		},
		{
			CommandBytes:      []byte{0xc0, 0xc1, 0xc2, 0xc3, 0x01, 0x01},
			CommandName:       "RSSI Last Response",
			InterpretFunction: interpretRSSIResponse,
			Description:       "Read second RSSI register.\nThe signal strength of the most recently received LoRa packet.",
		},
	}

	for _, c := range commands {
		SendCommand(c, e90)
	}
}

func SendCommand(comm E90Command, e90 *E90Device) {

	fmt.Println("Testing : ", comm.CommandName)
	fmt.Println(comm.Description)
	fmt.Println("Sending command: ", hex.EncodeToString(comm.CommandBytes))

	err := e90.SendUDPCommand(comm.CommandBytes)
	if err != nil {
		log.Println("Failed to send command: ", err)
	} else {
		fmt.Println("Command sent successfully")
	}

	if comm.SecondsTimeout == 0 {
		comm.SecondsTimeout = 5
	}
	response, err := e90.ReceiveUDPResponseWithTimeout(time.Duration(comm.SecondsTimeout) * time.Second)

	if err != nil {
		log.Println("Failed to receive response: ", err)
	} else {
		fmt.Println("Received response: ", hex.EncodeToString(response))
		fmt.Println("Response bytes: ", response)
		fmt.Println("Interpreted: ", comm.InterpretFunction(response))
	}
	fmt.Println("===")
}

// return is 4 bytes: C1 + [Address] + [Read Length] + [RSSI Values]
func interpretRSSIResponse(response []byte) string {
	if len(response) < 4 {
		return "Invalid response"
	}

	rssiValue := int(response[3])
	return fmt.Sprintf("RSSI: %d dBm", rssiValue/2)
}
