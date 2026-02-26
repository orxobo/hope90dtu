// Usage:
// client := NewATClient(e90)
// -
// // 1. Direct Method Usage
// model, err := client.GetModel()
//
//	if err != nil {
//		fmt.Println("error getting model: ", err.Error())
//	}
//
// fmt.Println("Model: ", model)
// -
// // 2. Enum Usage
// port, err := client.Run(CmdLPort)
//
//	if err != nil {
//		fmt.Println("error getting port: ", err.Error())
//	}
//
// fmt.Println("Port: ", port)
// -
// // 3. Raw String Usage (Relies on stringer generated strings)
// //// -
// client.ListCommands()
package atcommands

import (
	"fmt"
	"hope90dtu/device"
	"strings"
	"time"
	"unicode"
)

type ATClient struct {
	device         *device.E90Device
	deviceATPrefix string
}

// Initialises the AT command set handler
func NewATClient(device *device.E90Device, prefix string) (*ATClient, error) {

	var result strings.Builder
	for _, currentRune := range prefix {
		if unicode.IsLetter(currentRune) {
			result.WriteRune(currentRune)
			continue
		}
	}
	if result.Len() == 0 {
		return nil, fmt.Errorf("no AT prefix defined")
	}

	if device == nil {
		return nil, fmt.Errorf("device not initialised")
	}

	return &ATClient{device: device, deviceATPrefix: result.String()}, nil
}

// Run executes an AT command
// args are for saving the parameters to the device
func (c *ATClient) Run(cmd ATCmd, args ...string) (any, error) {
	atCommand, err := GetCommand(cmd)
	if err != nil {
		return nil, err
	}
	return atCommand.Action(c, args...)
}

// execute is the generic core. It accepts ATCmd and uses .String() for the wire protocol.
func execute[T any](c *ATClient, cmd ATCmd, parser func(string) (T, error), args ...string) (T, error) {
	var zero T

	fullCmd := cmd.String()
	if len(args) > 0 && args[0] != "" {
		fullCmd = fmt.Sprintf("%s=%s", fullCmd, args[0])
	}

	if err := c.sendATCommand(fullCmd); err != nil {
		return zero, err
	}

	rawBytes, err := c.receiveATResponse()
	if err != nil {
		return zero, err
	}

	payload, err := parseProtocolHeader(rawBytes)
	if err != nil {
		return zero, err
	}

	return parser(payload)
}

func (c *ATClient) sendATCommand(command string) error {
	command = fmt.Sprint(c.deviceATPrefix, "+", command, "\r\n")
	return c.device.SendUDPASCIICommand(command)
}

func (c *ATClient) receiveATResponse() (string, error) {
	rawBytes, err := c.device.ReceiveUDPResponseWithTimeout(4 * time.Second)
	if err != nil {
		return "", fmt.Errorf("no response: %w", err)
	}
	return string(rawBytes), nil
}
