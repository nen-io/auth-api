package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// ErrEntropy is returned if the system CSPRNG fails (extremely rare).
var ErrEntropy = fmt.Errorf("pkce: cannot read random bytes")

// Pair returns (codeVerifier, codeChallengeS256).
func PCKE_Pair() (string, string, error) {
	// 1 –  Generate 32 cryptographically-secure random bytes.
	const bytesLen = 32 // 256-bit
	buf := make([]byte, bytesLen)
	if n, err := rand.Read(buf); err != nil || n != bytesLen {
		return "", "", ErrEntropy
	}

	// 2 –  Base64URL-encode *without* padding.
	//     RawURLEncoding removes "=" so browsers & IdPs treat it as opaque.
	verifier := base64.RawURLEncoding.EncodeToString(buf)

	// 3 –  SHA-256 hash the ASCII verifier string,
	//     then Base64URL-encode that digest -> the S256 challenge.
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])

	return verifier, challenge, nil
}
