package atcommands

import (
	"fmt"
	"strconv"
)

func (c *ATClient) GetModel() (string, error) { return execute(c, CmdModel, ParseString) }

func (c *ATClient) GetName() (string, error) { return execute(c, CmdName, ParseString) }

func (c *ATClient) SetName(v string) (string, error) { return execute(c, CmdName, ParseString, v) }

func (c *ATClient) GetSN() (string, error) { return execute(c, CmdSn, ParseString) }

func (c *ATClient) SetSN(v string) (string, error) { return execute(c, CmdSn, ParseString, v) }

func (c *ATClient) Reboot() (bool, error) { return execute(c, CmdReboot, ParseBool) }

func (c *ATClient) Restore() (bool, error) { return execute(c, CmdRestore, ParseBool) }

func (c *ATClient) GetVersion() (string, error) { return execute(c, CmdVer, ParseString) }

func (c *ATClient) GetMac() (string, error) { return execute(c, CmdMac, ParseString) }

func (c *ATClient) GetLora() (LoraParams, error) { return execute(c, CmdLora, ParseLora) }

func (c *ATClient) SetLora(p LoraParams) (LoraParams, error) {
	arg := fmt.Sprintf("%d,%d,%d,%d,%v,%s,%d,%v,%s,%v,%v,%s,%d,%d",
		p.ModuleAddress,
		p.NetId,
		p.AirBaud,
		p.PacketLength,
		p.ChannelRssi,
		p.TransmissionPower,
		p.Channel,
		p.PacketRssi,
		p.TransmissionMode,
		p.RelayEnable,
		p.LBTEnable,
		p.WorRole,
		p.WorCycle,
		p.Key,
	)
	return execute(c, CmdLora, ParseLora, arg)
}

func (c *ATClient) GetWan() (string, error) { return execute(c, CmdWan, ParseString) }

func (c *ATClient) SetWan(v string) (string, error) {
	return execute(c, CmdWan, ParseString, v)
}

func (c *ATClient) GetLocalPort() (int, error) { return execute(c, CmdLPort, ParseInt) }

func (c *ATClient) SetLocalPort(v int) (int, error) {
	return execute(c, CmdLPort, ParseInt, strconv.Itoa(v))
}

func (c *ATClient) GetSocket() (string, error) { return execute(c, CmdSock, ParseString) }

func (c *ATClient) SetSocket(v string) (string, error) {
	return execute(c, CmdSock, ParseString, v)
}

func (c *ATClient) GetLinkStatus() (string, error) {
	return execute(c, CmdLinkSta, ParseString)
}

func (c *ATClient) ClearUARTCache(enable bool) (bool, error) {
	arg := "OFF"
	if enable {
		arg = "ON"
	}
	return execute(c, CmdUartClr, ParseBool, arg)
}

func (c *ATClient) GetRegMode() (string, error) { return execute(c, CmdRegMod, ParseString) }

func (c *ATClient) SetRegMode(v string) (string, error) { return execute(c, CmdRegMod, ParseString, v) }

func (c *ATClient) GetRegInfo() (string, error) { return execute(c, CmdRegInfo, ParseString) }

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

func (c *ATClient) GetHttpRequestMode() (string, error) {
	return execute(c, CmdHtpReqMode, ParseString)
}
func (c *ATClient) SetHttpRequestMode(v string) (string, error) {
	return execute(c, CmdHtpReqMode, ParseString, v)
}
func (c *ATClient) GetHttpUrl() (string, error) { return execute(c, CmdHtpUrl, ParseString) }

func (c *ATClient) SetHttpUrl(v string) (string, error) { return execute(c, CmdHtpUrl, ParseString, v) }

func (c *ATClient) GetHttpHeader() (string, error) { return execute(c, CmdHtpHead, ParseString) }

func (c *ATClient) SetHttpHeader(v string) (string, error) {
	return execute(c, CmdHtpHead, ParseString, v)
}

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

func (c *ATClient) GetMqttSub() (string, error) { return execute(c, CmdMqtSub, ParseString) }

func (c *ATClient) SetMqttSub(v string) (string, error) {
	return execute(c, CmdMqtSub, ParseString, v)
}

func (c *ATClient) GetMqttPub() (string, error) { return execute(c, CmdMqtPub, ParseString) }

func (c *ATClient) SetMqttPub(v string) (string, error) {
	return execute(c, CmdMqtPub, ParseString, v)
}
