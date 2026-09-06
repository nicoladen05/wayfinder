package main

import (
	"fmt"
	"net"
)

func main() {
	address := net.UDPAddr{
		IP:   net.IP{127, 0, 0, 1},
		Port: 1053,
	}

	listener, err := net.ListenUDP("udp", &address)

	if err != nil {
		fmt.Println("Error starting UDP server:", err)
	}

	for {
		buf := make([]byte, 1024)

		n, err := listener.Read(buf)

		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
		}

		fmt.Println(parseHeader(buf[:12]).string())

		question, _ := parseQuestion(buf[12:n])
		fmt.Println(question.string())
	}
}
