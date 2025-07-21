package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/aa3447/pokedex/internal/commands"
	"github.com/aa3447/pokedex/internal/global_structs"
)


 func main() {
	commandMap := commands.CommandMapGen()
	config := &global_structs.Config{
		Next: "",
		Prev: "",
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		
		if !scanner.Scan(){
			fmt.Printf("Error reading input: %v", scanner.Err())
			return
		}

		cleanedInput := cleanInput(scanner.Text())

		command, exist := commandMap[cleanedInput[0]]
		if exist{
			command.Callback(config)
			if command.Name == "help"{
				fmt.Println(command.Description)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", cleanedInput[0])
		}

	}
 }

 func cleanInput(text string) []string {
	loweredTrimmedText := strings.ToLower(strings.TrimSpace(text))

	splitTextOnWhiteSpace := strings.Fields(loweredTrimmedText)

	return splitTextOnWhiteSpace
 }
