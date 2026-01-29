package main

import (
	"net"
	"net/netip"
	"time"
)

type E90Device struct {
	conn          *net.UDPConn
	ATCommandMode bool
}

func NewE90Device(device netip.AddrPort) (*E90Device, error) {
	addr := net.UDPAddrFromAddrPort(device)

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &E90Device{conn, false}, nil
}

func (e *E90Device) sendUDPCommand(command []byte) error {
	_, err := e.conn.Write(command)
	return err
}

func (e *E90Device) sendUDPASCIICommand(command string) error {
	return e.sendUDPCommand([]byte(command))
}

func (e *E90Device) receiveUDPResponseWithTimeout(timeout time.Duration) ([]byte, error) {
	e.conn.SetReadDeadline(time.Now().Add(timeout))

	buffer := make([]byte, 1024)
	n, err := e.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}
