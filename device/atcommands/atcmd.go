// Complete list of AT Commands accepted by the E90
package atcommands

import "strings"

//go:generate stringer -type=ATCmd -linecomment
type ATCmd int

const (
	CmdInvalid ATCmd = iota // INVALID

	// --- Basic ---
	CmdModel   // MODEL
	CmdName    // NAME
	CmdSn      // SN
	CmdReboot  // REBT
	CmdRestore // RESTORE
	CmdVer     // VER
	CmdMac     // MAC

	// --- Wireless & Network ---
	CmdLora    // LORA
	CmdWan     // WAN
	CmdLPort   // LPORT
	CmdSock    // SOCK
	CmdLinkSta // LINKSTA
	CmdUartClr // UARTCLR

	// --- Registration & Heartbeat ---
	CmdRegMod    // REGMOD
	CmdRegInfo   // REGINFO
	CmdHeartMod  // HEARTMOD
	CmdHeartInfo // HEARTINFO

	// --- Timing ---
	CmdShortM     // SHORTM
	CmdTmoRst     // TMORST
	CmdTmoLink    // TMOLINK
	CmdWebCfgPort // WEBCFGPORT

	// --- Modbus ---
	CmdModWkMod   // MODWKMOD
	CmdModPtcl    // MODPTCL
	CmdModGtwyTm  // MODGTWYTM
	CmdModCmdEdit // MODCMDEDIT

	// --- HTTP ---
	CmdHtpReqMode // HTPREQMODE
	CmdHtpUrl     // HTPURL
	CmdHtpHead    // HTPHEAD

	// --- MQTT ---
	CmdMqttCloud  // MQTTCLOUD
	CmdMqtKpAlive // MQTKPALIVE
	CmdMqtDevId   // MQTDEVID
	CmdMqtUser    // MQTUSER
	CmdMqtPass    // MQTPASS
	CmdMqtPrdKey  // MQTTPRDKEY
	CmdMqtSub     // MQTSUB
	CmdMqtPub     // MQTPUB
)

// For UI
func ATCmds() []string {
	returnSlice := []string{}
	for i := CmdInvalid; i <= CmdMqtPub; i++ {
		returnSlice = append(returnSlice, i.String())
	}
	return returnSlice
}

// uses a map, as just returning it when found does not seem as elequent and more readable
func ATCmdFromString(atCmdString string) ATCmd {
	m := make(map[string]ATCmd)
	for i, s := range ATCmds() {
		m[s] = ATCmd(i)
	}
	val, ok := m[strings.TrimSpace(atCmdString)]
	if ok {
		return val
	}
	return CmdInvalid
}
