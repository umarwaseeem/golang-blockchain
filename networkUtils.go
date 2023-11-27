package main

import (
	"log"
	"net"
)

// Get local ip of this machine
func GetLocalIP() net.IP {
	// udp connection to google dns
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	// after connection done, we get the local ip to which the socket is bound
	return localAddr.IP
}
