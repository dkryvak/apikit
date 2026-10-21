package crypto

import (
	"crypto/sha512"
	"encoding/hex"
)

// SignaturePrefix is the scheme used in the Alea Digest header.
const SignaturePrefix = "SHA-512="

// Digest builds the Alea request signature header value:
//
//	"SHA-512=" + hex(sha512(payload + secretKey))
//
// For a grant the payload is the exact JSON body sent; for a cancel it is the
// casinoBonusId.
func Digest(secretKey, payload string) string {
	sum := sha512.Sum512([]byte(payload + secretKey))
	return SignaturePrefix + hex.EncodeToString(sum[:])
}
