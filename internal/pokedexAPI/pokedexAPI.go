package pokedexapi

import (
	"encoding/json"
	"io"
	"net/http"
)

// var baseUrl string = "https://pokeapi.co/api/v2/"

type LocationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocation(nextURL string) (LocationResponse, error) {
	response, err := http.Get(nextURL)
	if err != nil {
		return LocationResponse{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return LocationResponse{}, err
	}

	location := LocationResponse{}
	if err := json.Unmarshal(body, &location); err != nil {
		return location, err
	}

	return location, nil
}
