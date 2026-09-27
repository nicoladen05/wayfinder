package dnsmsg

import (
	"reflect"
	"testing"
)

func TestParseNameCompressed(t *testing.T) {
	bytes := []byte{
		// Offset 0
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e', // Label: example
		0x03, 'c', 'o', 'm', // Label: com
		0x00, // End of name

		// Offset 13
		0x03, 'w', 'w', 'w', // Label: www
		0xC0, 0x00, // Pointer to offset 0 (example.com)
	}

	name, length, err := parseName(bytes, 13)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := Name{"www", "example", "com"}
	if !reflect.DeepEqual(name, expected) {
		t.Errorf("Returned name did not match.\nExpected: %s\nGot: %s", expected, name)
	}

	// The label and the pointer take up 6 bytes; the pointed-to suffix does not count.
	if length != 6 {
		t.Errorf("Expected length 6, got %d", length)
	}
}

func TestParseNameInvalidPointer(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{"self pointer", []byte{
			0xC0, 0x00, // Pointer to offset 0 (itself)
		}},
		{"pointer loop", []byte{
			0x03, 'w', 'w', 'w', // Label: www
			0xC0, 0x06, // Pointer to offset 6
			0xC0, 0x00, // Pointer back to offset 0
		}},
		{"pointer past end", []byte{
			0xC0, 0xFF, // Pointer to offset 255 (name is only 2 bytes long)
		}},
		{"truncated pointer", []byte{
			0xC0, // First byte of a pointer, name ends here
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, _, err := parseName(tt.bytes, 0)
			if err == nil {
				t.Errorf("Expected an error for an invalid pointer, got name: %s", name)
			}
		})
	}
}
