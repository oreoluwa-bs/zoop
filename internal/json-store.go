package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONStore struct {
	filePath string
	data     map[string]string
}

func NewJSONStore(filePath string) *JSONStore {
	datasource, err := loadJSONDatasource(filePath)
	if err != nil {
		panic(err)
	}

	return &JSONStore{
		filePath: filePath,
		data:     datasource.Data,
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

	err := writeJSONDatasource(s.filePath, JSONDatasource{
		Data: s.data,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *JSONStore) Delete(key string) error {
	delete(s.data, key)
	return nil
}

type JSONDatasource struct {
	Data map[string]string `json:"data"`
}

func loadJSONDatasource(path string) (JSONDatasource, error) {
	var src JSONDatasource

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// file missing → create with defaults
		src = JSONDatasource{Data: make(map[string]string)}
		data, err := json.MarshalIndent(src, "", "  ")
		if err != nil {
			return src, err
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return src, err
		}
		return src, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return src, err
	}

	if err := json.Unmarshal(data, &src); err != nil {
		return src, err
	}

	return src, nil
}

func writeJSONDatasource(path string, src JSONDatasource) error {
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}
