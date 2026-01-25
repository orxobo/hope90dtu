package main

import (
	"net"
	"time"
)

func connectToUDPDevice(device Device) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", device.IP+":"+device.Port)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func sendUDPCommand(conn *net.UDPConn, command []byte) error {
	_, err := conn.Write(command)
	return err
}

func receiveUDPResponse(conn *net.UDPConn, timeout time.Duration) ([]byte, error) {
	// Set a timeout for reading
	conn.SetReadDeadline(time.Now().Add(timeout))

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}
