package commands
// Package commands provides the command handling functionality for the Pokedex application.

import (
	"fmt"
	"os"
	"github.com/aa3447/pokedex/internal/http"
	"github.com/aa3447/pokedex/internal/global_structs"
)

type CommandTemplate struct {
	Name string
	Description string
	Callback func(*global_structs.Config) error
}

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
func commandExit(config *global_structs.Config) error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
 }

 func commandHelp(config *global_structs.Config) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	return nil
 }

func commandMap(config *global_structs.Config) error {
	mapData ,err := http.GetMap(config, true)
	if err != nil {
		fmt.Println("Error fetching map data:", err)
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapBack(config *global_structs.Config) error {
	mapData ,err := http.GetMap(config, false)
	if err != nil {
		fmt.Println("Error fetching map data:", err)
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}