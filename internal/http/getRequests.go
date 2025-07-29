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

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order int `json:"order"`
		} `json:"version_group_details"`
	} `json:"moves"`
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
	var mapData MapStruct
	
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

		if err = json.NewDecoder(resp.Body).Decode(&mapData); err != nil {
			return MapStruct{}, err
		}

		cacheMap , err = json.Marshal(mapData)
		if err != nil {
			return MapStruct{}, fmt.Errorf("error marshalling map data: %v", err)
		}
		cache.Add(config.Next, cacheMap) 

		
	} else {
		if err = json.Unmarshal(cacheMap, &mapData); err != nil {
			return MapStruct{}, fmt.Errorf("error unmarshalling cached map data: %v", err)
		}
	}

	config.Next = mapData.Next
	config.Prev = mapData.Previous
	return mapData, nil
}

func GetLocationPokemonCache(config *global_structs.Config, location string) (LocationPokemon, error){
	var err error
	var resp *http.Response
	var cacheExists bool
	var cachePoke []byte
	var pokemonData LocationPokemon
	
	if cache == nil {
		cache = NewCache(5 * time.Minute)
	}
	
	locationUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", location)
	cachePoke, cacheExists = cache.Get(locationUrl)
	
	if !cacheExists{
		resp, err = http.Get(locationUrl)
		if err != nil {
			return LocationPokemon{},err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return LocationPokemon{}, fmt.Errorf("received status code %d", resp.StatusCode)
		}

		
		if err = json.NewDecoder(resp.Body).Decode(&pokemonData); err != nil {
			return LocationPokemon{}, err
		}

		cachePoke , err = json.Marshal(pokemonData)
		if err != nil {
			return LocationPokemon{}, fmt.Errorf("error marshalling map data: %v", err)
		}
		cache.Add(locationUrl, cachePoke) 
	} else {
		if err = json.Unmarshal(cachePoke, &pokemonData); err != nil {
			return LocationPokemon{}, fmt.Errorf("error unmarshalling cached map data: %v", err)
		}
	}

	return pokemonData, nil
}

func GetPokemonDetailsCache(config *global_structs.Config, pokemon string) (Pokemon, error){
	var err error
	var resp *http.Response
	var cacheExists bool
	var cachePoke []byte
	var pokemonData Pokemon
	
	if cache == nil {
		cache = NewCache(5 * time.Minute)
	}
	
	pokemonUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemon)
	cachePoke, cacheExists = cache.Get(pokemonUrl)
	
	if !cacheExists{
		resp, err = http.Get(pokemonUrl)
		if err != nil {
			return Pokemon{},err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return Pokemon{}, fmt.Errorf("received status code %d", resp.StatusCode)
		}

		if err = json.NewDecoder(resp.Body).Decode(&pokemonData); err != nil {
			return Pokemon{}, err
		}

		cachePoke , err = json.Marshal(pokemonData)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error marshalling map data: %v", err)
		}
		cache.Add(pokemonUrl, cachePoke) 
	} else {
		if err = json.Unmarshal(cachePoke, &pokemonData); err != nil {
			return Pokemon{}, fmt.Errorf("error unmarshalling cached map data: %v", err)
		}	
	}
	
	return pokemonData, nil
}
