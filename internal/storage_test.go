package internal

import (
	"testing"
)

func TestStoreManager_New(t *testing.T) {
	// Create a mock store for testing
	mockStore := NewInMemoryStore()

	manager := NewStoreManager(mockStore)

	if manager.Store != mockStore {
		t.Error("StoreManager should hold the provided store")
	}
}
