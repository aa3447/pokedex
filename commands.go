package main

import (
	"fmt"
	"os"
)

type commandTemplate struct {
	name string
	description string
	callback func(config *config) error
}

type config struct {
	next string
	prev string
}

// commandMapGen generates a map of commands with their names, descriptions, and callbacks.
func commandMapGen() map[string]commandTemplate{
	commandMap := map[string]commandTemplate{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "Displays the names of 20 location areas in the Pokemon",
			callback: commandMap,
		},
	}

	// Dynamically generate help command description
	var description string
	for _, command := range commandMap {
		description += fmt.Sprintf("%s: %s\n", command.name, command.description)
	}
	description += "help: Show this help message"
	commandMap["help"] = commandTemplate{
			name: "help",
			description: description,
			callback: commandHelp,
	}

	return commandMap
}

//Command Handlers
func commandExit(config *config) error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
 }

 func commandHelp(config *config) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	return nil
 }

func commandMap(config *config) error {
	mapData ,err := getMap(config, true)
	if err != nil {
		fmt.Println("Error fetching map data:", err)
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	return nil
}