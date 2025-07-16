package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

 func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		
		if !scanner.Scan(){
			fmt.Printf("Error reading input: %v", scanner.Err())
			return
		}

		cleanedInput := cleanInput(scanner.Text())

		fmt.Printf("Your command was: %s\n", cleanedInput[0])

	}
 }

 func cleanInput(text string) []string {
	loweredTrimmedText := strings.ToLower(strings.TrimSpace(text))

	splitTextOnWhiteSpace := strings.Fields(loweredTrimmedText)

	return splitTextOnWhiteSpace
 }