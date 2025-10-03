package internal

import "fmt"

type InMemoryStore struct {
	data map[string]string
}

// For testing purposes
func NewInMemoryStore() *InMemoryStore {
	data := make(map[string]string)

	return &InMemoryStore{
		data: data,
	}
}

func (s *InMemoryStore) Get(key string) (string, error) {
	val, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("Item with key (%s) not found", key)
	}

	return val, nil
}

func (s *InMemoryStore) Set(key string, value string) error {
	s.data[key] = value

	return nil
}

func (s *InMemoryStore) Delete(key string) error {
	delete(s.data, key)
	return nil
}

func (s *InMemoryStore) GetAllKeys() ([]string, error) {
	keys := make([]string, 0)
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys, nil
}
