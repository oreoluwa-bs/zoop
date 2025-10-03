package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"filippo.io/age"
)

type Cipher interface {
	Decrypt(ciphertext string) (string, error)
	Encrypt(plaintext string) (string, error)
}

type AgeCipher struct {
	identity *age.X25519Identity
}

func NewAgeCipher(identityStr string) (*AgeCipher, error) {
	identity, err := age.ParseX25519Identity(identityStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse identity: %w", err)
	}
	return &AgeCipher{identity: identity}, nil
}

func (c *AgeCipher) Encrypt(plaintext string) (string, error) {
	buf := &bytes.Buffer{}
	w, err := age.Encrypt(buf, c.identity.Recipient())
	if err != nil {
		return "", fmt.Errorf("failed to encrypt: %w", err)
	}
	_, err = io.WriteString(w, plaintext)
	if err != nil {
		return "", fmt.Errorf("failed to write plaintext: %w", err)
	}
	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (c *AgeCipher) Decrypt(ciphertext string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	buf := bytes.NewBuffer(decoded)
	r, err := age.Decrypt(buf, c.identity)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}
	decrypted, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to read decrypted data: %w", err)
	}
	return string(decrypted), nil
}
