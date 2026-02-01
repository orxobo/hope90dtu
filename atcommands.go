package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//go:generate stringer -type=ATCmd -linecomment
type ATCmd int

const (
	// --- Basic ---
	CmdModel   ATCmd = iota // MODEL
	CmdName                 // NAME
	CmdSn                   // SN
	CmdReboot               // REBT
	CmdRestore              // RESTORE
	CmdVer                  // VER
	CmdMac                  // MAC

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

// Data Types

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

// Client Definition

type ATClient struct {
	device      *E90Device
	RegistryMap map[string]ATCmd
}

func NewATClient(device *E90Device) *ATClient {

	// Build a registry map for ease of searching and speed
	registryMap := map[string]ATCmd{}
	for _, meta := range CommandRegistry {
		registryMap[meta.Cmd.String()] = meta.Cmd //TODO: impliment this for searching, ie if v,ok := c.RegistryMap["mykey"]:ok {}
	}
	return &ATClient{device: device, RegistryMap: registryMap}
}

// Dynamic Command Registry

type CommandMeta struct {
	Cmd         ATCmd
	Description string
	Action      func(c *ATClient, args ...string) (any, error)
}

// CommandRegistry binds Enums to Methods.
var CommandRegistry = []CommandMeta{
	// --- Basic ---
	{CmdModel, "Query model", func(c *ATClient, a ...string) (any, error) { return c.GetModel() }},
	{CmdName, "Query/Set Name", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetName(a[0])
		}
		return c.GetName()
	}},
	{CmdSn, "Query/Set ID", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetSN(a[0])
		}
		return c.GetSN()
	}},
	{CmdReboot, "Reboot device", func(c *ATClient, a ...string) (any, error) { return c.Reboot() }},
	{CmdRestore, "Factory reset", func(c *ATClient, a ...string) (any, error) { return c.Restore() }},
	{CmdVer, "Query Version", func(c *ATClient, a ...string) (any, error) { return c.GetVersion() }},
	{CmdMac, "Query MAC", func(c *ATClient, a ...string) (any, error) { return c.GetMac() }},

	// --- Wireless ---
	{CmdLora, "Query/Set Lora (Arg: CSV)", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			p, err := ParseLora(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetLora(p)
		}
		return c.GetLora()
	}},
	{CmdWan, "Query/Set WAN (Arg: DHCP/STATIC...)", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			p, err := ParseWan(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetWan(p)
		}
		return c.GetWan()
	}},
	{CmdLPort, "Query/Set Local Port", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			v, err := strconv.Atoi(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetLocalPort(v)
		}
		return c.GetLocalPort()
	}},
	{CmdSock, "Query/Set Socket", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			p, err := ParseSocket(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetSocket(p)
		}
		return c.GetSocket()
	}},
	{CmdLinkSta, "Query Link Status", func(c *ATClient, a ...string) (any, error) { return c.GetLinkStatus() }},
	{CmdUartClr, "Query/Set UART Clear", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			v, err := strconv.ParseBool(a[0])
			if err != nil {
				return nil, err
			}
			return c.ClearUARTCache(v)
		}
		return execute(c, CmdUartClr, ParseBool)
	}},

	// --- Registration ---
	{CmdRegMod, "Query/Set Reg Mode", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetRegMode(a[0])
		}
		return c.GetRegMode()
	}},
	{CmdRegInfo, "Query/Set Reg Info", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetRegInfo(a[0])
		}
		return c.GetRegInfo()
	}},
	{CmdHeartMod, "Query/Set Heart Mode", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetHeartbeatMode(a[0])
		}
		return c.GetHeartbeatMode()
	}},
	{CmdHeartInfo, "Query/Set Heart Info", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetHeartbeatInfo(a[0])
		}
		return c.GetHeartbeatInfo()
	}},

	// --- Timing ---
	{CmdShortM, "Query/Set Short Conn", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			v, err := strconv.Atoi(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetShortConnection(v)
		}
		return c.GetShortConnection()
	}},
	{CmdTmoRst, "Query/Set Tmo Restart", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			v, err := strconv.Atoi(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetTimeoutRestart(v)
		}
		return c.GetTimeoutRestart()
	}},
	{CmdTmoLink, "Query/Set Tmo Link", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetTimeoutLink(a[0])
		}
		return c.GetTimeoutLink()
	}},
	{CmdWebCfgPort, "Query/Set Web Port", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			v, err := strconv.Atoi(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetWebConfigPort(v)
		}
		return c.GetWebConfigPort()
	}},

	// --- Modbus ---
	{CmdModWkMod, "Query/Set Modbus Work Mode", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetModbusWorkMode(a[0])
		}
		return c.GetModbusWorkMode()
	}},
	{CmdModPtcl, "Query/Set Modbus Protocol", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			v, err := strconv.ParseBool(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetModbusProtocol(v)
		}
		return c.GetModbusProtocol()
	}},
	{CmdModGtwyTm, "Query/Set Modbus Gateway", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetModbusGatewayTime(a[0])
		}
		return c.GetModbusGatewayTime()
	}},
	{CmdModCmdEdit, "Query/Set Modbus Cmd", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetModbusCmdEdit(a[0])
		}
		return c.GetModbusCmdEdit()
	}},

	// --- HTTP ---
	{CmdHtpReqMode, "Query/Set HTTP Req Mode", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetHttpRequestMode(a[0])
		}
		return c.GetHttpRequestMode()
	}},
	{CmdHtpUrl, "Query/Set HTTP URL", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetHttpUrl(a[0])
		}
		return c.GetHttpUrl()
	}},
	{CmdHtpHead, "Query/Set HTTP Header", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetHttpHeader(a[0])
		}
		return c.GetHttpHeader()
	}},

	// --- MQTT ---
	{CmdMqttCloud, "Query/Set MQTT Cloud", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetMqttCloud(a[0])
		}
		return c.GetMqttCloud()
	}},
	{CmdMqtKpAlive, "Query/Set Keep Alive", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			v, err := strconv.Atoi(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetMqttKeepAlive(v)
		}
		return c.GetMqttKeepAlive()
	}},
	{CmdMqtDevId, "Query/Set Device ID", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetMqttDeviceId(a[0])
		}
		return c.GetMqttDeviceId()
	}},
	{CmdMqtUser, "Query/Set MQTT User", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetMqttUser(a[0])
		}
		return c.GetMqttUser()
	}},
	{CmdMqtPass, "Query/Set MQTT Pass", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetMqttPass(a[0])
		}
		return c.GetMqttPass()
	}},
	{CmdMqtPrdKey, "Query/Set MQTT Prod Key", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			return c.SetMqttProductKey(a[0])
		}
		return c.GetMqttProductKey()
	}},
	{CmdMqtSub, "Query/Set Sub (Arg: Qos,Topic)", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			p, err := ParseMqttTopic(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetMqttSub(p)
		}
		return c.GetMqttSub()
	}},
	{CmdMqtPub, "Query/Set Pub (Arg: Qos,Topic)", func(c *ATClient, a ...string) (any, error) {
		if len(a) > 0 {
			p, err := ParseMqttTopic(a[0])
			if err != nil {
				return nil, err
			}
			return c.SetMqttPub(p)
		}
		return c.GetMqttPub()
	}},
}

// ListCommands prints the available commands to stdout using the Registry.
func (c *ATClient) ListCommands() {
	fmt.Println("\n--- Available AT Commands ---")
	for _, meta := range CommandRegistry {
		fmt.Printf("%-12s : %s\n", meta.Cmd, meta.Description)
	}
	fmt.Println("-----------------------------")
}

// Run executes a command using the Safe Enum.
// usage: client.Run(CmdLPort, "8080")
func (c *ATClient) Run(cmd ATCmd, args ...string) (any, error) {
	for _, meta := range CommandRegistry {
		if meta.Cmd == cmd {
			return meta.Action(c, args...)
		}
	}
	return nil, fmt.Errorf("command enum %d not registered", cmd)
}

// RunRaw executes a command by String name (useful for CLI/User Input).
// It iterates the enums to find a match.
// usage: client.RunRaw("LPORT", "8080")
func (c *ATClient) RunRaw(name string, args ...string) (any, error) {
	name = strings.ToUpper(strings.TrimSpace(name))
	for _, meta := range CommandRegistry {
		if meta.Cmd.String() == name {
			return meta.Action(c, args...)
		}
	}
	return nil, fmt.Errorf("command string '%s' not found", name)
}

// Internal Engine (Private)

// execute is the generic core. It accepts ATCmd (int) and uses .String() for the wire protocol.
func execute[T any](c *ATClient, cmd ATCmd, parser func(string) (T, error), args ...string) (T, error) {
	var zero T

	// Use stringer generated value (e.g., "MODEL" or "LPORT")
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

// Public Type-Safe API (Methods)

// --- Basic ---

func (c *ATClient) GetModel() (string, error)        { return execute(c, CmdModel, ParseString) }
func (c *ATClient) GetName() (string, error)         { return execute(c, CmdName, ParseString) }
func (c *ATClient) SetName(v string) (string, error) { return execute(c, CmdName, ParseString, v) }
func (c *ATClient) GetSN() (string, error)           { return execute(c, CmdSn, ParseString) }
func (c *ATClient) SetSN(v string) (string, error)   { return execute(c, CmdSn, ParseString, v) }
func (c *ATClient) Reboot() (bool, error)            { return execute(c, CmdReboot, ParseBool) }
func (c *ATClient) Restore() (bool, error)           { return execute(c, CmdRestore, ParseBool) }
func (c *ATClient) GetVersion() (string, error)      { return execute(c, CmdVer, ParseString) }
func (c *ATClient) GetMac() (string, error)          { return execute(c, CmdMac, ParseString) }

// --- Wireless & Network ---

func (c *ATClient) GetLora() (LoraParams, error) { return execute(c, CmdLora, ParseLora) }
func (c *ATClient) SetLora(p LoraParams) (LoraParams, error) {
	arg := fmt.Sprintf("%s,%d,%d,%d,%s,%d", p.Mode, p.Addr, p.Baud, p.Ch, p.PwrMode, p.Wor)
	return execute(c, CmdLora, ParseLora, arg)
}

func (c *ATClient) GetWan() (WanParams, error) { return execute(c, CmdWan, ParseWan) }
func (c *ATClient) SetWan(p WanParams) (WanParams, error) {
	mode := "STATIC"
	if p.IsDHCP {
		mode = "DHCP"
	}
	arg := fmt.Sprintf("%s,%s,%s,%s", mode, p.IP, p.Mask, p.GW)
	return execute(c, CmdWan, ParseWan, arg)
}

func (c *ATClient) GetLocalPort() (int, error) { return execute(c, CmdLPort, ParseInt) }
func (c *ATClient) SetLocalPort(v int) (int, error) {
	return execute(c, CmdLPort, ParseInt, strconv.Itoa(v))
}

func (c *ATClient) GetSocket() (SocketParams, error) { return execute(c, CmdSock, ParseSocket) }
func (c *ATClient) SetSocket(p SocketParams) (SocketParams, error) {
	arg := fmt.Sprintf("%s,%s,%d", p.Protocol, p.IP, p.Port)
	return execute(c, CmdSock, ParseSocket, arg)
}

func (c *ATClient) GetLinkStatus() (LinkStatus, error) {
	return execute(c, CmdLinkSta, ParseLinkStatus)
}

func (c *ATClient) ClearUARTCache(enable bool) (bool, error) {
	arg := "OFF"
	if enable {
		arg = "ON"
	}
	return execute(c, CmdUartClr, ParseBool, arg)
}

// --- Registration & Heartbeat ---

func (c *ATClient) GetRegMode() (string, error)         { return execute(c, CmdRegMod, ParseString) }
func (c *ATClient) SetRegMode(v string) (string, error) { return execute(c, CmdRegMod, ParseString, v) }
func (c *ATClient) GetRegInfo() (string, error)         { return execute(c, CmdRegInfo, ParseString) }
func (c *ATClient) SetRegInfo(v string) (string, error) {
	return execute(c, CmdRegInfo, ParseString, v)
}
func (c *ATClient) GetHeartbeatMode() (string, error) { return execute(c, CmdHeartMod, ParseString) }
func (c *ATClient) SetHeartbeatMode(v string) (string, error) {
	return execute(c, CmdHeartMod, ParseString, v)
}
func (c *ATClient) GetHeartbeatInfo() (string, error) { return execute(c, CmdHeartInfo, ParseString) }
func (c *ATClient) SetHeartbeatInfo(v string) (string, error) {
	return execute(c, CmdHeartInfo, ParseString, v)
}

// --- Timing & Connection ---

func (c *ATClient) GetShortConnection() (int, error) { return execute(c, CmdShortM, ParseInt) }
func (c *ATClient) SetShortConnection(v int) (int, error) {
	return execute(c, CmdShortM, ParseInt, strconv.Itoa(v))
}
func (c *ATClient) GetTimeoutRestart() (int, error) { return execute(c, CmdTmoRst, ParseInt) }
func (c *ATClient) SetTimeoutRestart(v int) (int, error) {
	return execute(c, CmdTmoRst, ParseInt, strconv.Itoa(v))
}
func (c *ATClient) GetTimeoutLink() (string, error) { return execute(c, CmdTmoLink, ParseString) }
func (c *ATClient) SetTimeoutLink(v string) (string, error) {
	return execute(c, CmdTmoLink, ParseString, v)
}
func (c *ATClient) GetWebConfigPort() (int, error) { return execute(c, CmdWebCfgPort, ParseInt) }
func (c *ATClient) SetWebConfigPort(v int) (int, error) {
	return execute(c, CmdWebCfgPort, ParseInt, strconv.Itoa(v))
}

// --- Modbus ---

func (c *ATClient) GetModbusWorkMode() (string, error) { return execute(c, CmdModWkMod, ParseString) }
func (c *ATClient) SetModbusWorkMode(v string) (string, error) {
	return execute(c, CmdModWkMod, ParseString, v)
}
func (c *ATClient) GetModbusProtocol() (bool, error) { return execute(c, CmdModPtcl, ParseBool) }
func (c *ATClient) SetModbusProtocol(enable bool) (bool, error) {
	arg := "OFF"
	if enable {
		arg = "ON"
	}
	return execute(c, CmdModPtcl, ParseBool, arg)
}
func (c *ATClient) GetModbusGatewayTime() (string, error) {
	return execute(c, CmdModGtwyTm, ParseString)
}
func (c *ATClient) SetModbusGatewayTime(v string) (string, error) {
	return execute(c, CmdModGtwyTm, ParseString, v)
}
func (c *ATClient) GetModbusCmdEdit() (string, error) { return execute(c, CmdModCmdEdit, ParseString) }
func (c *ATClient) SetModbusCmdEdit(v string) (string, error) {
	return execute(c, CmdModCmdEdit, ParseString, v)
}

// --- HTTP ---

func (c *ATClient) GetHttpRequestMode() (string, error) {
	return execute(c, CmdHtpReqMode, ParseString)
}
func (c *ATClient) SetHttpRequestMode(v string) (string, error) {
	return execute(c, CmdHtpReqMode, ParseString, v)
}
func (c *ATClient) GetHttpUrl() (string, error)         { return execute(c, CmdHtpUrl, ParseString) }
func (c *ATClient) SetHttpUrl(v string) (string, error) { return execute(c, CmdHtpUrl, ParseString, v) }
func (c *ATClient) GetHttpHeader() (string, error)      { return execute(c, CmdHtpHead, ParseString) }
func (c *ATClient) SetHttpHeader(v string) (string, error) {
	return execute(c, CmdHtpHead, ParseString, v)
}

// --- MQTT ---

func (c *ATClient) GetMqttCloud() (string, error) { return execute(c, CmdMqttCloud, ParseString) }
func (c *ATClient) SetMqttCloud(v string) (string, error) {
	return execute(c, CmdMqttCloud, ParseString, v)
}
func (c *ATClient) GetMqttKeepAlive() (int, error) { return execute(c, CmdMqtKpAlive, ParseInt) }
func (c *ATClient) SetMqttKeepAlive(v int) (int, error) {
	return execute(c, CmdMqtKpAlive, ParseInt, strconv.Itoa(v))
}
func (c *ATClient) GetMqttDeviceId() (string, error) { return execute(c, CmdMqtDevId, ParseString) }
func (c *ATClient) SetMqttDeviceId(v string) (string, error) {
	return execute(c, CmdMqtDevId, ParseString, v)
}
func (c *ATClient) GetMqttUser() (string, error) { return execute(c, CmdMqtUser, ParseString) }
func (c *ATClient) SetMqttUser(v string) (string, error) {
	return execute(c, CmdMqtUser, ParseString, v)
}
func (c *ATClient) GetMqttPass() (string, error) { return execute(c, CmdMqtPass, ParseString) }
func (c *ATClient) SetMqttPass(v string) (string, error) {
	return execute(c, CmdMqtPass, ParseString, v)
}
func (c *ATClient) GetMqttProductKey() (string, error) { return execute(c, CmdMqtPrdKey, ParseString) }
func (c *ATClient) SetMqttProductKey(v string) (string, error) {
	return execute(c, CmdMqtPrdKey, ParseString, v)
}

func (c *ATClient) GetMqttSub() (MqttTopic, error) { return execute(c, CmdMqtSub, ParseMqttTopic) }
func (c *ATClient) SetMqttSub(t MqttTopic) (MqttTopic, error) {
	arg := fmt.Sprintf("%d,%s", t.QoS, t.Topic)
	return execute(c, CmdMqtSub, ParseMqttTopic, arg)
}
func (c *ATClient) GetMqttPub() (MqttTopic, error) { return execute(c, CmdMqtPub, ParseMqttTopic) }
func (c *ATClient) SetMqttPub(t MqttTopic) (MqttTopic, error) {
	arg := fmt.Sprintf("%d,%s", t.QoS, t.Topic)
	return execute(c, CmdMqtPub, ParseMqttTopic, arg)
}

// Parsers

func ParseString(d string) (string, error) { return d, nil }
func ParseInt(d string) (int, error)       { return strconv.Atoi(d) }
func ParseBool(d string) (bool, error)     { return true, nil }

func ParseLinkStatus(d string) (LinkStatus, error) {
	return LinkStatus{Connected: strings.Contains(d, "Connect"), Msg: d}, nil
}

func ParseLora(d string) (LoraParams, error) {
	p := strings.Split(d, ",")
	if len(p) < 6 {
		return LoraParams{}, errors.New("bad lora format")
	}

	addr, _ := strconv.Atoi(p[1])
	baud, _ := strconv.Atoi(p[2])
	ch, _ := strconv.Atoi(p[3])
	wor, _ := strconv.Atoi(p[5])

	pwr := "PWMAX"
	if len(p) > 6 {
		pwr = p[6]
	}

	return LoraParams{
		Mode: p[0], Addr: addr, Baud: baud, Ch: ch, PwrMode: pwr, Wor: wor,
	}, nil
}

func ParseWan(d string) (WanParams, error) {
	p := strings.Split(d, ",")
	if len(p) < 4 {
		return WanParams{}, errors.New("bad wan format")
	}
	return WanParams{
		IsDHCP: p[0] == "DHCP", IP: p[1], Mask: p[2], GW: p[3],
	}, nil
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
		if idx := strings.Index(raw, "="); idx != -1 {
			return strings.TrimSpace(raw[idx+1:]), nil
		}
		return "", nil
	}

	return "", fmt.Errorf("invalid response: %s", raw)
}

// Transport / Helpers

const DEVICE_AT_CMD = "NETWK"

func (c *ATClient) sendATCommand(command string) error {
	command = fmt.Sprint(DEVICE_AT_CMD, "+", command, "\r\n")
	byteCommand := []byte(command)
	fmt.Printf(">> Sending: %s", strings.TrimSpace(command))
	return c.device.sendUDPCommand(byteCommand)
}

func (c *ATClient) receiveATResponse() (string, error) {
	rawBytes, err := c.device.receiveUDPResponseWithTimeout(4 * time.Second)
	if err != nil {
		return "", fmt.Errorf("no response: %w", err)
	}
	fmt.Printf(" << Received: %s\n", strings.TrimSpace(string(rawBytes)))
	return string(rawBytes), nil
}
