package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func CreateRefreshToken() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)

	return hex.EncodeToString(key), nil
}
