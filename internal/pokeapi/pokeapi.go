package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/dd0lynx/pokedex-cli/internal/pokecache"
)

const API = "https://pokeapi.co/api/v2/"

var cache = pokecache.NewCache(1 * time.Minute)

type NamedResourceList struct {
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []NamedResource `json:"results"`
}

type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func callAPI[T any](url string, t *T) error {
	// check cache
	data, ok := cache.Get(url)
	if !ok {
		// no cache for url
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return errors.New("Request failed")
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		cache.Add(url, data)
	}

	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	return nil
}

type LocationArea struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon NamedResource `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// returns a LocationArea if specified otherwise a NamedResourceList
func GetLocationArea[T any](url string, la *T) error {
	return callAPI(url, la)
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		Stat     NamedResource `json:"stat"`
		BaseStat int           `json:"base_stat"`
	} `json:"stats"`
	Types []struct {
		Type NamedResource `json:"type"`
		Slot int           `json:"slot"`
	} `json:"types"`
}

func GetPokemon[T any](url string, p *T) error {
	return callAPI(url, p)
}
