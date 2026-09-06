package main

import "fmt"

func isCompressed(label byte) bool {
	// A name is compressed, if the first two bits are 1.
	return label&0b11000000 == 0b11000000
}

func readNames(name []byte, offset int, seen map[int]bool) (names []string, totalLength int, err error) {
	if offset >= len(name) {
		return nil, 0, fmt.Errorf("offset %d is out of bounds for name length %d", offset, len(name))
	}

	if seen[offset] {
		return nil, 0, fmt.Errorf("circular reference detected at offset %d", offset)
	}
	seen[offset] = true // Mark this offset as seen to prevent circular references

	for {
		pos := offset + totalLength
		if pos >= len(name) {
			return nil, 0, fmt.Errorf("truncated name")
		}

		currentByte := name[pos]

		// The first byte indicates the nextNameLength of the upcoming label
		nextNameLength := int(currentByte)

		// If the length is 0, we have reached the end of the qname section
		if nextNameLength == 0 {
			totalLength++ // Move past the zero-length byte
			break
		}

		if isCompressed(currentByte) {
			if pos+1 >= len(name) {
				return nil, 0, fmt.Errorf("pointer at offset %d is out of bounds for name length %d", pos+1, len(name))
			}
			nextByte := name[pos+1]

			offset := (int(currentByte&0b00111111))<<8 | int(nextByte) // Get the offset from the current and next byte

			// Recursively read the names from the section the pointer references
			suffix, _, err := readNames(name, int(offset), seen)
			if err != nil {
				return nil, 0, err
			}

			totalLength += 2 // Move past the two bytes of the pointer
			return append(names, suffix...), totalLength, err
		}

		totalLength++ // Move to the next byte after the length byte

		start := pos + 1
		end := start + nextNameLength
		if end > len(name) {
			return nil, 0, fmt.Errorf(
				"label length %d exceeds remaining bytes %d",
				nextNameLength, len(name)-start,
			)
		}

		labelBytes := name[totalLength+offset : totalLength+offset+nextNameLength]
		label := string(labelBytes)
		names = append(names, label)

		totalLength += nextNameLength // Move to the next label length byte
	}

	return names, totalLength, nil
}

func parseDomainName(name []byte, offset int) (names []string, totalLength int, err error) {
	return readNames(name, offset, make(map[int]bool))
}

func buildDomainName(labels []string) (bytes []byte) {
	for _, label := range labels {
		length := uint8(len(label))
		label := []byte(label)

		bytes = append(bytes, length)
		bytes = append(bytes, label...)
	}

	// Append the zero-length label to indicate the end of the domain name
	bytes = append(bytes, 0)

	return bytes
}
