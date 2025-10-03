package internal

import (
	"testing"

	"filippo.io/age"
)

func TestAgeCipher_EncryptDecrypt(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}

	cipher, err := NewAgeCipher(identity.String())
	if err != nil {
		t.Fatal(err)
	}

	plaintext := "hello world"

	encrypted, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	if encrypted == plaintext {
		t.Error("Encrypted text should be different from plaintext")
	}

	decrypted, err := cipher.Decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
	}

	if decrypted != plaintext {
		t.Errorf("Got %s; want %s", decrypted, plaintext)
	}
}
