package http
// Package http provides the HTTP client functionality to fetch data from the PokeAPI.


import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/aa3447/pokedex/internal/global_structs"
)

type MapStruct struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"results"`
}

// getMap fetches the map data from the API and updates the config with the next and previous URLs.
// If nextOrPrev is true, it fetches the next page; if false, it fetches the previous page.
func GetMap(config *global_structs.Config, nextOrPrev bool) (MapStruct, error){
	var err error
	var resp *http.Response

	if config.Next == "" {
		config.Next = "https://pokeapi.co/api/v2/location-area?limit=20"
	}

	if nextOrPrev{
		if config.Next == "" {
			return MapStruct{}, fmt.Errorf("no next page available")
		}
		resp, err = http.Get(config.Next)
	}else{
		if config.Prev == "" || config.Prev == "null" {
			return MapStruct{}, fmt.Errorf("you're on the first page")
		}
		resp, err = http.Get(config.Prev)
	}
	defer resp.Body.Close()

	if err != nil {
		return MapStruct{},err
	}
	if resp.StatusCode != http.StatusOK {
		return MapStruct{}, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	var mapData MapStruct
	if err := json.NewDecoder(resp.Body).Decode(&mapData); err != nil {
		return MapStruct{}, err
	}

	config.Next = mapData.Next
	config.Prev = mapData.Previous
	
	return mapData, nil
}