package atcommands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// TODO: add enums
type LoraParams struct {
	ModuleAddress     int    // [0]Module Address = 0-65535 (65535 = broadcast)
	NetId             int    // [1]Net ID = 0-255
	AirBaud           int    // [2]Air Baud in bits/second = 2400, 4800, 9600, 19200, 38400, 62500
	PacketLength      int    // [3]Packet Length in bytes = 32, 64, 128, 240
	ChannelRssi       bool   // [4]RSSI for the channel = RSCHON, RSCHOFF
	TransmissionPower string // [5]Transmission Power = PWMAX(hight),PWMID(middle),PWLOW(low),PWMIN(lowest)
	Channel           int    // [6]Channel = 0-80 for 900 series (850.125Mhz - 930.125Mhz in 1Mhz increments per channel)
	PacketRssi        bool   // [7]RSSI for the packet = RSDATON, RSDATOFF
	TransmissionMode  string // [8]Transmission Mode = TRNOR(normal), TRFIX(fixed)
	RelayEnable       bool   // [9]Relay Enable = RLYOFF, RLYON
	LBTEnable         bool   // [10]LBT Enable = LBTOFF, LBTON
	WorRole           string // [11]WOR Role = WOROFF, WORRX, WORTX
	WorCycle          int    // [12]WOR Cycle in miliseconds = 500-4000 in 500ms increments
	Key               int    // [13]KEY = 0-65535
}

// Parsers
func ParseString(d string) (string, error) { return d, nil }
func ParseInt(d string) (int, error)       { return strconv.Atoi(d) }
func ParseBool(d string) (bool, error)     { return strconv.ParseBool(d) }

// example: +OK=65535,0,2400,240,RSCHON,PWMAX,76,RSDATON,TRNOR,RLYOFF,LBTOFF,WOROFF,2000,0
func ParseLora(d string) (LoraParams, error) {
	params := strings.Split(d, ",")
	if len(params) < 13 {
		return LoraParams{}, errors.New("bad lora format")
	}

	loraparams := LoraParams{}

	var err error
	loraparams.ModuleAddress, err = strconv.Atoi(params[0])
	if err != nil {
		return LoraParams{}, errors.New("bad module address")
	}

	loraparams.NetId, err = strconv.Atoi(params[1])
	if err != nil {
		return LoraParams{}, errors.New("bad net ID")
	}

	loraparams.AirBaud, err = strconv.Atoi(params[2])
	if err != nil {
		return LoraParams{}, errors.New("bad air baud")
	}

	loraparams.PacketLength, err = strconv.Atoi(params[3])
	if err != nil {
		return LoraParams{}, errors.New("bad packet length")
	}

	loraparams.ChannelRssi = (params[4] == "RSCHON")
	loraparams.TransmissionPower = params[5]

	loraparams.Channel, err = strconv.Atoi(params[6])
	if err != nil {
		return LoraParams{}, errors.New("bad channel")
	}

	loraparams.PacketRssi = (params[7] == "RSDATON")
	loraparams.TransmissionMode = params[8]
	loraparams.RelayEnable = (params[9] == "RLYON")
	loraparams.LBTEnable = (params[10] == "LBTON")
	loraparams.WorRole = params[11]

	loraparams.WorCycle, err = strconv.Atoi(params[12])
	if err != nil {
		return LoraParams{}, errors.New("bad wor cycle")
	}

	loraparams.Key, err = strconv.Atoi(params[13])
	if err != nil {
		return LoraParams{}, errors.New("bad key")
	}

	return loraparams, nil
}

func parseProtocolHeader(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	if strings.Contains(raw, "+ERR=") {
		parts := strings.Split(raw, "=")
		code := 0
		if len(parts) > 1 {
			code, _ = strconv.Atoi(parts[1])
		}
		return "", NewATError(code, "")
	}

	if strings.HasPrefix(raw, "+OK") {
		if _, after, ok := strings.Cut(raw, "="); ok {
			return strings.TrimSpace(after), nil
		}
		return "", nil
	}

	return "", fmt.Errorf("invalid response: %s", raw)
}
