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

func parseQuestion(question []byte) (DNSQuestion, int, error) {
	qname, qnameLength := parseDomainName(question)

	if qnameLength+4 != len(question) {
		return DNSQuestion{}, 0, errors.New("invalid question length")
	}

	return DNSQuestion{
		QNAME:  qname,
		QTYPE:  binary.BigEndian.Uint16(question[qnameLength : qnameLength+2]),
		QCLASS: binary.BigEndian.Uint16(question[qnameLength+2 : qnameLength+4]),
	}, qnameLength + 4, nil
}
