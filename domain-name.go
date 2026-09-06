package main

func parseDomainName(qname []byte) (labels []string, i int) {
	for {
		// The first byte indicates then length of the upcoming label
		length := int(qname[i])
		i++
		// If the length is 0, we have reached the end of the qname section
		if length == 0 {
			break
		}

		// Read the label of the specified length
		label := string(qname[i : i+length])
		labels = append(labels, label)

		i += length
	}

	return labels, i
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
