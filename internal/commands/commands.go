package commands

// Package commands provides the command handling functionality for the Pokedex application.

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aa3447/pokedex/internal/global_structs"
	"github.com/aa3447/pokedex/internal/http"
)

type CommandTemplate struct {
	Name string
	Description string
	Callback func(config *global_structs.Config,augments ...string) error
}

type Pokedex struct {
	Dex map[string]http.Pokemon
}

var pokedex *Pokedex

// commandMapGen generates a map of commands with their names, descriptions, and callbacks.
func CommandMapGen() map[string]CommandTemplate{
	commandMap := map[string]CommandTemplate{
		"exit": {
			Name: "exit",
			Description: "Exit the Pokedex",
			Callback: commandExit,
		},
		"map": {
			Name: "map",
			Description: "Displays the names of 20 location areas in the Pokemon",
			Callback: commandMap,
		},
		"mapb": {
			Name: "mapb",
			Description: "Displays the previous names of 20 location areas in the Pokemon",
			Callback: commandMapBack,
		},
		"explore": {
			Name: "explore",
			Description: "List all available pokemon in a location",
			Callback: commandExplore,
		},
		"catch": {
			Name: "catch",
			Description: "Catch a pokemon by name",
			Callback: catchPokemon,
		},
	}

	// Dynamically generate help command description
	var description string
	for _, command := range commandMap {
		description += fmt.Sprintf("%s: %s\n", command.Name, command.Description)
	}
	description += "help: Show this help message"
	commandMap["help"] = CommandTemplate{
			Name: "help",
			Description: description,
			Callback: commandHelp,
	}

	return commandMap
}

//Command Handlers
func commandExit(config *global_structs.Config, augments ...string) error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
 }

 func commandHelp(config *global_structs.Config, augments ...string) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	return nil
 }

func commandMap(config *global_structs.Config, augments ...string) error {
	mapData ,err := http.GetMapCache(config, true)
	if err != nil {
		fmt.Println("Error fetching map data:", err)
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapBack(config *global_structs.Config, augments ...string) error {
	mapData ,err := http.GetMapCache(config, false)
	if err != nil {
		fmt.Println("Error fetching map data:", err)
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(config *global_structs.Config, augments ...string) error {
	if len(augments) < 1 {
		fmt.Println("Please provide a location")
		return nil
	}
	pokeData, err := http.GetLocationPokemonCache(config, augments[0])
	
	if err != nil {
		fmt.Println("Error fetching location pokemon data:", err)
		return err
	}
	
	fmt.Printf("Listing all pokemon in %s:\n", augments[0])
	for _, pokemon := range pokeData.Pokemon_encounters{	
		fmt.Println("- " + pokemon.Pokemon.Name)
	} 
	return nil
}

func catchPokemon(config *global_structs.Config, augments ...string) error {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	if pokedex == nil {
		pokedex = &Pokedex{
			Dex: make(map[string]http.Pokemon),
		}
	}

	if len(augments) < 1 {
		fmt.Println("Please provide a pokemon name to catch.")
		return nil
	}
	pokemon, err := http.GetPokemonDetailsCache(config, augments[0])

	if err != nil {
		fmt.Println("Error fetching pokemon data:", err)
		return err
	}
	
	pokemonName := augments[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	catchRate := 100.0 - (float64(pokemon.BaseExperience)/10.0)
	if randGen.Float64()*100.0 < catchRate {
		fmt.Printf("Congratulations! You caught a %s!\n", pokemonName)
		pokedex.Dex[pokemonName] = pokemon
	} else {
		fmt.Printf("Oh no! The %s escaped!\n", pokemonName)
	}
	
	return nil
}

func GetPokedex() *Pokedex {
	if pokedex == nil {
		pokedex = &Pokedex{
			Dex: make(map[string]http.Pokemon),
		}
	}
	return pokedex
}