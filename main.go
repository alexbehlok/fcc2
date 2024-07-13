package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Location struct {
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []mapResults `json:"results"`
}

type Config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(Config) error
}

type mapResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var config = Config{}

func main() {
	runningREPL := true

	for runningREPL {
		// Prompt the user for input
		fmt.Print(">> ")

		reader := bufio.NewScanner(os.Stdin)
		input := reader.Scan()

	}
}

func mapCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Prints the help menu",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the program",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
	}
}

func commandExit() error {
	// Exiting REPL
	return nil
}

func commandHelp() error {
	// Display the help message
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Available commands:")
	fmt.Println("help: Displays this help message")
	fmt.Println("exit: Exit the Pokedex\n")
	return nil
}

func commandMap(config *Config) error {
	url := config.Next
	// If no next URL is set, use the initial endpoint
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	location := Location{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Update config with the new pagination URLs
	config.Next = location.Next
	config.Previous = location.Previous

	// Display location names
	for _, loc := range location.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapb(config Config) error {
	url := config.Previous
	// If there is no previous page, print an error
	if url == "" {
		fmt.Println("No previous page available.")
		return nil
	}

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	location := Location{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		log.Fatal(err)
	}

	// Update config with the new pagination URLs
	config.Next = location.Next
	config.Previous = location.Previous

	// Display location names
	for _, loc := range location.Results {
		fmt.Println(loc.Name)
	}

	return nil
}
