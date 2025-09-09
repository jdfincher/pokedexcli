// Package pokeapi implements the api calls for the pokedex and caches results using package pokecache
package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetConfig(url string) (*Config, error) {
	locations := Config{}
	res, err := http.Get(url)
	if err != nil {
		return &locations, fmt.Errorf("error making http GET request: %w", err)
	}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locations); err != nil {
		return &locations, fmt.Errorf("error decoding response body: %w", err)
	}
	return &locations, nil
}
