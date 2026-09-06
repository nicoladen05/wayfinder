package main

import (
	"encoding/binary"
	"errors"
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

func parseResource(rawResource []byte) (DNSResource, int, error) {
	name, nameLength := parseDomainName(rawResource)

	headLength := nameLength + 10 // 2 bytes for TYPE, 2 bytes for CLASS, 4 bytes for TTL, 2 bytes for RDLENGTH

	if headLength > len(rawResource) {
		return DNSResource{}, 0, errors.New("invalid resource length")
	}

	dataLength := binary.BigEndian.Uint16(rawResource[nameLength+8 : nameLength+10]) // RDLENGTH

	if headLength+int(dataLength) > len(rawResource) {
		return DNSResource{}, 0, errors.New("invalid resource length")
	}

	resource := DNSResource{
		NAME:     name,
		TYPE:     binary.BigEndian.Uint16(rawResource[nameLength : nameLength+2]),
		CLASS:    binary.BigEndian.Uint16(rawResource[nameLength+2 : nameLength+4]),
		TTL:      binary.BigEndian.Uint32(rawResource[nameLength+4 : nameLength+8]),
		RDLENGTH: dataLength,
		RDATA:    rawResource[headLength : headLength+int(dataLength)],
	}

	return resource, nameLength + 10 + int(resource.RDLENGTH), nil
}

func buildResource(resource DNSResource) (bytes []byte, error error) {
	nameBytes := buildDomainName(resource.NAME)
	nameLength := len(nameBytes)

	bytes = append(bytes, nameBytes...)
	bytes = append(bytes, make([]byte, 10)...) // Reserve space for TYPE, CLASS, TTL, RDLENGTH

	binary.BigEndian.PutUint16(bytes[nameLength:nameLength+2], resource.TYPE)
	binary.BigEndian.PutUint16(bytes[nameLength+2:nameLength+4], resource.CLASS)
	binary.BigEndian.PutUint32(bytes[nameLength+4:nameLength+8], resource.TTL)
	binary.BigEndian.PutUint16(bytes[nameLength+8:nameLength+10], resource.RDLENGTH)

	bytes = append(bytes, resource.RDATA...)

	return bytes, nil
}
