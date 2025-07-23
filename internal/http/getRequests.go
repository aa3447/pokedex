package http

// Package http provides the HTTP client functionality to fetch data from the PokeAPI.
// This file deals with fetching map data from the API.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

var cache *Cache

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

func GetMapCache(config *global_structs.Config, nextOrPrev bool) (MapStruct, error){
	var err error
	var resp *http.Response
	var cacheExists bool
	
	if cache == nil {
		cache = NewCache(5 * time.Minute)
	}

	if config.Next == "" {
		config.Next = "https://pokeapi.co/api/v2/location-area?limit=20"
	}

	
	if nextOrPrev{
		if config.Next == "" {
			return MapStruct{}, fmt.Errorf("no next page available")
		}
		
		cacheExists = cache.Exists(config.Next)
		if !cacheExists {
			resp, err = http.Get(config.Next)
		}
	}else{
		if config.Prev == "" || config.Prev == "null" {
			return MapStruct{}, fmt.Errorf("you're on the first page")
		}
		
		cacheExists = cache.Exists(config.Prev)
		if !cacheExists {
			resp, err = http.Get(config.Prev)
		}
	}
	
	
	if !cacheExists{
		if err != nil {
			return MapStruct{},err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return MapStruct{}, fmt.Errorf("received status code %d", resp.StatusCode)
		}

		var mapData MapStruct
		if err := json.NewDecoder(resp.Body).Decode(&mapData); err != nil {
			return MapStruct{}, err
		}

		cacheMap , err := json.Marshal(mapData)
		if err != nil {
			return MapStruct{}, fmt.Errorf("error marshalling map data: %v", err)
		}
		cache.Add(config.Next, cacheMap) 

		config.Next = mapData.Next
		config.Prev = mapData.Previous
		return mapData, nil
	} else {
		var cacheMap []byte
		var exists bool

		if nextOrPrev{
			cacheMap, exists = cache.Get(config.Next)
		} else {
			cacheMap, exists = cache.Get(config.Prev)
		}
		
		if !exists {
			return MapStruct{}, fmt.Errorf("cache entry not found for %s", config.Next)
		}

		var mapData MapStruct
		if err := json.Unmarshal(cacheMap, &mapData); err != nil {
			return MapStruct{}, fmt.Errorf("error unmarshalling cached map data: %v", err)
		}

		config.Next = mapData.Next
		config.Prev = mapData.Previous
		return mapData, nil
	}
}