package internal

import (
	"testing"

	"filippo.io/age"
)

func TestEncryptedStore_SetGet(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}

	cipher, err := NewAgeCipher(identity.String())
	if err != nil {
		t.Fatal(err)
	}

	baseStore := NewInMemoryStore()
	encryptedStore := NewEncryptedStore(baseStore, cipher)

	err = encryptedStore.Set("key", "value")
	if err != nil {
		t.Fatal(err)
	}

	encryptedValue, err := baseStore.Get("key")
	if err != nil {
		t.Fatal(err)
	}
	if encryptedValue == "value" {
		t.Error("Base store should contain encrypted data")
	}

	decryptedValue, err := encryptedStore.Get("key")
	if err != nil {
		t.Fatal(err)
	}
	if decryptedValue != "value" {
		t.Errorf("Got %s; want %s", decryptedValue, "value")
	}
}

func TestEncryptedStore_GetAllKeys(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}

	cipher, err := NewAgeCipher(identity.String())
	if err != nil {
		t.Fatal(err)
	}

	baseStore := NewInMemoryStore()
	encryptedStore := NewEncryptedStore(baseStore, cipher)

	err = encryptedStore.Set("key1", "value1")
	if err != nil {
		t.Fatal(err)
	}
	err = encryptedStore.Set("key2", "value2")
	if err != nil {
		t.Fatal(err)
	}

	keys, err := encryptedStore.GetAllKeys()
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != 2 {
		t.Errorf("Got %d keys; want 2", len(keys))
	}

	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}
	if !keyMap["key1"] || !keyMap["key2"] {
		t.Errorf("Keys not found: %v", keys)
	}
}

func TestEncryptedStore_Delete(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}

	cipher, err := NewAgeCipher(identity.String())
	if err != nil {
		t.Fatal(err)
	}

	baseStore := NewInMemoryStore()
	encryptedStore := NewEncryptedStore(baseStore, cipher)

	err = encryptedStore.Set("key", "value")
	if err != nil {
		t.Fatal(err)
	}

	err = encryptedStore.Delete("key")
	if err != nil {
		t.Fatal(err)
	}

	_, err = encryptedStore.Get("key")
	if err == nil {
		t.Error("Expected error after delete")
	}
}
