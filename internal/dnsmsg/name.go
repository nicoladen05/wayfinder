package dnsmsg

import (
	"fmt"
	"strings"
)

// A domain name consiting of 1:many labels.
type Name []string

// Return the string representation of the label (i.e. www.google.com)
func (n Name) String() string {
	return strings.Join(n, ".")
}

func isCompressed(label byte) bool {
	// A name is compressed, if the first two bits are 1.
	return label&0b11000000 == 0b11000000
}

func readName(b []byte, offset int, seen map[int]bool) (name Name, totalLength int, err error) {
	if offset >= len(b) {
		return nil, 0, fmt.Errorf("offset %d is out of bounds for name length %d", offset, len(b))
	}

	if seen[offset] {
		return nil, 0, fmt.Errorf("circular reference detected at offset %d", offset)
	}
	seen[offset] = true // Mark this offset as seen to prevent circular references

	for {
		pos := offset + totalLength
		if pos >= len(b) {
			return nil, 0, fmt.Errorf("truncated name")
		}

		currentByte := b[pos]

		// The first byte indicates the nextNameLength of the upcoming label
		nextNameLength := int(currentByte)

		// If the length is 0, we have reached the end of the qname section
		if nextNameLength == 0 {
			totalLength++ // Move past the zero-length byte
			break
		}

		if isCompressed(currentByte) {
			if pos+1 >= len(b) {
				return nil, 0, fmt.Errorf("pointer at offset %d is out of bounds for name length %d", pos+1, len(b))
			}
			nextByte := b[pos+1]

			offset := (int(currentByte&0b00111111))<<8 | int(nextByte) // Get the offset from the current and next byte

			// Recursively read the names from the section the pointer references
			suffix, _, err := readName(b, int(offset), seen)
			if err != nil {
				return nil, 0, err
			}

			totalLength += 2 // Move past the two bytes of the pointer
			return append(name, suffix...), totalLength, err
		}

		totalLength++ // Move to the next byte after the length byte

		start := pos + 1
		end := start + nextNameLength
		if end > len(b) {
			return nil, 0, fmt.Errorf(
				"label length %d exceeds remaining bytes %d",
				nextNameLength, len(b)-start,
			)
		}

		labelBytes := b[totalLength+offset : totalLength+offset+nextNameLength]
		label := string(labelBytes)
		name = append(name, label)

		totalLength += nextNameLength // Move to the next label length byte
	}

	return name, totalLength, nil
}

// Parse a domain name from a byte slice.
func parseName(nameBytes []byte, offset int) (name Name, totalLength int, err error) {
	return readName(nameBytes, offset, make(map[int]bool))
}

// Returns the byte representation of a domain name.
func (n Name) Bytes(compress bool) (bytes []byte) {
	for _, label := range n {
		length := uint8(len(label))
		label := []byte(label)

		bytes = append(bytes, length)
		bytes = append(bytes, label...)
	}

	// Append the zero-length label to indicate the end of the domain name
	bytes = append(bytes, 0)

	return bytes
}
