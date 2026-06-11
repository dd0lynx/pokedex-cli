package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/dd0lynx/pokedex-cli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, args []string) error
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
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Shows 20 location areas in the Pokemon world. Each call will show the next 20 location",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Like the map command but shows the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Lists the pokemon that can be encountered in an area. takes a location or id i.e. explore canalave-city-area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon. Caught pokemon can be inpected.",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Shows information about a caught pokemon.",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows a list of all caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func commandHelp(c *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage: ")
	fmt.Println("")

	for k, v := range getCommands() {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func commandExit(c *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(c *config, args []string) error {
	locationAreas := pokeapi.NamedResourceList{}
	if err := pokeapi.GetLocationArea(getNextMapURL(c), &locationAreas); err != nil {
		return err
	}

	c.Next = locationAreas.Next
	c.Previous = locationAreas.Previous
	for _, la := range locationAreas.Results {
		fmt.Println(la.Name)
	}

	return nil
}

func getNextMapURL(c *config) string {
	if c.Next == "" {
		// return "https://pokeapi.co/api/v2/location-area/"
		return pokeapi.API + "location-area/"
	}
	return c.Next
}

func commandMapb(c *config, args []string) error {
	if c.Previous == "" {
		return errors.New("you're on the first page")
	}
	locationAreas := pokeapi.NamedResourceList{}
	if err := pokeapi.GetLocationArea(getPreviousMapURL(c), &locationAreas); err != nil {
		return err
	}

	c.Next = locationAreas.Next
	c.Previous = locationAreas.Previous
	for _, la := range locationAreas.Results {
		fmt.Println(la.Name)
	}

	return nil
}

func getPreviousMapURL(c *config) string {
	if c.Previous == "" {
		return pokeapi.API + "location-area/"
	}
	return c.Previous
}

func commandExplore(c *config, args []string) error {
	if len(args) < 1 {
		return errors.New("Please input a location to explore")
	}
	fmt.Printf("Exploring %s...", args[1])

	locationArea := pokeapi.LocationArea{}
	if err := pokeapi.GetLocationArea(pokeapi.API+"location-area/"+args[1], &locationArea); err != nil {
		return err
	}

	fmt.Println("Found Pokemon: ")
	for _, pokemon := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config, args []string) error {
	if len(args) < 1 {
		return errors.New("Please input a pokemon to catch")
	}
	if _, ok := c.Pokemon[args[1]]; ok {
		fmt.Printf("%s has already been caught!\n", args[1])
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", args[1])

	pokemon := pokeapi.Pokemon{}
	if err := pokeapi.GetLocationArea(pokeapi.API+"pokemon/"+args[1], &pokemon); err != nil {
		return err
	}

	const catchNum = 30
	if rand.Intn(pokemon.BaseExperience) < catchNum {
		// caught
		fmt.Printf("%s was caught!\n", args[1])
		c.Pokemon[args[1]] = pokemon
	} else {
		// escaped
		fmt.Printf("%s escaped!\n", args[1])
	}
	return nil
}

func commandInspect(c *config, args []string) error {
	if len(args) < 1 {
		return errors.New("Please input a pokemon to inspect")
	}
	pokemon, ok := c.Pokemon[args[1]]
	if !ok {
		fmt.Printf("%s hasn't been caught yet. Go catch them first!\n", args[1])
		return nil
	}

	fmt.Printf("Height: %d\nWeight: %d\n", pokemon.Height, pokemon.Weight)
	fmt.Println("Stats: ")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types: ")
	for _, stat := range pokemon.Types {
		fmt.Printf("  -%s\n", stat.Type.Name)
	}

	return nil
}

func commandPokedex(c *config, args []string) error {
	fmt.Println("Your Pokedex: ")
	if len(c.Pokemon) < 1 {
		fmt.Println("Your Pokedex is empty, go catch some pokemon!")
		return nil
	}

	for pokemon := range c.Pokemon {
		fmt.Printf("  - %s\n", pokemon)
	}

	return nil
}
