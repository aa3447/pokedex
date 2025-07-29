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
	Name        string
	Description string
	Callback    func(config *global_structs.Config, arguments ...string) error
}

type Pokedex struct {
	Dex map[string]http.Pokemon
}

var pokedex *Pokedex

// commandMapGen generates a map of commands with their names, descriptions, and callbacks.
func CommandMapGen() map[string]CommandTemplate {
	commandMap := map[string]CommandTemplate{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Displays the names of 20 location areas in the Pokemon",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the previous names of 20 location areas in the Pokemon",
			Callback:    commandMapBack,
		},
		"explore": {
			Name:        "explore",
			Description: "List all available pokemon in a location",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Catch a pokemon by name",
			Callback:    catchPokemon,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a pokemon by name",
			Callback:    inspectPokemon,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "List all pokemon you have caught",
			Callback:    listOwnedPokemon,
		},
	}

	// Dynamically generate help command description
	var description string
	for _, command := range commandMap {
		description += fmt.Sprintf("%s: %s\n", command.Name, command.Description)
	}
	description += "help: Show this help message"
	commandMap["help"] = CommandTemplate{
		Name:        "help",
		Description: description,
		Callback:    commandHelp,
	}

	return commandMap
}

// Command Handlers

// commandExit exits the Pokedex application.
func commandExit(config *global_structs.Config, arguments ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
// commandHelp displays the help message for all commands.
func commandHelp(config *global_structs.Config, arguments ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	return nil
}

// commandMap fetches and displays the names of 20 location areas in the Pokemon.
func commandMap(config *global_structs.Config, arguments ...string) error {
	mapData, err := http.GetMapCache(config, true)
	if err != nil {
		fmt.Println("Error fetching map data:", err)
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}

// commandMapBack fetches and displays the previous names of 20 location areas in the Pokemon.
func commandMapBack(config *global_structs.Config, arguments ...string) error {
	mapData, err := http.GetMapCache(config, false)
	if err != nil {
		fmt.Println("Error fetching map data:", err)
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}

// commandExplore lists all available pokemon in a specified location.
// It requires a location name as an argument.
func commandExplore(config *global_structs.Config, arguments ...string) error {
	if len(arguments)  < 1 {
		fmt.Println("Please provide a location")
		return nil
	}
	pokeData, err := http.GetLocationPokemonCache(config, arguments[0])

	if err != nil {
		fmt.Println("Error fetching location pokemon data:", err)
		return err
	}

	fmt.Printf("Listing all pokemon in %s:\n", arguments[0])
	for _, pokemon := range pokeData.Pokemon_encounters {
		fmt.Println("- " + pokemon.Pokemon.Name)
	}
	return nil
}

// catchPokemon attempts to catch a pokemon by name.
// It requires a pokemon name as an argument.
func catchPokemon(config *global_structs.Config, arguments ...string) error {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	if pokedex == nil {
		pokedex = &Pokedex{
			Dex: make(map[string]http.Pokemon),
		}
	}

	if len(arguments)  < 1 {
		fmt.Println("Please provide a pokemon name to catch.")
		return nil
	}
	pokemon, err := http.GetPokemonDetailsCache(config, arguments[0])

	if err != nil {
		fmt.Println("Error fetching pokemon data:", err)
		return err
	}

	pokemonName := arguments[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	catchRate := 100.0 - ((float64(pokemon.BaseExperience) * 1.5) / 10.0)
	if catchRate < 10.0 {
		catchRate = 10.0
	}

	if randGen.Float64()*100.0 < catchRate {
		fmt.Printf("Congratulations! You caught a %s!\n", pokemonName)
		pokedex.Dex[pokemonName] = pokemon
	} else {
		fmt.Printf("Oh no! The %s escaped!\n", pokemonName)
	}

	return nil
}

// inspectPokemon inspects a pokemon by name and displays its details.
// It requires a pokemon name as an argument.
func inspectPokemon(config *global_structs.Config, arguments ...string) error {
	if len(arguments)  < 1 {
		fmt.Println("Please provide a pokemon name to inspect.")
		return nil
	}
	pokemonName := arguments[0]
	pokemon, exists := pokedex.Dex[pokemonName]
	if !exists {
		fmt.Printf("No pokemon named %s found in your Pokedex.\n", pokemonName)
		return nil
	}

	fmt.Printf("Inspecting %s:\n", pokemonName)
	fmt.Printf("Base Experience: %d\n", pokemon.BaseExperience)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("- %s\n", t.Type.Name)
	}
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("- %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Abilities:")
	for _, ability := range pokemon.Abilities {
		fmt.Printf("- %s\n", ability.Ability.Name)
	}
	return nil
}

// listOwnedPokemon lists all pokemon that have been caught by the user.
func listOwnedPokemon(config *global_structs.Config, arguments ...string) error {
	if pokedex == nil  {
		pokedex = &Pokedex{
			Dex: make(map[string]http.Pokemon),
		}
		fmt.Println("You have no pokemon in your Pokedex.")
		return nil
	}
	if len(pokedex.Dex) == 0 {
		fmt.Println("You have no pokemon in your Pokedex.")
		return nil
	}

	fmt.Println("Your Pokedex contains the following pokemon:")
	for name := range pokedex.Dex {
		fmt.Println("- " + name)
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
