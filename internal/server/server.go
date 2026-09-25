package server

import (
	"fmt"
	"wayfinder/internal/dnsmsg"
)

type Resolver interface {
	Resolve(q dnsmsg.Question) ([]dnsmsg.Resource, error)
}

type Server struct {
	addr     string
	resolver Resolver
}

func New(addr string, resolver Resolver) *Server {
	return &Server{
		addr:     addr,
		resolver: resolver,
	}
}

func (s Server) Listen() {
	question := dnsmsg.Question{
		QNAME:  dnsmsg.Name{"www", "google", "com"},
		QTYPE:  1,
		QCLASS: 1,
	}

	resp, err := s.resolver.Resolve(question)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	for _, r := range resp {
		fmt.Println(r.String())
	}
}
