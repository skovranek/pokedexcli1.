package pokiapi

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"errors"
	"encoding/json"
)

type Locations struct {
	Count    int      `json:"count"`
	Next     *string  `json:"next"`
	Previous *string  `json:"previous"`
	Results  []struct {
		Name string   `json:"name"`
		URL  string   `json:"url"`
	}                 `json:"results"`
}

// TODO: verify correct pointer to struct, and if outside func
var locations *Locations = &Locations{}

func Fetch(url *string) (*Locations, error) {
	response, err := http.Get(*url)
	if err != nil {
		return &Locations{}, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode > 299 {
		errString := fmt.Sprintf("Response failed:\nstatus code: %d\nbody: %s", response.StatusCode, body)
		err = errors.New(errString)
		return &Locations{}, err
	}
	if err != nil {
		return &Locations{}, err
	}
	err = json.Unmarshal(body, locations)
	if err != nil {
		return &Locations{}, err
	}
	return locations, nil
}