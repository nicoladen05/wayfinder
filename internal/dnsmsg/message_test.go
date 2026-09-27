package dnsmsg

import (
	"bytes"
	"reflect"
	"testing"
)

// A standard query header with recursion desired and no sections.
var exampleEmptyMessageBytes = []byte{
	0x12, 0x34, // ID: 0x1234
	0x01, 0x00, // Flags: QR=0, OPCODE=0, AA=0, TC=0, RD=1, RA=0, Z=0, AD=0, CD=0, RCODE=0
	0x00, 0x00, // QDCOUNT: 0
	0x00, 0x00, // ANCOUNT: 0
	0x00, 0x00, // NSCOUNT: 0
	0x00, 0x00, // ARCOUNT: 0
}

// The Message that exampleEmptyMessage parses to.
var exampleEmptyMessage = Message{
	Header: Header{
		ID:      0x1234,
		QR:      0,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      1,
		RA:      0,
		Z:       0,
		AD:      0,
		CD:      0,
		RCODE:   0,
		QDCOUNT: 0,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	},
	Questions:   []Question{},
	Answers:     []Resource{},
	Authorities: []Resource{},
	Additionals: []Resource{},
}

func TestParseEmptyMessage(t *testing.T) {
	msg, _ := Parse(exampleEmptyMessageBytes)

	if !reflect.DeepEqual(msg, exampleEmptyMessage) {
		t.Errorf("Returned message did not match.\nExpected:\n%s\nGot:\n%s", exampleEmptyMessage, msg)
	}
}

func TestPackEmptyMessage(t *testing.T) {
	b, err := exampleEmptyMessage.Pack()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !bytes.Equal(b, exampleEmptyMessageBytes) {
		t.Errorf("Returned bytes did not match.\nExpected:\n% x\nGot:\n% x", exampleEmptyMessageBytes, b)
	}
}

// A standard query for the A record of www.example.com with recursion desired.
var exampleQuestionMessageBytes = []byte{
	0x12, 0x34, // ID: 0x1234
	0x01, 0x00, // Flags: QR=0, OPCODE=0, AA=0, TC=0, RD=1, RA=0, Z=0, AD=0, CD=0, RCODE=0
	0x00, 0x01, // QDCOUNT: 1
	0x00, 0x00, // ANCOUNT: 0
	0x00, 0x00, // NSCOUNT: 0
	0x00, 0x00, // ARCOUNT: 0

	// Question (offset 12)
	0x03, 'w', 'w', 'w', // Label: www
	0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e', // Label: example
	0x03, 'c', 'o', 'm', // Label: com
	0x00,       // End of name
	0x00, 0x01, // QTYPE: 1 (A)
	0x00, 0x01, // QCLASS: 1 (IN)
}

// The Message that exampleQuestionMessageBytes parses to.
var exampleQuestionMessage = Message{
	Header: Header{
		ID:      0x1234,
		QR:      0,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      1,
		RA:      0,
		Z:       0,
		AD:      0,
		CD:      0,
		RCODE:   0,
		QDCOUNT: 1,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	},
	Questions: []Question{
		{QNAME: Name{"www", "example", "com"}, QTYPE: 1, QCLASS: 1},
	},
	Answers:     []Resource{},
	Authorities: []Resource{},
	Additionals: []Resource{},
}

func TestParseQuestionMessage(t *testing.T) {
	msg, _ := Parse(exampleQuestionMessageBytes)

	if !reflect.DeepEqual(msg, exampleQuestionMessage) {
		t.Errorf("Returned message did not match.\nExpected:\n%s\nGot:\n%s", exampleQuestionMessage, msg)
	}
}

func TestPackQuestionMessage(t *testing.T) {
	b, err := exampleQuestionMessage.Pack()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !bytes.Equal(b, exampleQuestionMessageBytes) {
		t.Errorf("Returned bytes did not match.\nExpected:\n% x\nGot:\n% x", exampleQuestionMessageBytes, b)
	}
}

// A recursive resolver's response to the query above, answering with a single A record.
// The answer's name is compressed and points back to the question name.
var exampleAnswerMessageBytes = []byte{
	0x12, 0x34, // ID: 0x1234
	0x81, 0x80, // Flags: QR=1, OPCODE=0, AA=0, TC=0, RD=1, RA=1, Z=0, AD=0, CD=0, RCODE=0
	0x00, 0x01, // QDCOUNT: 1
	0x00, 0x01, // ANCOUNT: 1
	0x00, 0x00, // NSCOUNT: 0
	0x00, 0x00, // ARCOUNT: 0

	// Question (offset 12)
	0x03, 'w', 'w', 'w', // Label: www
	0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e', // Label: example
	0x03, 'c', 'o', 'm', // Label: com
	0x00,       // End of name
	0x00, 0x01, // QTYPE: 1 (A)
	0x00, 0x01, // QCLASS: 1 (IN)

	// Answer (offset 33)
	0xC0, 0x0C, // NAME: pointer to offset 12 (www.example.com)
	0x00, 0x01, // TYPE: 1 (A)
	0x00, 0x01, // CLASS: 1 (IN)
	0x00, 0x00, 0x0E, 0x10, // TTL: 3600
	0x00, 0x04, // RDLENGTH: 4
	0x5D, 0xB8, 0xD8, 0x22, // RDATA: 93.184.216.34
}

// The Message that exampleAnswerMessageBytes parses to.
var exampleAnswerMessage = Message{
	Header: Header{
		ID:      0x1234,
		QR:      1,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      1,
		RA:      1,
		Z:       0,
		AD:      0,
		CD:      0,
		RCODE:   0,
		QDCOUNT: 1,
		ANCOUNT: 1,
		NSCOUNT: 0,
		ARCOUNT: 0,
	},
	Questions: []Question{
		{QNAME: Name{"www", "example", "com"}, QTYPE: 1, QCLASS: 1},
	},
	Answers: []Resource{
		{
			NAME:     Name{"www", "example", "com"},
			TYPE:     1,
			CLASS:    1,
			TTL:      3600,
			RDLENGTH: 4,
			RDATA:    []byte{0x5D, 0xB8, 0xD8, 0x22},
		},
	},
	Authorities: []Resource{},
	Additionals: []Resource{},
}

func TestParseAnswerMessage(t *testing.T) {
	msg, _ := Parse(exampleAnswerMessageBytes)

	if !reflect.DeepEqual(msg, exampleAnswerMessage) {
		t.Errorf("Returned message did not match.\nExpected:\n%s\nGot:\n%s", exampleAnswerMessage, msg)
	}
}

func TestPackAnswerMessage(t *testing.T) {
	b, err := exampleAnswerMessage.Pack()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !bytes.Equal(b, exampleAnswerMessageBytes) {
		t.Errorf("Returned bytes did not match.\nExpected:\n% x\nGot:\n% x", exampleAnswerMessageBytes, b)
	}
}

// A TLD server's referral for example.com to its nameservers, as seen during iterative resolution.
// The authority section holds the NS records and the additional section holds the glue A records.
// The glue record names point into the RDATA of the NS records.
var exampleReferralMessageBytes = []byte{
	0xAB, 0xCD, // ID: 0xABCD
	0x80, 0x00, // Flags: QR=1, OPCODE=0, AA=0, TC=0, RD=0, RA=0, Z=0, AD=0, CD=0, RCODE=0
	0x00, 0x01, // QDCOUNT: 1
	0x00, 0x00, // ANCOUNT: 0
	0x00, 0x02, // NSCOUNT: 2
	0x00, 0x02, // ARCOUNT: 2

	// Question (offset 12)
	0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e', // Label: example
	0x03, 'c', 'o', 'm', // Label: com
	0x00,       // End of name
	0x00, 0x01, // QTYPE: 1 (A)
	0x00, 0x01, // QCLASS: 1 (IN)

	// Authority 1 (offset 29)
	0xC0, 0x0C, // NAME: pointer to offset 12 (example.com)
	0x00, 0x02, // TYPE: 2 (NS)
	0x00, 0x01, // CLASS: 1 (IN)
	0x00, 0x02, 0xA3, 0x00, // TTL: 172800
	0x00, 0x06, // RDLENGTH: 6
	0x03, 'n', 's', '1', 0xC0, 0x0C, // RDATA (offset 41): ns1 + pointer to offset 12 (ns1.example.com)

	// Authority 2 (offset 47)
	0xC0, 0x0C, // NAME: pointer to offset 12 (example.com)
	0x00, 0x02, // TYPE: 2 (NS)
	0x00, 0x01, // CLASS: 1 (IN)
	0x00, 0x02, 0xA3, 0x00, // TTL: 172800
	0x00, 0x06, // RDLENGTH: 6
	0x03, 'n', 's', '2', 0xC0, 0x0C, // RDATA (offset 59): ns2 + pointer to offset 12 (ns2.example.com)

	// Additional 1 (offset 65)
	0xC0, 0x29, // NAME: pointer to offset 41 (ns1.example.com)
	0x00, 0x01, // TYPE: 1 (A)
	0x00, 0x01, // CLASS: 1 (IN)
	0x00, 0x02, 0xA3, 0x00, // TTL: 172800
	0x00, 0x04, // RDLENGTH: 4
	0xC0, 0x00, 0x02, 0x01, // RDATA: 192.0.2.1

	// Additional 2 (offset 81)
	0xC0, 0x3B, // NAME: pointer to offset 59 (ns2.example.com)
	0x00, 0x01, // TYPE: 1 (A)
	0x00, 0x01, // CLASS: 1 (IN)
	0x00, 0x02, 0xA3, 0x00, // TTL: 172800
	0x00, 0x04, // RDLENGTH: 4
	0xC0, 0x00, 0x02, 0x02, // RDATA: 192.0.2.2
}

// The Message that exampleReferralMessageBytes parses to.
var exampleReferralMessage = Message{
	Header: Header{
		ID:      0xABCD,
		QR:      1,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      0,
		RA:      0,
		Z:       0,
		AD:      0,
		CD:      0,
		RCODE:   0,
		QDCOUNT: 1,
		ANCOUNT: 0,
		NSCOUNT: 2,
		ARCOUNT: 2,
	},
	Questions: []Question{
		{QNAME: Name{"example", "com"}, QTYPE: 1, QCLASS: 1},
	},
	Answers: []Resource{},
	Authorities: []Resource{
		{
			NAME:     Name{"example", "com"},
			TYPE:     2,
			CLASS:    1,
			TTL:      172800,
			RDLENGTH: 6,
			RDATA:    []byte{0x03, 'n', 's', '1', 0xC0, 0x0C},
		},
		{
			NAME:     Name{"example", "com"},
			TYPE:     2,
			CLASS:    1,
			TTL:      172800,
			RDLENGTH: 6,
			RDATA:    []byte{0x03, 'n', 's', '2', 0xC0, 0x0C},
		},
	},
	Additionals: []Resource{
		{
			NAME:     Name{"ns1", "example", "com"},
			TYPE:     1,
			CLASS:    1,
			TTL:      172800,
			RDLENGTH: 4,
			RDATA:    []byte{0xC0, 0x00, 0x02, 0x01},
		},
		{
			NAME:     Name{"ns2", "example", "com"},
			TYPE:     1,
			CLASS:    1,
			TTL:      172800,
			RDLENGTH: 4,
			RDATA:    []byte{0xC0, 0x00, 0x02, 0x02},
		},
	},
}

func TestParseReferralMessage(t *testing.T) {
	msg, _ := Parse(exampleReferralMessageBytes)

	if !reflect.DeepEqual(msg, exampleReferralMessage) {
		t.Errorf("Returned message did not match.\nExpected:\n%s\nGot:\n%s", exampleReferralMessage, msg)
	}
}

func TestPackReferralMessage(t *testing.T) {
	b, err := exampleReferralMessage.Pack()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !bytes.Equal(b, exampleReferralMessageBytes) {
		t.Errorf("Returned bytes did not match.\nExpected:\n% x\nGot:\n% x", exampleReferralMessageBytes, b)
	}
}

// A message that is cut off in the middle of the header.
var exampleShortHeaderMessageBytes = []byte{
	0x12, 0x34, // ID: 0x1234
	0x01, 0x00, // Flags: QR=0, OPCODE=0, AA=0, TC=0, RD=1, RA=0, Z=0, AD=0, CD=0, RCODE=0
	0x00, 0x01, // QDCOUNT: 1
	// ANCOUNT, NSCOUNT and ARCOUNT are missing
}

func TestParseShortHeader(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{"empty", []byte{}},
		{"partial header", exampleShortHeaderMessageBytes},
		{"one byte short", exampleEmptyMessageBytes[:11]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := Parse(tt.bytes)
			if err == nil {
				t.Errorf("Expected an error for a message shorter than 12 bytes, got message:\n%s", msg)
			}
		})
	}
}
