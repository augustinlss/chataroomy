package pkg

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomToken(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
