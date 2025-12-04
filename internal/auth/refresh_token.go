package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)

	encodedKey := hex.EncodeToString(key)
	return encodedKey, nil
}
