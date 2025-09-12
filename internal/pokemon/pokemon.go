// Package pokemon implements the structs and methods of a Pokemon and Pokedex
package pokemon

import (
	"fmt"
	"math/rand"
)

type Pokedex struct {
	Dex     map[string]Pokemon
	Target  Pokemon
	Current Pokemon
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
	} `json:"held_items"`
	Moves []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
	} `json:"moves"`
	Name    string `json:"name"`
	Order   int    `json:"order"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func (p *Pokedex) Add() {
	p.Dex[p.Target.Name] = p.Target
}

func (p *Pokedex) Catch() bool {
	chance := rand.Intn(p.Target.BaseExperience)
	base := rand.Intn(p.Target.BaseExperience)
	if chance >= base {
		p.Add()
		return true
	}
	return false
}

func (p *Pokedex) Find(name string) bool {
	v, ok := p.Dex[name]
	if !ok {
		fmt.Println("You have not caught " + name + " yet!")
		return false
	}
	p.Current = v
	return true
}

func NewPokedex() *Pokedex {
	m := make(map[string]Pokemon)
	p := &Pokedex{
		Dex:     m,
		Target:  *new(Pokemon),
		Current: *new(Pokemon),
	}
	return p
}
