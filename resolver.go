package main

import (
	"fmt"
	"math/rand"
	"net"
)

var ROOT_SERVERS = []net.IP{
	net.IPv4(198, 41, 0, 4),    // a.root-servers.net
	net.IPv4(170, 247, 170, 2), // b.root-servers.net
}

func resolveNameserverRecursively(domain []string, qType uint16) error {
	// Send the first query to one of the root servers

	// Build the header
	header, err := buildHeader(DNSHeader{
		ID:      uint16(rand.Uint64()),
		QR:      0,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      0,
		AD:      0,
		CD:      0,
		RCODE:   0,
		QDCOUNT: 1,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	})
	if err != nil {
		return err
	}

	// Build the question
	question, err := buildQuestion(DNSQuestion{
		QNAME:  domain,
		QTYPE:  qType,
		QCLASS: 1, // Hardcoded to IN (Internet) for now
	})
	if err != nil {
		return err
	}

	// Build the udp connection
	conn, err := net.DialUDP(
		"udp",
		nil,
		&net.UDPAddr{
			IP:   ROOT_SERVERS[0],
			Port: 53,
		})
	if err != nil {
		return err
	}

	request := header
	request = append(request, question...)

	_, err = conn.Write(request)
	if err != nil {
		return err
	}

	// Read the response
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)

	// Parse the response
	respHeader, err := parseHeader(buf[:12])
	fmt.Println("\n\nROOT SERVER\n\nHeader received: \n", respHeader.string())

	lenBytes := 12

	if respHeader.QDCOUNT > 0 {
		question, len, _ := parseQuestion(buf[lenBytes:n])
		fmt.Println("Question received: ", question.string())

		lenBytes += len
	}

	if respHeader.NSCOUNT > 0 {
		answer, len, _ := parseResource(buf[lenBytes:n])
		fmt.Println("Nameserver received: ", answer.string())

		lenBytes += len
	}

	return nil
}
