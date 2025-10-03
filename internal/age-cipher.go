package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"filippo.io/age"
)

type AgeCipher struct {
	identity  age.Identity
	recipient age.Recipient
}

func NewAgeCipherWithPassphrase(passphrase string) (*AgeCipher, error) {
	recipient, err := age.NewScryptRecipient(passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to create recipient: %v", err)
	}

	identity, err := age.NewScryptIdentity(passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	return &AgeCipher{
		identity:  identity,
		recipient: recipient,
	}, nil
}

func NewAgeCipherWithKeyFile(keyFile string) (*AgeCipher, error) {

	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	identities, err := age.ParseIdentities(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse identities: %w", err)
	}

	if len(identities) == 0 {
		return nil, fmt.Errorf("no identities found in key file: %s", keyFile)
	}

	identity := identities[0]

	recipients, err := age.ParseRecipients(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipients: %w", err)
	}

	if len(recipients) == 0 {
		return nil, fmt.Errorf("no recipients found in key file: %s", keyFile)
	}

	recipient := recipients[0]

	return &AgeCipher{
		identity:  identity,
		recipient: recipient,
	}, nil
}

func (a *AgeCipher) Encrypt(plaintext string) (string, error) {
	out := &bytes.Buffer{}

	w, err := age.Encrypt(out, a.recipient)

	if err != nil {
		return "", fmt.Errorf(
			"Failed to create encrypted item: %v", err)
	}
	if _, err := io.WriteString(w, plaintext); err != nil {
		return "", fmt.Errorf(
			"Failed to write encrypted item: %v", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf(
			"Failed to close encrypted item: %v", err)
	}

	return base64.StdEncoding.EncodeToString(out.Bytes()), nil
}

func (a *AgeCipher) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	f := bytes.NewReader(data)

	r, err := age.Decrypt(f, a.identity)
	if err != nil {
		return "", fmt.Errorf(
			"Failed to read encrypted item: %v", err)
	}
	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return "", fmt.Errorf(
			"Failed to read encrypted item: %v", err)
	}

	return string(out.Bytes()), nil
}
