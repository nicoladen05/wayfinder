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

	conn, err := net.ListenUDP("udp", &address)

	if err != nil {
		fmt.Println("Error starting UDP server:", err)
	}

	for {
		buf := make([]byte, 1024)

		n, respAddr, err := conn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
		}

		header := parseHeader(buf[:12])
		fmt.Println(header.string())

		question, _ := parseQuestion(buf[12:n])
		fmt.Println(question.string())

		// Send a mock response
		responseHeader := DNSHeader{
			ID:      header.ID,
			QR:      1, // Response
			OPCODE:  header.OPCODE,
			AA:      0,
			TC:      0,
			RD:      header.RD,
			RA:      0,
			Z:       0,
			AD:      0,
			CD:      0,
			RCODE:   0,
			QDCOUNT: 0,
			ANCOUNT: 1,
			NSCOUNT: 0,
			ARCOUNT: 0,
		}

		responseBody := DNSResource{
			NAME:     []string{"www", "google", "com"},
			TYPE:     1,                  // A record
			CLASS:    1,                  // IN class
			TTL:      300,                // Time to live in seconds
			RDLENGTH: 4,                  // Length of the RDATA field (IPv4 address)
			RDATA:    []byte{8, 8, 8, 8}, // Example IPv4 address for www.google.com
		}

		rawResponseHeader := buildHeader(responseHeader)
		rawResponseBody := buildResource(responseBody)

		response := rawResponseHeader
		response = append(response, rawResponseBody...)

		conn.WriteToUDP(response, respAddr)

		fmt.Println("Sent response for ID:", header.ID)
	}
}
