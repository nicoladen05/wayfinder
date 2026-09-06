package main

import (
	"encoding/binary"
	"fmt"
)

type DNSResource struct {
	NAME     []string // The domain name represented as a sequence of labels (i.e. [www google com]).
	TYPE     uint16   // The type of the resource record (16 bit)
	CLASS    uint16   // The class of the data record (16 bit)
	TTL      uint32   // The time to live of the resource record (32 bit)
	RDLENGTH uint16   // The length of the RDATA field (16 bit)
	RDATA    []byte   // The data of the resource record (variable length)
}

func (r DNSResource) string() string {
	return fmt.Sprintf(`
NAME: %s, 
TYPE: %d,
CLASS: %d,
TTL: %d,
RDLENGTH: %d,
RDATA: %v`,
		r.NAME,
		r.TYPE,
		r.CLASS,
		r.TTL,
		r.RDLENGTH,
		r.RDATA,
	)
}

func parseResource(resource []byte) DNSResource {
	name, nameLength := parseQName(resource)

	return DNSResource{
		NAME:     name,
		TYPE:     binary.BigEndian.Uint16(resource[nameLength : nameLength+2]),
		CLASS:    binary.BigEndian.Uint16(resource[nameLength+2 : nameLength+4]),
		TTL:      binary.BigEndian.Uint32(resource[nameLength+4 : nameLength+8]),
		RDLENGTH: binary.BigEndian.Uint16(resource[nameLength+8 : nameLength+10]),
		RDATA:    resource[nameLength+10:],
	}
}
