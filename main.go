package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

type config struct {
	Next     string
	Previous string
}

type NamedAPIResourceList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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
	}
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage: \n")

	for k, v := range getCommands() {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(c *config) error {
	locationAreas := NamedAPIResourceList{}
	if err := getLocationAreas(getNextMapURL(c), &locationAreas); err != nil {
		return err
	}

	c.Next = locationAreas.Next
	c.Previous = locationAreas.Previous
	for _, la := range locationAreas.Results {
		fmt.Println(la.Name)
	}

	return nil
}

func commandMapb(c *config) error {
	if c.Previous == "" {
		return errors.New("you're on the first page")
	}
	locationAreas := NamedAPIResourceList{}
	if err := getLocationAreas(getPreviousMapURL(c), &locationAreas); err != nil {
		return err
	}

	c.Next = locationAreas.Next
	c.Previous = locationAreas.Previous
	for _, la := range locationAreas.Results {
		fmt.Println(la.Name)
	}

	return nil
}

func getLocationAreas(url string, la *NamedAPIResourceList) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, la); err != nil {
		return err
	}
	return nil
}

func getNextMapURL(c *config) string {
	if c.Next == "" {
		return "https://pokeapi.co/api/v2/location-area/"
	}
	return c.Next
}

func getPreviousMapURL(c *config) string {
	if c.Previous == "" {
		return "https://pokeapi.co/api/v2/location-area/"
	}
	return c.Previous
}

func main() {
	commands := getCommands()
	var config config
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// wait for input
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := CleanInput(scanner.Text())

		// run command
		command, ok := commands[text[0]]
		if ok {
			if err := command.callback(&config); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
