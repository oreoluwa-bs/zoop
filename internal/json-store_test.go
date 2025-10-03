package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func initTest(t *testing.T) *JSONStore {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "data-test.json")

	s, err := NewJSONStore(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestJSONStore_SetGet(t *testing.T) {
	s := initTest(t)

	err := s.Set("set", "get")
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.Get("set")
	if err != nil {
		t.Fatal(err)
	}
	if val != "get" {
		t.Errorf("Got %s; want %s", "get", "set")
	}
}

func TestJSONStore_SetDelete(t *testing.T) {
	s := initTest(t)

	err := s.Set("set", "delete")
	if err != nil {
		t.Fatal(err)
	}

	err = s.Delete("set")
	if err != nil {
		t.Fatal(err)
	}

}

func TestJSONStore_GetNotFound(t *testing.T) {
	s := initTest(t)
	_, err := s.Get("notset")
	if err == nil {
		t.Error("Got nil; want err")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Got '%s'; want '%s'", err.Error(), "Item with key (notset) not found")
	}
}

func TestJSONStore_NewFileCreation(t *testing.T) {
	// Use temp dir so we don't pollute filesystem
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "new_store.json")

	// Verify file doesn't exist initially
	if _, err := os.Stat(filePath); err == nil {
		t.Fatal("File should not exist initially")
	}

	// Create store - should create file automatically
	_, err := NewJSONStore(filePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("File should have been created")
	}

	// Verify file has correct structure
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	var datasource JSONDatasource
	if err := json.Unmarshal(data, &datasource); err != nil {
		t.Fatalf("Invalid JSON structure: %v", err)
	}

	if len(datasource.Data) != 0 {
		t.Fatalf("Expected empty data map, got %d items", len(datasource.Data))
	}
}

func TestJSONStore_ConcurrentAccess(t *testing.T) {
	s := initTest(t)
	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperations)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d_%d", goroutineID, j)
				value := fmt.Sprintf("value_%d_%d", goroutineID, j)

				if err := s.Set(key, value); err != nil {
					errors <- fmt.Errorf("goroutine %d, op %d: %v", goroutineID, j, err)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}

	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numOperations; j++ {
			key := fmt.Sprintf("key_%d_%d", i, j)
			expectedValue := fmt.Sprintf("value_%d_%d", i, j)

			actualValue, err := s.Get(key)
			if err != nil {
				t.Errorf("Failed to get %s: %v", key, err)
			}
			if actualValue != expectedValue {
				t.Errorf("Expected %s, got %s", expectedValue, actualValue)
			}
		}
	}
}
