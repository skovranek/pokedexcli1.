package pokeapi

import (
	"internal/pokecache"
	"time"
	"net/http"
	"io/ioutil"
	"fmt"
	"errors"
	"encoding/json"
)

var cache pokecache.Cache = pokecache.NewCache(20 * time.Second) //5 * time.Minute

type LocationAreas struct {
	Count    int      `json:"count"`
	Next     *string  `json:"next"`
	Previous *string  `json:"previous"`
	Results  []struct {
		Name string   `json:"name"`
		URL  string   `json:"url"`
	}                 `json:"results"`
}

// TODO: verify if this should be outside func
var locationAreas *LocationAreas = &LocationAreas{}

func GetLocationAreas(url *string) (*LocationAreas, error) {
	body, ok := cache.Get(*url)
	if !ok {
		response, err := http.Get(*url)
		if err != nil {
			return &LocationAreas{}, err
		}
		defer response.Body.Close()
		body, err = ioutil.ReadAll(response.Body)
		if response.StatusCode > 299 {
			errString := fmt.Sprintf("Response failed:\nstatus code: %d\nbody: %s", response.StatusCode, body)
			err = errors.New(errString)
			return &LocationAreas{}, err
		}
		if err != nil {
			return &LocationAreas{}, err
		}
		cache.Add(*url, body)
	}
	err := json.Unmarshal(body, locationAreas)
	if err != nil {
		return &LocationAreas{}, err
	}
	return locationAreas, nil
}

type ExploredArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []string `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

var exploredArea *ExploredArea = &ExploredArea{}

func ExploreArea(url *string) (*ExploredArea, error) {
	body, ok := cache.Get(*url)
	if !ok {
		response, err := http.Get(*url)
		if err != nil {
			return &ExploredArea{}, err
		}
		defer response.Body.Close()
		body, err = ioutil.ReadAll(response.Body)
		if response.StatusCode > 299 {
			errString := fmt.Sprintf("Response failed:\nstatus code: %d\nbody: %s", response.StatusCode, body)
			err = errors.New(errString)
			return &ExploredArea{}, err
		}
		if err != nil {
			return &ExploredArea{}, err
		}
		cache.Add(*url, body)
	}
	err := json.Unmarshal(body, exploredArea)
	if err != nil {
		return &ExploredArea{}, err
	}
	return exploredArea, nil
}
