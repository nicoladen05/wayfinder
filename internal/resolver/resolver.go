package resolver

import (
	"math/rand"
	"net"
	"wayfinder/internal/dnsmsg"
)

type Resolver struct {
	roots []net.IP
}

func New(roots []net.IP) *Resolver {
	return &Resolver{
		roots: roots,
	}
}

var RootServers = []net.IP{
	net.IPv4(198, 41, 0, 4),    // a.root-servers.net
	net.IPv4(170, 247, 170, 2), // b.root-servers.net
}

// Recursively resolve a DNS Query
func (r Resolver) Resolve(question dnsmsg.Question) ([]dnsmsg.Resource, error) {
	// Send the first query to one of the root servers

	header := dnsmsg.Header{
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
	}

	message := dnsmsg.Message{
		Header:    header,
		Questions: []dnsmsg.Question{question},
	}

	messageBytes, err := message.Pack()
	if err != nil {
		return nil, err
	}

	// Build the udp connection
	conn, err := net.DialUDP(
		"udp",
		nil,
		&net.UDPAddr{
			IP:   r.roots[0], // TODO: Select from one of the 13
			Port: 53,
		})
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(messageBytes)
	if err != nil {
		return nil, err
	}

	// Read the response
	buf := make([]byte, 1024)
	respLen, _, err := conn.ReadFromUDP(buf)

	// Parse the response
	response, err := dnsmsg.Parse(buf[0:respLen])
	if err != nil {
		return nil, err
	}

	return response.Answers, nil
}
