package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


 func main() {
	commandMap := commandMapGen()
	config := &config{
		next: "",
		prev: "",
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
			command.callback(config)
			if command.name == "help"{
				fmt.Println(command.description)
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
