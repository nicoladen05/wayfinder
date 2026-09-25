package dnsmsg

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Resource struct {
	NAME     Name   // The domain name represented as a sequence of labels (i.e. [www google com]).
	TYPE     uint16 // The type of the resource record (16 bit)
	CLASS    uint16 // The class of the data record (16 bit)
	TTL      uint32 // The time to live of the resource record (32 bit)
	RDLENGTH uint16 // The length of the RDATA field (16 bit)
	RDATA    []byte // The data of the resource record (variable length)
}

// Return the string representation of the resource
func (r Resource) String() string {
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

// Parse bytes into a resource
//
// Parameters:
//   - bytes: the bytes to parse
//   - offset: the offset in bytes at which to start parsing
//
// Returns the parsed Resource and the length of parsed section in bytes or an error.
func parseResource(bytes []byte, offset int) (Resource, int, error) {
	name, nameLength, err := parseName(bytes, offset)
	if err != nil {
		return Resource{}, 0, err
	}

	fieldOffset := offset + nameLength
	headEnd := fieldOffset + 10 // 2 bytes for TYPE, 2 bytes for CLASS, 4 bytes for TTL, 2 bytes for RDLENGTH

	if headEnd > len(bytes) {
		return Resource{}, 0, errors.New("invalid resource length")
	}

	dataLength := binary.BigEndian.Uint16(bytes[fieldOffset+8 : fieldOffset+10])
	dataEnd := headEnd + int(dataLength)

	if dataEnd > len(bytes) {
		return Resource{}, 0, errors.New("invalid resource length")
	}

	resource := Resource{
		NAME:     name,
		TYPE:     binary.BigEndian.Uint16(bytes[fieldOffset : fieldOffset+2]),
		CLASS:    binary.BigEndian.Uint16(bytes[fieldOffset+2 : fieldOffset+4]),
		TTL:      binary.BigEndian.Uint32(bytes[fieldOffset+4 : fieldOffset+8]),
		RDLENGTH: dataLength,
		RDATA:    bytes[headEnd:dataEnd],
	}

	return resource, nameLength + 10 + int(resource.RDLENGTH), nil
}

// Return the byte representation of the resource
func (r Resource) Bytes() (bytes []byte) {
	nameBytes := r.NAME.Bytes()
	nameLength := len(nameBytes)

	bytes = append(bytes, nameBytes...)
	bytes = append(bytes, make([]byte, 10)...) // Reserve space for TYPE, CLASS, TTL, RDLENGTH

	binary.BigEndian.PutUint16(bytes[nameLength:nameLength+2], r.TYPE)
	binary.BigEndian.PutUint16(bytes[nameLength+2:nameLength+4], r.CLASS)
	binary.BigEndian.PutUint32(bytes[nameLength+4:nameLength+8], r.TTL)
	binary.BigEndian.PutUint16(bytes[nameLength+8:nameLength+10], r.RDLENGTH)

	bytes = append(bytes, r.RDATA...)

	return bytes
}
