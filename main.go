package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dd0lynx/pokedex-cli/internal/pokeapi"
)

type config struct {
	Next     string
	Previous string
	Pokemon  map[string]pokeapi.Pokemon
}

func main() {
	commands := getCommands()
	var config config
	config.Pokemon = map[string]pokeapi.Pokemon{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// wait for input
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := CleanInput(scanner.Text())

		// no input
		if len(text) < 1 {
			fmt.Println("Unknown command")
			continue
		}

		// run command
		command, ok := commands[text[0]]
		if ok {
			if err := command.callback(&config, text); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
