package tokens

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecureToken(length int) (string, error) {
	bytesNeeded := length / 2
	b := make([]byte, bytesNeeded)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
