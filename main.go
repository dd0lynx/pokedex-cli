package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := CleanInput(scanner.Text())

		fmt.Printf("Your command was: %s\n", text[0])
	}
}
