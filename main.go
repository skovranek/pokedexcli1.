package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
    "internal/pokeapi"
    "errors"
    "log"
)

type Config struct {
	scanner  *bufio.Scanner
	ch       chan string
	next     *string
	previous *string
}

type cliCommand struct {
	name		string
	description	string
	callback	func(*Config) error
}

func newCommandsMap() map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:		 "map",
			description: "Displays the next 10 locations",
			callback:	 commandMap,
		},
		"mapb": {
			name:		 "mapb",
			description: "Displays back the previous 10 locations",
			callback:	 commandMapBack,
		},
		"help": {
			name:		 "help",
			description: "Lists available commands",
			callback:	 commandHelp,
		},
		"exit": {
			name:		 "exit",
			description: "Exit the Pokedex",
			callback:	 commandExit,
		},
		"q": {
			name:		 "q",
			description: "Quit",
			callback: func(*Config) error {
				fmt.Println("Goodbye!")
				play = false
				return nil
			},
		},
	}
}

func commandMap(config *Config) error {
	fmt.Println("---------------------\nThe Next 10 Locations\n---------------------")
	url := config.next
	if url == nil {
		err := errors.New("Error: end of map")
		return err
	}
	return printLocations(config, url)
}

func commandMapBack(config *Config) error {
	fmt.Println("-------------------------\nThe Previous 10 Locations\n-------------------------")
	url := config.previous
	if url == nil {
		err := errors.New("Error: cannot map back from start")
		return err
	}
	return printLocations(config, url)
}

func printLocations(config *Config, url *string) error {
	locations, err := pokeapi.Fetch(url)
	if err != nil {
		return err
	}
	config.next = locations.Next
	config.previous = locations.Previous
	for _, location := range locations.Results {
		fmt.Println(strings.Title(strings.ReplaceAll(location.Name, "-", " ")))
	}
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("----------------\nPokedex commands\n----------------")
	for _, command := range commands {
		fmt.Fprintln(os.Stderr, command.name, ":", command.description)
	}
	return nil
}

func commandExit(config *Config) error {
	fmt.Print("Exit? (y/n) ")
	go getInput(config)
	input := <- config.ch
	input = parseString(input)
	if input == "y" || input == "yes" {
		play = false
		fmt.Println("\nGoodbye!")	
	}
	return nil
}

func getInput(config *Config) {
	for config.scanner.Scan() {
		config.ch <- config.scanner.Text()
	}
	if err := config.scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func parseString(toParse string) string {
	toParse = strings.TrimSpace(toParse)
	toParse = strings.ToLower(toParse)
	return toParse
}

var play bool = true
var commands map[string]cliCommand
var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
var ch = make(chan string)
var pokeapiLocationAreaEndpoint string = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=10"
var config *Config = &Config{
	scanner: scanner,
	ch: ch,
	next: &pokeapiLocationAreaEndpoint,
	previous: nil,
}

func main() {
	commands = newCommandsMap()
	fmt.Print("\n")
	for ; play == true; {
		fmt.Print("pokedex > ")
		go getInput(config)
		inputString := <- ch
		parsedString := parseString(inputString)
		command, ok := commands[parsedString]
		fmt.Print("\n")
		if !ok {
			fmt.Println("Invalid input, please try again.\n(Enter 'help' to list available commands)\n(Enter 'q' to quit)")
		} else {
			err := command.callback(config)
			if err != nil {
				log.Println(err)
			}
		}
		fmt.Print("\n")
	}
}