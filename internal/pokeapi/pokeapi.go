package pokeapi

import (
	"internal/pokecache"
	"net/http"
	"io/ioutil"
	"fmt"
	"errors"
	"encoding/json"
)

var cache pokecache.Cache = pokecache.NewCache(pokecache.Interval) //20 * time.Second, import "time"

func GetLocationAreas(URL string) (LocationAreas, error) {
	body, ok := cache.Get(URL)

	if !ok {
		response, err := http.Get(URL)
		if err != nil {
			return LocationAreas{}, err
		}
		defer response.Body.Close()

		body, err = ioutil.ReadAll(response.Body)
		if response.StatusCode > 299 {
			msg := fmt.Sprintf("Error: response failed. status-code: %d, body: %s", response.StatusCode, body)
			err = errors.New(msg)
			return LocationAreas{}, err
		}
		if err != nil {
			return LocationAreas{}, err
		}
	}

	locationAreasResult := LocationAreas{}
	err := json.Unmarshal(body, &locationAreasResult)
	if err != nil {
		return LocationAreas{}, err
	}

	if !ok {
		cache.Add(URL, body)
	}

	return locationAreasResult, nil
}

func ExploreArea(URL string) (ExploredArea, error) {
	body, ok := cache.Get(URL)
	
	if !ok {
		response, err := http.Get(URL)
		if err != nil {
			return ExploredArea{}, err
		}
		defer response.Body.Close()

		body, err = ioutil.ReadAll(response.Body)
		if response.StatusCode > 299 {
			errString := fmt.Sprintf("Error: response failed. status code: %d, body: %s", response.StatusCode, body)
			err = errors.New(errString)
			return ExploredArea{}, err
		}
		if err != nil {
			return ExploredArea{}, err
		}
	}

	exploredAreaResult := ExploredArea{}
	err := json.Unmarshal(body, &exploredAreaResult)
	if err != nil {
		return ExploredArea{}, err
	}
	if !ok {
		cache.Add(URL, body)
	}

	return exploredAreaResult, nil
}

func GetPokemon(URL string) (Pokemon, error) {
	response, err := http.Get(URL)
	if err != nil {
		return Pokemon{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode > 299 {
		errString := fmt.Sprintf("Error: response failed. status code: %d, body: %s", response.StatusCode, body)
		err = errors.New(errString)
		return Pokemon{}, err
	}
	if err != nil {
		return Pokemon{}, err
	}

	pokemonResult := Pokemon{}
	err = json.Unmarshal(body, &pokemonResult)
	if err != nil {
		return Pokemon{}, err
	}
	return pokemonResult, nil
}

type LocationAreas struct {
	Next     *string  `json:"next"`
	Previous *string  `json:"previous"`
	Results  []struct {
		Name string   `json:"name"`
	}                 `json:"results"`
}

type ExploredArea struct {
	PokemonEncounters []struct {
		Pokemon       struct {
			Name      string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string  `json:"name"`
	Order 		   int     `json:"order"`
	BaseExperience int     `json:"base_experience"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Stats          []struct {
		BaseStat   int     `json:"base_stat"`
		Effort     int     `json:"effort"`
		Stat       struct {
			Name   string  `json:"name"`
			URL    string  `json:"url"`
		} 				   `json:"stat"`
	} 				       `json:"stats"`
	Types          []struct {
		Slot       int     `json:"slot"`
		Type       struct {
			Name   string  `json:"name"`
			URL    string  `json:"url"`
		}                  `json:"type"`
	}                      `json:"types"`
	Abilities      []struct {
		Ability    struct {
			Name   string  `json:"name"`
			URL    string  `json:"url"`
		}                  `json:"ability"`
		IsHidden   bool    `json:"is_hidden"`
		Slot       int     `json:"slot"`
	}                      `json:"abilities"`
	//LocationAreaEncounters string `json:"location_area_encounters"`
}