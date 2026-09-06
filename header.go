package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type DNSHeader struct {
	ID      uint16 // Identification of the DNS query (16 bit)
	QR      uint8  // Query/Response Flag (1 bit) (0 = query, 1 = response)
	OPCODE  uint8  // Kind of the query (4 bit) (1 = Standard Query, 2 = Inverse Query, 3 = Server status request)
	AA      uint8  // Authoritative Answer (1 bit) (1 = Authoritative, 0 = Not authoritative)
	TC      uint8  // Truncation (1 bit) (1 = Message truncated, 0 = Message not truncated)
	RD      uint8  // Recursion Desired (1 bit) (1 = Recursion desired, 0 = Recursion not desired)
	RA      uint8  // Recursion Available (1 bit) (1 = Recursion available, 0 = Recursion not available)
	Z       uint8  // Reserved (1 bit) (Must be 0)
	AD      uint8  // Authentic Data (1 bit) (1 = Authentic data, 0 = Not authentic data)
	CD      uint8  // Checking Disabled (1 bit) (1 = Checking disabled, 0 = Checking enabled)
	RCODE   uint8  // Response code (4 bit) (0 = No error, 1 = Format error, 2 = Server failure, 3 = Name error, 4 = Not implemented, 5 = Refused)
	QDCOUNT uint16 // Number of questions (16 bit)
	ANCOUNT uint16 // Number of answers (16 bit)
	NSCOUNT uint16 // Number of authority records (16 bit)
	ARCOUNT uint16 // Number of additional records (16 bit)
}

func (h DNSHeader) isValid() bool {
	return !(h.QR > 1 ||
		h.OPCODE > 15 ||
		h.AA > 1 ||
		h.TC > 1 ||
		h.RD > 1 ||
		h.RA > 1 ||
		h.Z > 1 ||
		h.AD > 1 ||
		h.CD > 1 ||
		h.RCODE > 15)
}

func (h DNSHeader) string() string {
	return fmt.Sprintf(
		`ID: %d,
QR: %d,
OPCODE: %d,
AA: %d,
TC: %d,
RD: %d,
RA: %d,
Z: %d,
AD: %d,
CD: %d,
RCODE: %d,
QDCOUNT: %d,
ANCOUNT: %d,
NSCOUNT: %d,
ARCOUNT: %d`,
		h.ID,
		h.QR,
		h.OPCODE,
		h.AA,
		h.TC,
		h.RD,
		h.RA,
		h.Z,
		h.AD,
		h.CD,
		h.RCODE,
		h.QDCOUNT,
		h.ANCOUNT,
		h.NSCOUNT,
		h.ARCOUNT)
}

func parseHeader(header []byte) (DNSHeader, error) {
	if len(header) != 12 {
		return DNSHeader{}, errors.New("header must be exactly 12 bytes long")
	}

	return DNSHeader{
		ID:      binary.BigEndian.Uint16(header[0:2]),
		QR:      header[2] >> 7,
		OPCODE:  (header[2] >> 3) & 0x0F,
		AA:      (header[2] >> 2) & 0x01,
		TC:      (header[2] >> 1) & 0x01,
		RD:      header[2] & 0x01,
		RA:      (header[3] >> 7) & 0x01,
		Z:       (header[3] >> 6) & 0x01,
		AD:      (header[3] >> 5) & 0x01,
		CD:      (header[3] >> 4) & 0x01,
		RCODE:   header[3] & 0x0F,
		QDCOUNT: binary.BigEndian.Uint16(header[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(header[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(header[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(header[10:12]),
	}, nil
}

func buildHeader(header DNSHeader) ([]byte, error) {
	if !header.isValid() {
		return nil, errors.New("invalid header values")
	}

	initialHeader := make([]byte, 12)

	binary.BigEndian.PutUint16(initialHeader[0:2], header.ID)
	initialHeader[2] = (header.QR << 7) | (header.OPCODE << 3) | (header.AA << 2) | (header.TC << 1) | header.RD
	binary.BigEndian.PutUint16(initialHeader[4:6], header.QDCOUNT)
	binary.BigEndian.PutUint16(initialHeader[6:8], header.ANCOUNT)
	binary.BigEndian.PutUint16(initialHeader[8:10], header.NSCOUNT)
	binary.BigEndian.PutUint16(initialHeader[10:12], header.ARCOUNT)

	return initialHeader, nil
}
