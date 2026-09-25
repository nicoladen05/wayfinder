package dnsmsg

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Question struct {
	QNAME  Name   // The domain name represented as a sequence of labels (i.e. [www google com]).
	QTYPE  uint16 // The type of the query (16 bit)
	QCLASS uint16 // The class of the query (16 bit)
}

var qType = map[uint16]string{
	0:  "Reserved",
	1:  "A",
	2:  "NS",
	5:  "CNAME",
	16: "TXT",
	28: "AAAA",
}

// Return a string representation of the question
func (q Question) String() string {
	return fmt.Sprintf(`
QNAME: %s, 
QTYPE: %d (%s),
QCLASS: %d`,
		q.QNAME,
		q.QTYPE,
		qType[q.QTYPE],
		q.QCLASS)
}

// Parse a question from a byte slice
func parseQuestion(question []byte, offset int) (Question, int, error) {
	qname, qnameLength, err := parseName(question, offset)
	if err != nil {
		return Question{}, 0, err
	}

	fieldOffset := offset + qnameLength
	if fieldOffset+4 > len(question) {
		return Question{}, 0, errors.New("invalid question length")
	}

	return Question{
		QNAME:  qname,
		QTYPE:  binary.BigEndian.Uint16(question[fieldOffset : fieldOffset+2]),
		QCLASS: binary.BigEndian.Uint16(question[fieldOffset+2 : fieldOffset+4]),
	}, qnameLength + 4, nil
}

func (q Question) Bytes() (bytes []byte) {
	name := q.QNAME.Bytes()
	nameLength := len(name)

	bytes = append(bytes, name...)

	bytes = append(bytes, make([]byte, 4)...) // Reserve 4 bytes for QTYPE and QCLASS

	binary.BigEndian.PutUint16(bytes[nameLength:nameLength+2], q.QTYPE)
	binary.BigEndian.PutUint16(bytes[nameLength+2:nameLength+4], q.QCLASS)

	return bytes
}
