package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type DNSQuestion struct {
	QNAME  []string // The domain name represented as a sequence of labels (i.e. [www google com]).
	QTYPE  uint16   // The type of the query (16 bit)
	QCLASS uint16   // The class of the query (16 bit)
}

var qType = map[uint16]string{
	0:  "Reserved",
	1:  "A",
	2:  "NS",
	5:  "CNAME",
	16: "TXT",
	28: "AAAA",
}

func (q DNSQuestion) string() string {
	return fmt.Sprintf(`
QNAME: %s, 
QTYPE: %d (%s),
QCLASS: %d`,
		q.QNAME,
		q.QTYPE,
		qType[q.QTYPE],
		q.QCLASS)
}

func parseQuestion(question []byte, offset int) (DNSQuestion, int, error) {
	qname, qnameLength, err := parseDomainName(question, offset)
	if err != nil {
		return DNSQuestion{}, 0, err
	}

	fieldOffset := offset + qnameLength
	if fieldOffset+4 > len(question) {
		return DNSQuestion{}, 0, errors.New("invalid question length")
	}

	return DNSQuestion{
		QNAME:  qname,
		QTYPE:  binary.BigEndian.Uint16(question[fieldOffset : fieldOffset+2]),
		QCLASS: binary.BigEndian.Uint16(question[fieldOffset+2 : fieldOffset+4]),
	}, qnameLength + 4, nil
}

func buildQuestion(question DNSQuestion) (bytes []byte, error error) {
	name := buildDomainName(question.QNAME)
	nameLength := len(name)

	bytes = append(bytes, name...)

	bytes = append(bytes, make([]byte, 4)...) // Reserve 4 bytes for QTYPE and QCLASS

	binary.BigEndian.PutUint16(bytes[nameLength:nameLength+2], question.QTYPE)
	binary.BigEndian.PutUint16(bytes[nameLength+2:nameLength+4], question.QCLASS)

	return bytes, nil
}
