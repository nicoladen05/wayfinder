package dnsmsg

import (
	"fmt"
	"strings"
)

type Message struct {
	Header      Header
	Questions   []Question
	Answers     []Resource
	Authorities []Resource
	Additionals []Resource
}

type parser[T Resource | Question] func([]byte, int) (T, int, error)

// Return a string representation of a question or resource section
func sectionString[T fmt.Stringer](items []T) string {
	if len(items) == 0 {
		return "\n(none)"
	}

	var sb strings.Builder
	for i, item := range items {
		fmt.Fprintf(&sb, "\n[%d]%s", i, item)
	}

	return sb.String()
}

// Return a string representation of the message
func (m Message) String() string {
	return fmt.Sprintf(`HEADER:
%s

QUESTIONS (%d):%s

ANSWERS (%d):%s

AUTHORITIES (%d):%s

ADDITIONALS (%d):%s`,
		m.Header,
		len(m.Questions), sectionString(m.Questions),
		len(m.Answers), sectionString(m.Answers),
		len(m.Authorities), sectionString(m.Authorities),
		len(m.Additionals), sectionString(m.Additionals),
	)
}

// Parse a question or resource section
func parseSection[T Resource | Question](parser parser[T], b []byte, sectionCount int, offset *int) ([]T, error) {
	items := make([]T, sectionCount)

	for i := range sectionCount {
		item, itemLen, err := parser(b, *offset)
		if err != nil {
			return []T{}, err
		}

		items[i] = item
		*offset += itemLen
	}

	return items, nil
}

// Parse a DNS message from a byte slice.
func Parse(b []byte) (Message, error) {
	header, err := parseHeader(b)
	if err != nil {
		return Message{}, err
	}

	offset := 12 // Header is always 12 bytes long
	questions, err := parseSection(parseQuestion, b, int(header.QDCOUNT), &offset)
	if err != nil {
		return Message{}, err
	}

	answers, err := parseSection(parseResource, b, int(header.ANCOUNT), &offset)
	if err != nil {
		return Message{}, err
	}

	authorities, err := parseSection(parseResource, b, int(header.NSCOUNT), &offset)
	if err != nil {
		return Message{}, err
	}

	additionals, err := parseSection(parseResource, b, int(header.ARCOUNT), &offset)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Header:      header,
		Questions:   questions,
		Answers:     answers,
		Authorities: authorities,
		Additionals: additionals,
	}, nil

}

// Pack a DNS message into a byte slice.
func (m Message) Pack() (bytes []byte, err error) {
	headerBytes, err := m.Header.Bytes()
	if err != nil {
		return nil, err
	}
	bytes = append(bytes, headerBytes...)

	for _, question := range m.Questions {
		bytes = append(bytes, question.Bytes()...)
	}

	for _, answer := range m.Answers {
		bytes = append(bytes, answer.Bytes()...)
	}

	for _, authority := range m.Authorities {
		bytes = append(bytes, authority.Bytes()...)
	}

	for _, additional := range m.Additionals {
		bytes = append(bytes, additional.Bytes()...)
	}

	return bytes, nil
}
