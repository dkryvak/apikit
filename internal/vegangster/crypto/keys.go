package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

func ParsePrivateKey(privateKeyStr string) (*rsa.PrivateKey, error) {
	s := strings.TrimSpace(privateKeyStr)

	if strings.Contains(s, "BEGIN") {
		s = strings.ReplaceAll(s, "-----BEGIN PRIVATE KEY-----", "")
		s = strings.ReplaceAll(s, "-----END PRIVATE KEY-----", "")
		s = strings.ReplaceAll(s, "-----BEGIN RSA PRIVATE KEY-----", "")
		s = strings.ReplaceAll(s, "-----END RSA PRIVATE KEY-----", "")
		s = strings.Join(strings.Fields(s), "")

		der, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, fmt.Errorf("invalid RSA private key (PEM/base64): %w", err)
		}
		return parsePKCS8PrivateDER(der)
	}

	der, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("invalid RSA private key (PKCS#8 base64): %w", err)
	}
	return parsePKCS8PrivateDER(der)
}

func parsePKCS8PrivateDER(der []byte) (*rsa.PrivateKey, error) {
	parsed, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("invalid RSA private key (PKCS#8 DER): %w", err)
	}

	privateKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid RSA private key: expected *rsa.PrivateKey, got %T", parsed)
	}
	return privateKey, nil
}

func ParsePublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	s := strings.TrimSpace(publicKeyStr)

	if strings.Contains(s, "BEGIN") {
		if block, _ := pem.Decode([]byte(s)); block != nil {
			return parsePKIXPublicDER(block.Bytes)
		}

		s = strings.ReplaceAll(s, "-----BEGIN PUBLIC KEY-----", "")
		s = strings.ReplaceAll(s, "-----END PUBLIC KEY-----", "")
		s = strings.ReplaceAll(s, "-----BEGIN RSA PUBLIC KEY-----", "")
		s = strings.ReplaceAll(s, "-----END RSA PUBLIC KEY-----", "")
		s = strings.Join(strings.Fields(s), "")

		der, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, fmt.Errorf("invalid RSA public key (PEM/base64): %w", err)
		}
		return parsePKIXPublicDER(der)
	}

	der, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("invalid RSA public key (X.509 base64): %w", err)
	}
	return parsePKIXPublicDER(der)
}

func parsePKIXPublicDER(der []byte) (*rsa.PublicKey, error) {
	parsed, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("invalid RSA public key (X.509/PKIX DER): %w", err)
	}

	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid RSA public key: expected *rsa.PublicKey, got %T", parsed)
	}
	return publicKey, nil
}
