package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func RandString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64 .RawURLEncoding.EncodeToString(b)
}
