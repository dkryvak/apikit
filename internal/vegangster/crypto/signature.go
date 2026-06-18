package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func CreateSignature(payload any, key *rsa.PrivateKey) (string, error) {
	if key == nil {
		return "", fmt.Errorf("privateKey is nil")
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload to JSON: %w", err)
	}

	sum := sha256.Sum256(payloadJSON)

	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sum[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign payload: %w", err)
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}

func VerifySignature(payload any, signatureB64 string, key *rsa.PublicKey) error {
	if key == nil {
		return fmt.Errorf("publicKey is nil")
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload to JSON: %w", err)
	}

	sig, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return fmt.Errorf("failed to decode base64 signature: %w", err)
	}

	sum := sha256.Sum256(payloadJSON)

	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, sum[:], sig); err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}
	return nil
}
