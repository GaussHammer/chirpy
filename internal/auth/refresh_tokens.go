package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	array := make([]byte, 32)
	_, err := rand.Read(array)
	if err != nil {
		return "", err
	}
	finalString := hex.EncodeToString([]byte(array))
	return finalString, nil
}
