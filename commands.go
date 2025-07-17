package main

import (
	"fmt"
	"os"
)

type commandTemplate struct {
	name string
	description string
	callback func() error
}

// commandMapGen generates a map of commands with their names, descriptions, and callbacks.
func commandMapGen() map[string]commandTemplate{
	commandMap := map[string]commandTemplate{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
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
func commandExit() error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
 }

 func commandHelp() error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	return nil
 }