package main

import (
	"net"
	"wayfinder/internal/resolver"
	"wayfinder/internal/server"
)

var rootServers = []net.IP{
	net.IPv4(198, 41, 0, 4),    // a.root-servers.net
	net.IPv4(170, 247, 170, 2), // b.root-servers.net
}

func main() {
	r := resolver.New(rootServers)
	s := server.New("127.0.0.1:1053", r)
	s.Listen()
}
