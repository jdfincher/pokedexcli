// Package pokeapi implements the api calls for the pokedex and caches results using package pokecache
package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jdfincher/pokedexcli/internal/pokecache"
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

type Client struct {
	Cache   *pokecache.Cache
	BaseURL string
	Cfg     *Config
}

func NewClient(interval time.Duration) *Client {
	client := &Client{
		Cache:   pokecache.NewCache(interval),
		BaseURL: "https://pokeapi.co/api/v2/",
		Cfg:     new(Config),
	}
	return client
}

func (c *Client) Get(url string) (*Config, error) {
	if c == nil {
		return nil, errors.New("error: cache is not initialized or is nil")
	}
	v, ok := c.Cache.Find(url)
	if !ok {
		_ = v
		d, err := FetchData(url)
		if err != nil {
			return nil, err
		} else {
			c.Cache.Add(url, d)
			c.Cfg, err = configDecoder(d)
			if err != nil {
				return nil, err
			}
			return c.Cfg, nil
		}
	} else {
		var err error
		c.Cfg, err = configDecoder(v)
		if err != nil {
			return nil, err
		}
		return c.Cfg, nil
	}
}

func FetchData(url string) ([]byte, error) {
	var data []byte
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http.NewRequest: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error from http.Response: %w", err)
	}
	defer res.Body.Close()

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading data from response body: %w", err)
	}
	return data, nil
}

func configDecoder(data []byte) (*Config, error) {
	cfg := new(Config)
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error could not unmarshal data : %w", err)
	}
	return cfg, nil
}
