package main

import (
	"net/http"
	"encoding/json"
	"fmt"
)
func makeGetRequest[T any](url string, result *T) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return fmt.Errorf("couldn't read request: %w", err)
    }
	return nil
}
