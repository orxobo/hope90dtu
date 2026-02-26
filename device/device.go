package device

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"
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

// NewE90UDPDevice establishes a UDP connection to the E90 via UDP4
func NewE90UDPDevice(device netip.AddrPort) (*E90Device, error) {
	addr := net.UDPAddrFromAddrPort(device)

	connection, err := net.DialUDP("udp4", nil, addr)
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

func SerialPorts() []string {
	return []string{"COM1", "COM2", "COM3", "COM4", "COM5", "COM6",
		"/dev/ttyUSB0", "/dev/ttyACM0", "/dev/ttyS0"}
}

func BaudRates() []string {
	return []string{"9600", "19200", "38400", "57600", "115200"}
}

// NewE90SerialDevice is a future functionality placeholder
func NewE90SerialDevice(serialPort string, baudRate int) (*E90Device, error) {
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

	e.sendToMonitor(buffer[:n], false)
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
	monitorMessage := "Received"
	if sending {
		monitorMessage = "Sending"
	}
	monitorMessage = fmt.Sprintf("%s:\n[RAW]:%v,\n[HEX]:%X,\n[STR]: %s", monitorMessage, message, message, strings.TrimRight(string(message), "\r\n"))
	e.monitor(monitorMessage)
}

func (e *E90Device) SendRandomData(length int) (int, error) {
	if length < 0 {
		return 0, fmt.Errorf("length cannot be negative")
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return 0, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	err = e.SendUDPCommand(bytes)
	if err != nil {
		return 0, err
	}

	resp, err := e.ReceiveUDPResponse()
	if err != nil {
		return 0, err
	}
	return len(resp), nil
}

func (e *E90Device) UDPListener(ctx context.Context) {

	// irrelevant data to initiate UDP connection with E90 as it will only retrun data on the port that sent it data
	e.conn.Write([]byte("INIT"))

	for {
		select {
		case <-ctx.Done():
			return
		default:
			buffer := make([]byte, 1024)
			e.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			n, err := e.conn.Read(buffer)
			if err != nil {
				continue
			}
			if n > 0 {
				data := make([]byte, n)
				copy(data, buffer[:n])
				go e.sendToMonitor(data, false)
			}
		}

	}
}
