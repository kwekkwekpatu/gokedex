package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

var cliName string = "gokedex"

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
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Gokedex!")
	os.Exit(0)
	return nil
}

func displayNext() error {
	response, err := http.Get("")
	return nil
}

func displayPrevious() error {
	return nil
}
