package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	pokecache "github.com/kwekkwekpatu/gokedex/internal/pokecache"
	pokedexapi "github.com/kwekkwekpatu/gokedex/internal/pokedexAPI"
)

var cliName string = "gokedex"
var nextURL string = "https://pokeapi.co/api/v2/location-area"
var previousURL any
var commandHistory []string
var historyIndex int = -1

type Pokedex struct {
	pokedex map[string]pokedexapi.Pokemon
}

func NewPokedex() *Pokedex {
	return &Pokedex{pokedex: make(map[string]pokedexapi.Pokemon)}
}

func (p *Pokedex) AddPokemon(pokemon pokedexapi.Pokemon) {
	p.pokedex[pokemon.Name] = pokemon
}

func (p *Pokedex) GetPokemon(name string) (pokedexapi.Pokemon, bool) {
	pokemon, exists := p.pokedex[name]
	return pokemon, exists
}

type cliCommand struct {
	name        string
	description string
	callback    func(cache *pokecache.Cache, dex *Pokedex, args ...string) error
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
		"explore": {
			name:        "explore",
			description: "Explore the given location",
			callback:    exploreLocation,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch the selected pokemon",
			callback:    catch,
		},
		"inspect": {
			name:        "inspect",
			description: "Show the data of the selected pokemon if it's in the pokedex",
			callback:    inspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows all the pokemon in your pokedex",
			callback:    showPokedex,
		},
	}
}

func main() {
	commands := getCommands()

	dex := NewPokedex()
	cache := pokecache.NewCache(5 * time.Second)
	scanner := bufio.NewScanner(os.Stdin)

	bootGokedex()
	printPromt()
	for scanner.Scan() {

		input := scanner.Text()
		args := strings.Split(input, " ")
		commandName := args[0]
		args = append(args[:0], args[1:]...)

		command, exists := commands[commandName]
		if exists {
			err := command.callback(cache, dex, args...)
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

func commandHelp(cache *pokecache.Cache, dex *Pokedex, args ...string) error {
	fmt.Println("Welcome to the Gokedex!")
	fmt.Println("Usage:")
	fmt.Println("  help: Displays a help message")
	fmt.Println("  exit: Exit the Gokedex")
	fmt.Println("  map: Print the next 20 locations")
	fmt.Println("  mapb: Print the previous 20 locations")
	fmt.Println("  explore [location]: Shows the pokemon that can be found in the given location")
	fmt.Println("  catch [pokemon]: Attempts to catch the given pokemon")
	fmt.Println("  inspect [pokemon]: Shows the information of the selected pokemon if the pokemon has been added to the pokedex")
	fmt.Println("  pokedex: Show all the pokemon currently in your pokedex")
	return nil
}

func commandExit(ccache *pokecache.Cache, dex *Pokedex, args ...string) error {
	fmt.Println("Closing the Gokedex!")
	os.Exit(0)
	return nil
}

func displayNext(cache *pokecache.Cache, dex *Pokedex, args ...string) error {
	body, err := fetch(nextURL, cache)
	if err != nil {
		return err
	}
	locations, err := pokedexapi.UnmarshalLocations(body)
	if err != nil {
		return err
	}

	display(locations)
	display(locations)
	return nil
}

func displayPrevious(cache *pokecache.Cache, dex *Pokedex, args ...string) error {
	if previousURL == nil {
		fmt.Println("You are already at the first locations")
		return fmt.Errorf("No previousURL")
	}
	url, ok := previousURL.(string)
	if !ok {
		return fmt.Errorf("previousURL not a string")
	}
	body, err := fetch(url, cache)
	if err != nil {
		return err
	}
	locations, err := pokedexapi.UnmarshalLocations(body)
	if err != nil {
		return err
	}

	display(locations)
	return nil
}

func display(locations pokedexapi.LocationsResponse) error {
	nextURL = locations.Next
	previousURL = locations.Previous
	for _, location := range locations.Results {
		locationName := location.Name
		fmt.Println(locationName)
	}
	return nil
}

func exploreLocation(cache *pokecache.Cache, dex *Pokedex, args ...string) error {
	baseURL := "https://pokeapi.co/api/v2/location-area/"
	location := args[0]
	fmt.Println("Exploring " + location + "...")
	if location == "" {
		return fmt.Errorf("Invalid location name")
	}
	url := baseURL + location
	body, err := fetch(url, cache)
	if err != nil {
		return err
	}
	locationData, err := pokedexapi.UnmarshalLocation(body)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, encouter := range locationData.PokemonEncounters {
		pokemon := encouter.Pokemon
		println("- " + pokemon.Name)
	}
	return nil
}

func fetch(url string, cache *pokecache.Cache) ([]byte, error) {
	body, exists := cache.Get(url)
	defer cache.Add(url, body)
	if exists {
		return body, nil
	}
	return pokedexapi.Get(url)
}

func catch(cache *pokecache.Cache, dex *Pokedex, args ...string) error {
	baseUrl := "https://pokeapi.co/api/v2/pokemon/"
	nameOrId := args[0]
	if nameOrId == "" {
		return fmt.Errorf("No pokemon name or id given.")
	}
	url := baseUrl + nameOrId
	body, err := fetch(url, cache)
	if err != nil {
		return err
	}
	pokemonData, err := pokedexapi.UnmarshalPokemon(body)
	if err != nil {
		return err
	}
	pokemonName := pokemonData.Name
	fmt.Println("Throwing a pokeball at " + pokemonName + "...")
	if tryCatchPokemon(pokemonData) {
		dex.AddPokemon(pokemonData)
		fmt.Println(pokemonName + " was caught!")
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Println(pokemonName + " escaped!")
	}
	return nil
}

func tryCatchPokemon(pokemon pokedexapi.Pokemon) bool {
	baseExp := pokemon.BaseExperience
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	threshold := randGen.Intn(100)
	if threshold < 150-baseExp {
		return true
	}
	return false
}

func inspect(cache *pokecache.Cache, dex *Pokedex, args ...string) error {
	if len(args) == 0 || args[0] == "" {
		fmt.Println("No pokemon selected for inspection")
		return nil
	}
	name := args[0]
	pokemon, exists := dex.GetPokemon(name)
	if !exists {
		fmt.Println("you have not caught that pokemon.")
		fmt.Println("Try catching it with the catch command.")
		return nil
	}
	printPokemon(pokemon)
	return nil
}

func printPokemon(pokemon pokedexapi.Pokemon) {
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", pokemon.Name, pokemon.Height, pokemon.Weight)

	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokeType := range pokemon.Types {
		fmt.Printf(" - %s\n", pokeType.Type.Name)
	}
}

func showPokedex(cache *pokecache.Cache, dex *Pokedex, args ...string) error {
	if len(dex.pokedex) == 0 {
		fmt.Println("Your pokedex is empty.")
		fmt.Println("Try catching some pokemon first!")
		return nil
	}
	fmt.Println("Your pokedex:")
	for _, pokemon := range dex.pokedex {
		fmt.Println(" - " + pokemon.Name)
	}
	return nil
}

func addCommand(command string) {
	commandHistory = append(commandHistory, command)
}

func handleKeyPress(key string, currentCommand *string) {
	if key == "up" {
		if historyIndex < len(commandHistory)-1 {
			historyIndex++
			*currentCommand = commandHistory[len(commandHistory)-1-historyIndex]
		}
	} else if key == "down" {
		if historyIndex > 0 {
			historyIndex--
			*currentCommand = commandHistory[len(commandHistory)-1-historyIndex]
		} else {
			historyIndex = -1
			*currentCommand = ""
		}
	}
}
