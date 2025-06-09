package utils

import (
	"crypto/rand"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func UUIDv4() (string, error) {
	uuid := make([]byte, 16)

	// Read 16 random bytes
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant (10)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 10

	// Format as UUID string
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16],
	), nil
}

func NanoID() (string, error) {
	return gonanoid.New()
}
