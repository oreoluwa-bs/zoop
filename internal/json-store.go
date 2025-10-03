package internal

import "fmt"

type JSONStore struct {
	data map[string]string
}

func NewJSONStore() *JSONStore {
	data := make(map[string]string)

	return &JSONStore{
		data: data,
	}
}

func (s *JSONStore) Get(key string) (string, error) {
	val, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("Item with key (%s) not found", key)
	}

	return val, nil
}

func (s *JSONStore) Set(key string, value string) error {
	s.data[key] = value

	return nil
}

func (s *JSONStore) Delete(key string) error {
	delete(s.data, key)
	return nil
}
