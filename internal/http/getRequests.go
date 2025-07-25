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

type LocationPokemon struct {
	Encounter_method_rates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Encounter_method struct {
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"encounter_method"`
	Game_index int `json:"game_index"`
	Id int `json:"id"`
	Location struct {
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	Pokemon_encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL string `json:"url"`
		} `json:"pokemon"`
		Version_details []struct {
			Encounter_details []struct {
				Chance int `json:"chance"`
				Condition_values []struct {} `json:"condition_values"`
				Method struct {
					Name string `json:"name"`
					URL string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
			Max_chance int `json:"max_chance"`
			Version struct {
				Name string `json:"name"`
				URL string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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

// GetMapCache fetches the map data from the cache or the API if not cached.
func GetMapCache(config *global_structs.Config, nextOrPrev bool) (MapStruct, error){
	var err error
	var resp *http.Response
	var cacheExists bool
	var cacheMap []byte
	
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
		
		cacheMap, cacheExists = cache.Get(config.Next)
		if !cacheExists {
			resp, err = http.Get(config.Next)
		}
	}else{
		if config.Prev == "" || config.Prev == "null" {
			return MapStruct{}, fmt.Errorf("you're on the first page")
		}
		
		cacheMap, cacheExists = cache.Get(config.Prev)
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
		if err = json.NewDecoder(resp.Body).Decode(&mapData); err != nil {
			return MapStruct{}, err
		}

		cacheMap , err = json.Marshal(mapData)
		if err != nil {
			return MapStruct{}, fmt.Errorf("error marshalling map data: %v", err)
		}
		cache.Add(config.Next, cacheMap) 

		config.Next = mapData.Next
		config.Prev = mapData.Previous
		return mapData, nil
	} else {

		var mapData MapStruct
		if err = json.Unmarshal(cacheMap, &mapData); err != nil {
			return MapStruct{}, fmt.Errorf("error unmarshalling cached map data: %v", err)
		}

		config.Next = mapData.Next
		config.Prev = mapData.Previous
		return mapData, nil
	}
}

func GetLocationPokemonCache(config *global_structs.Config, location string) (LocationPokemon, error){
	var err error
	var resp *http.Response
	var cacheExists bool
	var cachePoke []byte
	
	if cache == nil {
		cache = NewCache(5 * time.Minute)
	}
	
	locationUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", location)
	cachePoke, cacheExists = cache.Get(locationUrl)
	if !cacheExists {
		resp, err = http.Get(locationUrl)
	}
	
	
	if !cacheExists{
		if err != nil {
			return LocationPokemon{},err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return LocationPokemon{}, fmt.Errorf("received status code %d", resp.StatusCode)
		}

		var pokemonData LocationPokemon
		if err = json.NewDecoder(resp.Body).Decode(&pokemonData); err != nil {
			return LocationPokemon{}, err
		}

		cachePoke , err = json.Marshal(pokemonData)
		if err != nil {
			return LocationPokemon{}, fmt.Errorf("error marshalling map data: %v", err)
		}
		cache.Add(locationUrl, cachePoke) 

		return pokemonData, nil
	} else {
		var pokemonData LocationPokemon
		if err = json.Unmarshal(cachePoke, &pokemonData); err != nil {
			return LocationPokemon{}, fmt.Errorf("error unmarshalling cached map data: %v", err)
		}
		return pokemonData, nil
	}
}