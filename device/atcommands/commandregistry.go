// Link ATCmd to the description and associated function
package atcommands

import (
	"fmt"
	"strconv"
)

type ATCommand struct {
	Cmd         ATCmd
	Description string
	Action      func(c *ATClient, args ...string) (any, error)
}

func GetCommand(cmd ATCmd) (*ATCommand, error) {
	for _, atCommand := range CommandRegistry {
		if atCommand.Cmd == cmd {
			return &atCommand, nil
		}
	}
	return nil, fmt.Errorf("command %s not found", cmd.String())
}

var CommandRegistry = []ATCommand{
	{CmdModel, "Query model",
		func(c *ATClient, a ...string) (any, error) {
			return c.GetModel()
		}},
	{CmdName, "Query/Set Name",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetName(a[0])
			}
			return c.GetName()
		}},
	{CmdSn, "Query/Set ID",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetSN(a[0])
			}
			return c.GetSN()
		}},
	{CmdReboot, "Reboot device",
		func(c *ATClient, a ...string) (any, error) {
			return c.Reboot()
		}},
	{CmdRestore, "Factory reset",
		func(c *ATClient, a ...string) (any, error) {
			return c.Restore()
		}},
	{CmdVer, "Query Version",
		func(c *ATClient, a ...string) (any, error) {
			return c.GetVersion()
		}},
	{CmdMac, "Query MAC",
		func(c *ATClient, a ...string) (any, error) {
			return c.GetMac()
		}},
	{CmdLora, "Query/Set Lora (Arg: CSV)",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				p, err := ParseLora(a[0])
				if err != nil {
					return nil, err
				}
				return c.SetLora(p)
			}
			return c.GetLora()
		}},
	{CmdWan, "Query/Set WAN (Arg: DHCP/STATIC...)",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetWan(a[0])
			}
			return c.GetWan()
		}},
	{CmdLPort, "Query/Set Local Port",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				v, err := strconv.Atoi(a[0])
				if err != nil {
					return nil, err
				}
				return c.SetLocalPort(v)
			}
			return c.GetLocalPort()
		}},
	{CmdSock, "Query/Set Socket",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetSocket(a[0])
			}
			return c.GetSocket()
		}},
	{CmdLinkSta, "Query Link Status",
		func(c *ATClient, a ...string) (any, error) {
			return c.GetLinkStatus()
		}},
	{CmdUartClr, "Query/Set UART Clear",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				v, err := strconv.ParseBool(a[0])
				if err != nil {
					return nil, err
				}
				return c.ClearUARTCache(v)
			}
			return execute(c, CmdUartClr, ParseBool)
		}},
	{CmdRegMod, "Query/Set Reg Mode",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetRegMode(a[0])
			}
			return c.GetRegMode()
		}},
	{CmdRegInfo, "Query/Set Reg Info",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetRegInfo(a[0])
			}
			return c.GetRegInfo()
		}},
	{CmdHeartMod, "Query/Set Heart Mode",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetHeartbeatMode(a[0])
			}
			return c.GetHeartbeatMode()
		}},
	{CmdHeartInfo, "Query/Set Heart Info",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetHeartbeatInfo(a[0])
			}
			return c.GetHeartbeatInfo()
		}},
	{CmdShortM, "Query/Set Short Conn",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				v, err := strconv.Atoi(a[0])
				if err != nil {
					return nil, err
				}
				return c.SetShortConnection(v)
			}
			return c.GetShortConnection()
		}},
	{CmdTmoRst, "Query/Set Tmo Restart",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				v, err := strconv.Atoi(a[0])
				if err != nil {
					return nil, err
				}
				return c.SetTimeoutRestart(v)
			}
			return c.GetTimeoutRestart()
		}},
	{CmdTmoLink, "Query/Set Tmo Link",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetTimeoutLink(a[0])
			}
			return c.GetTimeoutLink()
		}},
	{CmdWebCfgPort, "Query/Set Web Port",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				v, err := strconv.Atoi(a[0])
				if err != nil {
					return nil, err
				}
				return c.SetWebConfigPort(v)
			}
			return c.GetWebConfigPort()
		}},
	{CmdModWkMod, "Query/Set Modbus Work Mode",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetModbusWorkMode(a[0])
			}
			return c.GetModbusWorkMode()
		}},
	{CmdModPtcl, "Query/Set Modbus Protocol",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				v, err := strconv.ParseBool(a[0])
				if err != nil {
					return nil, err
				}
				return c.SetModbusProtocol(v)
			}
			return c.GetModbusProtocol()
		}},
	{CmdModGtwyTm, "Query/Set Modbus Gateway",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetModbusGatewayTime(a[0])
			}
			return c.GetModbusGatewayTime()
		}},
	{CmdModCmdEdit, "Query/Set Modbus Cmd",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetModbusCmdEdit(a[0])
			}
			return c.GetModbusCmdEdit()
		}},
	{CmdHtpReqMode, "Query/Set HTTP Req Mode",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetHttpRequestMode(a[0])
			}
			return c.GetHttpRequestMode()
		}},
	{CmdHtpUrl, "Query/Set HTTP URL",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetHttpUrl(a[0])
			}
			return c.GetHttpUrl()
		}},
	{CmdHtpHead, "Query/Set HTTP Header",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetHttpHeader(a[0])
			}
			return c.GetHttpHeader()
		}},
	{CmdMqttCloud, "Query/Set MQTT Cloud",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetMqttCloud(a[0])
			}
			return c.GetMqttCloud()
		}},
	{CmdMqtKpAlive, "Query/Set Keep Alive",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				v, err := strconv.Atoi(a[0])
				if err != nil {
					return nil, err
				}
				return c.SetMqttKeepAlive(v)
			}
			return c.GetMqttKeepAlive()
		}},
	{CmdMqtDevId, "Query/Set Device ID",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetMqttDeviceId(a[0])
			}
			return c.GetMqttDeviceId()
		}},
	{CmdMqtUser, "Query/Set MQTT User",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetMqttUser(a[0])
			}
			return c.GetMqttUser()
		}},
	{CmdMqtPass, "Query/Set MQTT Pass",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetMqttPass(a[0])
			}
			return c.GetMqttPass()
		}},
	{CmdMqtPrdKey, "Query/Set MQTT Prod Key",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetMqttProductKey(a[0])
			}
			return c.GetMqttProductKey()
		}},
	{CmdMqtSub, "Query/Set Sub (Arg: Qos,Topic)",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetMqttSub(a[0])
			}
			return c.GetMqttSub()
		}},
	{CmdMqtPub, "Query/Set Pub (Arg: Qos,Topic)",
		func(c *ATClient, a ...string) (any, error) {
			if len(a) > 0 {
				return c.SetMqttPub(a[0])
			}
			return c.GetMqttPub()
		}},
}
