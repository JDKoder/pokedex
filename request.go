package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/JDKoder/pokedex/internal"
	"time"
	"io"
)

var (
	cache_duration, _ = time.ParseDuration(cache_interval) 
	requestCache = internal.NewCache(cache_duration)
)

const (
	cache_interval =  "15s"
)

func makeGetRequest[T any](url string, result *T) error {
	//check the cache
	var foobar []byte
	foobar, ok := requestCache.Get(url)
	if !ok {
		//cache miss
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		
		foobar, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		//add to cache for next time
		go requestCache.Add(url, foobar)
	}
	unmarshalErr := json.Unmarshal(foobar, &result)
	if unmarshalErr != nil {
		return fmt.Errorf("couldn't read request: %w", unmarshalErr)
    }
	return nil
}
