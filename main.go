package main

import (
	"bufio"
	"fmt"
	"os"

	pokedexapi "github.com/kwekkwekpatu/gokedex/internal/pokedexAPI"
)

var cliName string = "gokedex"
var nextURL string = "https://pokeapi.co/api/v2/location-area"
var previousURL any

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Gokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 locations",
			callback:    displayNext,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations",
			callback:    displayPrevious,
		},
	}
}

func main() {
	commands := getCommands()

	scanner := bufio.NewScanner(os.Stdin)
	bootGokedex()
	printPromt()
	for scanner.Scan() {

		input := scanner.Text()

		command, exists := commands[input]
		if exists {
			err := command.callback()
			if err != nil {
				fmt.Println("Error executing command: ", err)
			}
		} else {
			fmt.Println("Unknown command: ", input)
		}
		printPromt()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func printPromt() error {
	fmt.Print(cliName, "> ")
	return nil
}

func bootGokedex() error {
	fmt.Println("Booting up the Gokedex!")
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Gokedex!")
	fmt.Println("Usage:")
	fmt.Println("  help: Displays a help message")
	fmt.Println("  exit: Exit the Gokedex")
	fmt.Println("  map: Print the next 20 locations")
	fmt.Println("  mapb: Print the previous 20 locations")
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Gokedex!")
	os.Exit(0)
	return nil
}

func displayNext() error {
	locations, err := pokedexapi.GetLocation(nextURL)
	if err != nil {
		return err
	}

	display(locations)
	return nil
}

func displayPrevious() error {
	if previousURL == nil {
		fmt.Println("You are already at the first locations")
		return fmt.Errorf("No previousURL")
	}
	url, ok := previousURL.(string)
	if !ok {
		return fmt.Errorf("previousURL not a string")
	}
	locations, err := pokedexapi.GetLocation(string(url))
	if err != nil {
		return err
	}

	display(locations)
	return nil
}

func display(locations pokedexapi.LocationResponse) error {
	nextURL = locations.Next
	previousURL = locations.Previous
	for _, location := range locations.Results {
		locationName := location.Name
		fmt.Println(locationName)
	}
	return nil
}
