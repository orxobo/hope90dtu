package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//go:generate stringer -type=ATCommand -linecomment
type ATCommand int

const (
	CMD_MODEL      ATCommand = iota // MODEL, Query model
	CMD_NAME                        // NAME, Query/Set Name
	CMD_SN                          // SN, Query/Set ID
	CMD_REBT                        // REBT, Reboot device
	CMD_RESTORE                     // RESTORE, Factory reset
	CMD_VER                         // VER, Query version information
	CMD_MAC                         // MAC, Querying the MAC address
	CMD_LORA                        // LORA, Query/set the wireless parameters of the machine
	CMD_WAN                         // WAN, Query/set network parameters
	CMD_LPORT                       // LPORT, Query/set the local port number
	CMD_SOCK                        // SOCK, Query/set the working mode of the machine and network parameters of the target device
	CMD_LINKSTA                     // LINKSTA, Query network link status
	CMD_UARTCLR                     // UARTCLR, Query/set serial port cache clearing status
	CMD_REGMOD                      // REGMOD, Query/Set Registration Package Mode
	CMD_REGINFO                     // REGINFO, Query/set custom registration package content
	CMD_HEARTMOD                    // HEARTMOD, Query/set the heartbeat packet mode
	CMD_HEARTINFO                   // HEARTINFO, Query/Set Heartbeat Data
	CMD_SHORTM                      // SHORTM, Query/Set Short Connection Time
	CMD_TMORST                      // TMORST, Query/set timeout restart time
	CMD_TMOLINK                     // TMOLINK, Query/set the time and times of disconnection and reconnection
	CMD_WEBCFGPORT                  // WEBCFGPORT, Web configuration port
	CMD_MODWKMOD                    // MODWKMOD, Query Modbus working mode and command timeout time
	CMD_MODPTCL                     // MODPTCL, Enable Modbus TCP to Modbus RTU protocol conversion
	CMD_MODGTWYTM                   // MODGTWYTM, Set Modbus gateway command storage time and automatic query interval
	CMD_MODCMDEDIT                  // MODCMDEDIT, Modbus configuration gateway pre-stored instruction query and editing
	CMD_HTPREQMODE                  // HTPREQMODE, Query/Set HTTP Request Method
	CMD_HTPURL                      // HTPURL, Query/Set HTTP URL Path
	CMD_HTPHEAD                     // HTPHEAD, Query/Set HTTP header
	CMD_MQTTCLOUD                   // MQTTCLOUD, Query/Set MQTT Target Platform
	CMD_MQTKPALIVE                  // MQTKPALIVE, Query/Set MQTT Keep-Alive Heartbeat Packet Sending Period
	CMD_MQTDEVID                    // MQTDEVID, Query/Set MQTT Device Name (Client ID)
	CMD_MQTUSER                     // MQTUSER, Query/Set MQTT User Name (User Name/Device Name)
	CMD_MQTPASS                     // MQTPASS, Query/set MQTT product password (MQTT password/Device Secret)
	CMD_MQTTPRDKEY                  // MQTTPRDKEY, Query/Set the Product Key of Alibaba Cloud MQTT
	CMD_MQTSUB                      // MQTSUB, Query/Set MQTT Subscription Topic
	CMD_MQTPUB                      // MQTPUB, Query/Set MQTT Publishing Topic
)

// CommandDef binds the String, Description, and Parser together
type CommandDef[T any] struct {
	CmdStr      string
	Description string
	Parser      func(string) (T, error)
}

// Registry: The single source of truth for Command -> Logic mapping
var Registry = map[ATCommand]any{
	// --- Basic Functions ---
	CMD_MODEL:   CommandDef[string]{"MODEL", "Query model", ParseString},
	CMD_NAME:    CommandDef[string]{"NAME", "Query/Set Name", ParseString},
	CMD_SN:      CommandDef[string]{"SN", "Query/Set ID", ParseString},
	CMD_REBT:    CommandDef[bool]{"REBT", "Reboot device", ParseBool},
	CMD_RESTORE: CommandDef[bool]{"RESTORE", "Factory reset", ParseBool},
	CMD_VER:     CommandDef[string]{"VER", "Query firmware version", ParseString},
	CMD_MAC:     CommandDef[string]{"MAC", "Query MAC address", ParseString},

	// --- Wireless & Network ---
	CMD_LORA:    CommandDef[LoraParams]{"LORA", "Query/set wireless parameters", ParseLora},
	CMD_WAN:     CommandDef[WanParams]{"WAN", "Query/set network parameters", ParseWan},
	CMD_LPORT:   CommandDef[int]{"LPORT", "Query/set local port number", ParseInt},
	CMD_SOCK:    CommandDef[SocketParams]{"SOCK", "Query/set target network parameters", ParseSocket},
	CMD_LINKSTA: CommandDef[LinkStatus]{"LINKSTA", "Query network link status", ParseLinkStatus},
	CMD_UARTCLR: CommandDef[bool]{"UARTCLR", "Query/set serial port cache clearing", ParseBool},

	// --- Registration & Heartbeat ---
	CMD_REGMOD:    CommandDef[string]{"REGMOD", "Query/Set Registration Package Mode", ParseString},
	CMD_REGINFO:   CommandDef[string]{"REGINFO", "Query/set custom registration content", ParseString},
	CMD_HEARTMOD:  CommandDef[string]{"HEARTMOD", "Query/set heartbeat packet mode", ParseString},
	CMD_HEARTINFO: CommandDef[string]{"HEARTINFO", "Query/Set Heartbeat Data", ParseString},

	// --- Timing & Connection ---
	CMD_SHORTM:     CommandDef[int]{"SHORTM", "Query/Set Short Connection Time", ParseInt},
	CMD_TMORST:     CommandDef[int]{"TMORST", "Query/set timeout restart time", ParseInt},
	CMD_TMOLINK:    CommandDef[string]{"TMOLINK", "Query/set disconnect/reconnect time", ParseString},
	CMD_WEBCFGPORT: CommandDef[int]{"WEBCFGPORT", "Web configuration port", ParseInt},

	// --- Modbus ---
	CMD_MODWKMOD:   CommandDef[string]{"MODWKMOD", "Query Modbus working mode", ParseString},
	CMD_MODPTCL:    CommandDef[bool]{"MODPTCL", "Enable Modbus TCP to RTU conversion", ParseBool},
	CMD_MODGTWYTM:  CommandDef[string]{"MODGTWYTM", "Set Modbus gateway timing", ParseString},
	CMD_MODCMDEDIT: CommandDef[string]{"MODCMDEDIT", "Modbus pre-stored instruction", ParseString},

	// --- HTTP & IoT ---
	CMD_HTPREQMODE: CommandDef[string]{"HTPREQMODE", "Query/Set HTTP Request Method", ParseString},
	CMD_HTPURL:     CommandDef[string]{"HTPURL", "Query/Set HTTP URL Path", ParseString},
	CMD_HTPHEAD:    CommandDef[string]{"HTPHEAD", "Query/Set HTTP header", ParseString},

	// --- MQTT ---
	CMD_MQTTCLOUD:  CommandDef[string]{"MQTTCLOUD", "Query/Set MQTT Target Platform", ParseString},
	CMD_MQTKPALIVE: CommandDef[int]{"MQTKPALIVE", "Query/Set MQTT Keep-Alive", ParseInt},
	CMD_MQTDEVID:   CommandDef[string]{"MQTDEVID", "Query/Set MQTT Device Name", ParseString},
	CMD_MQTUSER:    CommandDef[string]{"MQTUSER", "Query/Set MQTT User Name", ParseString},
	CMD_MQTPASS:    CommandDef[string]{"MQTPASS", "Query/set MQTT password", ParseString},
	CMD_MQTTPRDKEY: CommandDef[string]{"MQTTPRDKEY", "Query/Set Alibaba Product Key", ParseString},
	CMD_MQTSUB:     CommandDef[MqttTopic]{"MQTSUB", "Query/Set MQTT Subscription Topic", ParseMqttTopic},
	CMD_MQTPUB:     CommandDef[MqttTopic]{"MQTPUB", "Query/Set MQTT Publishing Topic", ParseMqttTopic},
}

type ATClient struct {
	device *E90Device
}

func NewATClient(device *E90Device) *ATClient {
	return &ATClient{device: device}
}

// Execute is the Generic entry point.
func Execute[T any](c *ATClient, cmd ATCommand, args ...string) (T, error) {
	var zero T

	defInterface, exists := Registry[cmd]
	if !exists {
		return zero, fmt.Errorf("command enum %d not found in registry", cmd)
	}

	def, ok := defInterface.(CommandDef[T])
	if !ok {
		return zero, fmt.Errorf("type mismatch: registry expects %T, caller requested %T", defInterface, zero)
	}

	fullCmd := def.CmdStr
	if len(args) > 0 && args[0] != "" {
		fullCmd = fmt.Sprintf("%s=%s", def.CmdStr, args[0])
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

	return def.Parser(payload)
}

func parseProtocolHeader(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	if strings.Contains(raw, "+ERR=") {
		parts := strings.Split(raw, "=")
		if len(parts) > 1 {
			code, _ := strconv.Atoi(parts[1])
			msg := "Unknown"
			switch code {
			case 1:
				msg = "Invalid format"
			case 2:
				msg = "Invalid command"
			case 3:
				msg = "Not defined"
			case 4:
				msg = "Invalid parameter"
			case 5:
				msg = "Not defined"
			}
			return "", fmt.Errorf("AT Error %d: %s", code, msg)
		}
		return "", errors.New("malformed error response")
	}

	if strings.HasPrefix(raw, "+OK") {
		if idx := strings.Index(raw, "="); idx != -1 {
			return strings.TrimSpace(raw[idx+1:]), nil
		}
		return "", nil
	}

	return "", fmt.Errorf("invalid response: %s", raw)
}

// --- Specific Data Types ---

type LoraParams struct {
	Mode, PwrMode       string
	Addr, Baud, Ch, Wor int
}

type WanParams struct {
	IsDHCP            bool
	IP, Mask, GW, DNS string
}

type SocketParams struct {
	Protocol, IP string
	Port         int
}

type MqttTopic struct {
	QoS   int
	Topic string
}

type LinkStatus struct {
	Connected bool
	Msg       string
}

// --- Parsers ---

func ParseString(d string) (string, error) { return d, nil }
func ParseInt(d string) (int, error)       { return strconv.Atoi(d) }
func ParseBool(d string) (bool, error)     { return true, nil }

func ParseLinkStatus(d string) (LinkStatus, error) {
	return LinkStatus{Connected: strings.Contains(d, "Connect"), Msg: d}, nil
}

func ParseLora(d string) (LoraParams, error) {
	// e.g. MODNOR,0,2400,23,TRFIX,250,PWMAX
	p := strings.Split(d, ",")
	if len(p) < 7 {
		return LoraParams{}, errors.New("bad lora format")
	}
	addr, _ := strconv.Atoi(p[1])
	baud, _ := strconv.Atoi(p[2])
	ch, _ := strconv.Atoi(p[3])
	wor, _ := strconv.Atoi(p[5])
	return LoraParams{Mode: p[0], Addr: addr, Baud: baud, Ch: ch, PwrMode: p[4], Wor: wor}, nil
}

func ParseWan(d string) (WanParams, error) {
	p := strings.Split(d, ",")
	if len(p) < 4 {
		return WanParams{}, errors.New("bad wan format")
	}
	return WanParams{IsDHCP: p[0] == "DHCP", IP: p[1], Mask: p[2], GW: p[3]}, nil
}

func ParseSocket(d string) (SocketParams, error) {
	p := strings.Split(d, ",")
	if len(p) < 3 {
		return SocketParams{}, errors.New("bad sock format")
	}
	port, _ := strconv.Atoi(p[2])
	return SocketParams{Protocol: p[0], IP: p[1], Port: port}, nil
}

func ParseMqttTopic(d string) (MqttTopic, error) {
	p := strings.Split(d, ",")
	if len(p) < 2 {
		return MqttTopic{}, errors.New("bad mqtt format")
	}
	q, _ := strconv.Atoi(p[0])
	return MqttTopic{QoS: q, Topic: p[1]}, nil
}

// --- Helpers ---

const DEVICE_AT_CMD = "NETWK" // Prefix on E90 when requesting AT command access

// UDP helper for logging and formatting
func (c *ATClient) sendATCommand(command string) error {
	command = fmt.Sprint(DEVICE_AT_CMD, "+", command, "\r\n") // (0x0D and 0x0A)
	byteCommand := []byte(command)
	fmt.Printf("\nSending command: %s,\nBytes: %v,\nHex: %v\n\n",
		strings.TrimSpace(command),
		byteCommand,
		hex.EncodeToString(byteCommand))
	return c.device.sendUDPCommand(byteCommand)
}

// UDP helper for logging and formatting
func (c *ATClient) receiveATResponse() (string, error) {

	rawBytes, err := c.device.receiveUDPResponseWithTimeout(4 * time.Second)
	if err != nil {
		return "", fmt.Errorf("No response received: %w", err)
	}

	fmt.Printf("Received: %s,\nHex: %v\n\n", strings.TrimSpace(string(rawBytes)), hex.EncodeToString(rawBytes))
	return string(rawBytes), nil
}
