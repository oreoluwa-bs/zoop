package internal

import (
	"fmt"
	"testing"
)

func TestStoreManager_New(t *testing.T) {
	// Create a mock store for testing
	mockStore := &MockStore{
		data: make(map[string]string),
	}

	manager := NewStoreManager(mockStore)

	if manager.Store != mockStore {
		t.Error("StoreManager should hold the provided store")
	}
}

type MockStore struct {
	data map[string]string
}

func (m *MockStore) Get(key string) (string, error) {
	val, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("Item with key (%s) not found", key)
	}
	return val, nil
}

func (m *MockStore) Set(key, value string) error {
	m.data[key] = value
	return nil
}

func (m *MockStore) Delete(key string) error {
	delete(m.data, key)
	return nil
}
