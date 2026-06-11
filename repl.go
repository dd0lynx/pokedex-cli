package main

import (
	"strings"
)

func CleanInput(text string) []string {
	cleaned := strings.Fields(strings.ToLower(text))

	return cleaned
}
