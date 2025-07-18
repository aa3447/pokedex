package main


import (
	"fmt"
	"net/http"
	"encoding/json"
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

// getMap fetches the map data from the API and updates the config with next and previous URLs.
// If nextOrPrev is true, it fetches the next page; if false, it fetches the previous page.
func getMap(config *config, nextOrPrev bool) (MapStruct, error){
	var err error
	var resp *http.Response

	if config.next == "" {
		config.next = "https://pokeapi.co/api/v2/location-area?limit=20"
	}

	if nextOrPrev{
		if config.next == "" {
			fmt.Println("No next page available.")
			return MapStruct{}, fmt.Errorf("no next page available")
		}
		resp, err = http.Get(config.next)
	}else{
		if config.prev == "" {
			fmt.Println("No previous page available.")
			return MapStruct{}, fmt.Errorf("no previous page available")
		}
		resp, err = http.Get(config.prev)
	}
	defer resp.Body.Close()
	
	if err != nil {
		fmt.Println("Error fetching map data:", err)
		return MapStruct{},err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		return MapStruct{}, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	var mapData MapStruct
	if err := json.NewDecoder(resp.Body).Decode(&mapData); err != nil {
		fmt.Println("Error decoding map data:", err)
		return MapStruct{}, err
	}

	config.next = mapData.Next
	config.prev = mapData.Previous
	
	return mapData, nil
}