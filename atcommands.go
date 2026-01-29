package main

import (
	"fmt"
	"strconv"
	"strings"
)

// === ATError ===

type ATError struct {
	Code    int
	Message string
	Command string
}

func (e ATError) Error() string {
	if e.Command != "" {
		return fmt.Sprintf("AT command '%s' error %d: %s", e.Command, e.Code, e.Message)
	}
	return fmt.Sprintf("AT error %d: %s", e.Code, e.Message)
}

func NewATError(code int, command string) ATError {
	message := errorCodes[code]
	if message == "" {
		message = "Unknown error"
	}
	return ATError{
		Code:    code,
		Message: message,
		Command: command,
	}
}

var errorCodes = map[int]string{
	-1: "invalid command format",
	-2: "invalid command",
	-3: "Not yet defined",
	-4: "invalid parameter",
	-5: "Not yet defined",
}

// === Core AT Response Parser ===

type ATResponse struct {
	Raw     string
	Success bool
	Data    string
}

func ParseResponse(raw string) ATResponse {
	raw = strings.TrimSpace(raw)
	resp := ATResponse{Raw: raw}

	if strings.HasPrefix(raw, "+OK") {
		resp.Success = true
		if idx := strings.Index(raw, "="); idx != -1 {
			resp.Data = strings.TrimSpace(raw[idx+1:])
		}
	}

	return resp
}

// === Generic Command Type ===

type Command[T any] struct {
	Name        string
	Description string
	Parser      func(string) T
}

// Execute command and parse response
func (c Command[T]) Execute(response string) (T, error) {
	var zero T
	resp := ParseResponse(response)
	if !resp.Success {
		return zero, fmt.Errorf("command failed: %s", resp.Raw)
	}
	return c.Parser(resp.Data), nil
}

// === Parser Factories ===

func StringParser() func(string) string {
	return func(data string) string {
		return data
	}
}

func IntParser() func(string) int {
	return func(data string) int {
		val, _ := strconv.Atoi(data)
		return val
	}
}

func BoolParser() func(string) bool {
	return func(data string) bool {
		return true // +OK means success
	}
}

func ModelParser() func(string) Model {
	return func(data string) Model {
		return Model{Name: data}
	}
}

func LORAParser() func(string) LORAParams {
	return func(data string) LORAParams {
		var params LORAParams
		parts := strings.Split(data, ",")
		if len(parts) >= 8 {
			params.Address, _ = strconv.Atoi(parts[0])
			params.NetworkID, _ = strconv.Atoi(parts[1])
			params.AirDataRate, _ = strconv.Atoi(parts[2])
			params.PacketLength, _ = strconv.Atoi(parts[3])
			params.Channel, _ = strconv.Atoi(parts[6])
		}
		return params
	}
}

func NetworkParser() func(string) NetworkParams {
	return func(data string) NetworkParams {
		var params NetworkParams
		parts := strings.Split(data, ",")
		if len(parts) >= 4 {
			params.Mode = strings.TrimSpace(parts[0])
			params.IP = strings.TrimSpace(parts[1])
			params.Gateway = strings.TrimSpace(parts[3])
		}
		return params
	}
}

func ModeParser() func(string) ModeParams {
	return func(data string) ModeParams {
		var params ModeParams
		parts := strings.Split(data, ",")
		if len(parts) >= 3 {
			params.Mode = strings.TrimSpace(parts[0])
			params.RemoteIP = strings.TrimSpace(parts[1])
			params.RemotePort, _ = strconv.Atoi(strings.TrimSpace(parts[2]))
		}
		return params
	}
}

func TopicParser() func(string) TopicParams {
	return func(data string) TopicParams {
		var params TopicParams
		parts := strings.Split(data, ",")
		if len(parts) >= 2 {
			params.QoS, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
			params.Topic = strings.TrimSpace(parts[1])
		}
		return params
	}
}

// === Data Types ===

type Model struct{ Name string }

type LORAParams struct {
	Address, NetworkID, AirDataRate, PacketLength, Channel int
}

type NetworkParams struct {
	Mode, IP, Gateway string
}

type ModeParams struct {
	Mode, RemoteIP string
	RemotePort     int
}

type TopicParams struct {
	QoS   int
	Topic string
}

// === Command Registry ===

type CommandRegistry struct {
	commands map[string]interface{}
}

func NewCommandRegistry() *CommandRegistry {
	reg := &CommandRegistry{
		commands: make(map[string]interface{}),
	}

	reg.commands["AT+EXAT"] = Command[bool]{
		Name:        "AT+EXAT",
		Description: "Exit AT mode",
		Parser:      BoolParser(),
	}

	reg.commands["AT+MODEL"] = Command[Model]{
		Name:        "AT+MODEL",
		Description: "Query model",
		Parser:      ModelParser(),
	}

	reg.commands["AT+NAME"] = Command[string]{
		Name:        "AT+NAME",
		Description: "Query/set name",
		Parser:      StringParser(),
	}

	reg.commands["AT+SN"] = Command[string]{
		Name:        "AT+SN",
		Description: "Query/set ID",
		Parser:      StringParser(),
	}

	reg.commands["AT+REBT"] = Command[bool]{
		Name:        "AT+REBT",
		Description: "Reboot device",
		Parser:      BoolParser(),
	}

	reg.commands["AT+LORA"] = Command[LORAParams]{
		Name:        "AT+LORA",
		Description: "LORA parameters",
		Parser:      LORAParser(),
	}

	reg.commands["AT+WAN"] = Command[NetworkParams]{
		Name:        "AT+WAN",
		Description: "Network parameters",
		Parser:      NetworkParser(),
	}

	reg.commands["AT+LPORT"] = Command[int]{
		Name:        "AT+LPORT",
		Description: "Local port",
		Parser:      IntParser(),
	}

	reg.commands["AT+SOCK"] = Command[ModeParams]{
		Name:        "AT+SOCK",
		Description: "Working mode",
		Parser:      ModeParser(),
	}

	reg.commands["AT+MQTTCLOUD"] = Command[string]{
		Name:        "AT+MQTTCLOUD",
		Description: "MQTT platform",
		Parser:      StringParser(),
	}

	reg.commands["AT+MQTSUB"] = Command[TopicParams]{
		Name:        "AT+MQTSUB",
		Description: "MQTT subscription",
		Parser:      TopicParser(),
	}

	reg.commands["AT+MQTPUB"] = Command[TopicParams]{
		Name:        "AT+MQTPUB",
		Description: "MQTT publish",
		Parser:      TopicParser(),
	}

	return reg
}

// Generic Execute function (as a standalone function)
func Execute[T any](registry *CommandRegistry, name, response string) (T, error) {
	cmdInterface, exists := registry.commands[name]
	if !exists {
		var zero T
		return zero, fmt.Errorf("unknown command: %s", name)
	}

	// Type assert to the correct command type
	cmd, ok := cmdInterface.(Command[T])
	if !ok {
		var zero T
		return zero, fmt.Errorf("type mismatch for command: %s", name)
	}

	return cmd.Execute(response)
}
