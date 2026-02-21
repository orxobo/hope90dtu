package device

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"time"
)

type E90Device struct {
	conn               *net.UDPConn
	disconnectCallBack func()
	monitor            func(content string)
}

// NewE90UDPDeviceFromIPAddressAndPort is a helper for ease of instacation from UI
func NewE90UDPDeviceFromIPAddressAndPort(ip string, port string) (*E90Device, error) {
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, fmt.Errorf("invalid ip address: %s, %w", ip, err)
	}

	convertedPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("%s is invalid for a port number: %w", port, err)
	}

	if convertedPort < 1 || convertedPort > 65535 {
		return nil, fmt.Errorf("port out of range")
	}

	deviceAddressPort := netip.AddrPortFrom(
		ipAddr,
		uint16(convertedPort),
	)
	return NewE90UDPDevice(deviceAddressPort)
}

func NewE90UDPDevice(device netip.AddrPort) (*E90Device, error) {
	addr := net.UDPAddrFromAddrPort(device)

	connection, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &E90Device{
		conn:               connection,
		disconnectCallBack: func() {},
		monitor: func(content string) {
			fmt.Println(content)
		},
	}, nil
}

// NewE90SerialDevice is a future functionality placeholder
func NewE90SerialDevice(serialPort string, baudRate int) (*E90Device, error) {
	baudrates := []string{"9600", "19200", "38400", "57600", "115200"}
	serialPorts := []string{"COM1", "COM2", "COM3", "COM4", "COM5", "COM6",
		"/dev/ttyUSB0", "/dev/ttyACM0", "/dev/ttyS0"}

	_ = baudrates
	_ = serialPorts
	// TODO: add serial connection to E90, github.com/hawkli-1994/serio looks good
	return nil, fmt.Errorf("not implemented")
}

func (e *E90Device) SendUDPCommand(command []byte) error {
	_, err := e.conn.Write(command)
	e.sendToMonitor(command, true)
	return err
}

func (e *E90Device) SendUDPASCIICommand(command string) error {
	return e.SendUDPCommand([]byte(command))
}

func (e *E90Device) SendUDPHexCommand(command string) error {
	hexValue, err := hex.DecodeString(command)
	if err != nil {
		return err
	}
	return e.SendUDPCommand(hexValue)
}

// ReceiveUDPResponse waits 5 seconds for a response before timing out
func (e *E90Device) ReceiveUDPResponse() ([]byte, error) {
	return e.ReceiveUDPResponseWithTimeout(5 * time.Second)
}

// ReceiveUDPResponseWithTimeout waits for a response for a given set of time
func (e *E90Device) ReceiveUDPResponseWithTimeout(timeout time.Duration) ([]byte, error) {
	e.conn.SetReadDeadline(time.Now().Add(timeout))

	buffer := make([]byte, 1024)
	n, err := e.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	e.sendToMonitor(buffer, false)
	return buffer[:n], nil
}

func (e *E90Device) Close() {
	e.conn.Close()
	e.disconnectCallBack()
}

func (e *E90Device) SetDisconnectCallback(disconnected func()) {
	e.disconnectCallBack = disconnected
}

func (e *E90Device) SetMonitor(monitor func(content string)) {
	e.monitor = monitor
}

// Send to monitor
func (e *E90Device) sendToMonitor(message []byte, sending bool) {
	monitorMessage := "Received @"
	if sending {
		monitorMessage = "Sending @"
	}
	timestamp := time.Now().Format("15:04:05.000")
	monitorMessage = fmt.Sprintf("%s %s [RAW]:%s, [HEX]:%s, [STR]: %s\n", monitorMessage, timestamp, message, hex.EncodeToString(message), string(message))
	e.monitor(monitorMessage)
}
