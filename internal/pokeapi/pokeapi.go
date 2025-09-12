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
	"github.com/jdfincher/pokedexcli/internal/pokemon"
)

type PokemonEncounters struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Client struct {
	Cache         *pokecache.Cache
	BaseURL       string
	Loc           *Locations
	PokEncounters *PokemonEncounters
	Pokedex       *pokemon.Pokedex
}

func NewClient(interval time.Duration) *Client {
	client := &Client{
		Cache:         pokecache.NewCache(interval),
		BaseURL:       "https://pokeapi.co/api/v2/",
		Loc:           new(Locations),
		PokEncounters: new(PokemonEncounters),
		Pokedex:       pokemon.NewPokedex(),
	}
	return client
}

func (c *Client) GetLocations(url string) (*Client, error) {
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
			c.Loc, err = locationsDecoder(d)
			if err != nil {
				return nil, err
			}
			return c, nil
		}
	} else {
		var err error
		c.Loc, err = locationsDecoder(v)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}

func (c *Client) GetPokEncounters(url string) (*Client, error) {
	if c == nil {
		return c, errors.New("error: cache is not initialized or is nil")
	}
	v, ok := c.Cache.Find(url)
	if !ok {
		_ = v
		d, err := FetchData(url)
		if err != nil {
			return c, err
		} else {
			c.Cache.Add(url, d)
			c.PokEncounters, err = pokemonEDecoder(d)
			if err != nil {
				return c, err
			}
			return c, nil
		}
	} else {
		var err error
		c.PokEncounters, err = pokemonEDecoder(v)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}

func (c *Client) GetPokemon(url string) (*Client, error) {
	if c == nil {
		return c, errors.New("error cache is not initialized or is nil")
	}
	v, ok := c.Cache.Find(url)
	if !ok {
		_ = v
		d, err := FetchData(url)
		if err != nil {
			return c, err
		}
		c.Cache.Add(url, d)
		c.Pokedex.Target, err = pokeDecoder(d)
		if err != nil {
			return c, err
		}
		return c, nil
	} else {
		var err error
		c.Pokedex.Target, err = pokeDecoder(v)
		if err != nil {
			return c, err
		}
		return c, nil
	}
}

func FetchData(url string) ([]byte, error) {
	var data []byte
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http.NewRequest: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error from http.Response: %w", err)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: http status code: %v", res.StatusCode)
	}
	defer res.Body.Close()

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading data from response body: %w", err)
	}
	return data, nil
}

func locationsDecoder(data []byte) (*Locations, error) {
	Loc := new(Locations)
	if err := json.Unmarshal(data, &Loc); err != nil {
		return nil, fmt.Errorf("error: could not unmarshal data\ndetails: %w", err)
	}
	return Loc, nil
}

func pokemonEDecoder(data []byte) (*PokemonEncounters, error) {
	Pok := new(PokemonEncounters)
	if err := json.Unmarshal(data, &Pok); err != nil {
		return nil, fmt.Errorf("error: could not unmarshal data\ndetails: %w", err)
	}
	return Pok, nil
}

func pokeDecoder(data []byte) (pokemon.Pokemon, error) {
	pok := *new(pokemon.Pokemon)
	if err := json.Unmarshal(data, &pok); err != nil {
		return pok, fmt.Errorf("error: could not unmarshal data\ndtails: %w", err)
	}
	return pok, nil
}
