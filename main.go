package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
    "internal/pokiapi"
    "errors"
    "log"
)

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

func commandHelp(config *Config) error {
	fmt.Println("\nPokedex commands:\n")
	for _, command := range commands {
		fmt.Fprintln(os.Stderr, command.name, ":", command.description)
	}
	fmt.Print("\n")
	return nil
}

func commandExit(config *Config) error {
	fmt.Print("Exit? (y/n) ")
	go getInput(config)
	input := <- config.ch
	input = parseString(input)
	if input == "y" || input == "yes" {
		play = false
	}
	return nil
}

func printLocations(config *Config, url *string) error {
	locations, err := pokiapi.Fetch(url)
	if err != nil {
		return err
	}
	config.next = locations.Next
	config.previous = locations.Previous
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMap(config *Config) error {
	url := config.next
	if url == nil {
		err := errors.New("Error: end of map")
		return err
	}
	return printLocations(config, url)
}

func commandMapBack(config *Config) error {
	url := config.previous
	if url == nil {
		err := errors.New("Error: cannot map back from start")
		return err
	}
	return printLocations(config, url)
}

// capitalize any struct/field needs exporting
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

func getCommandsMap() (map[string]cliCommand) {
	return map[string]cliCommand{
		"map": {
			name:		 "map",
			description: "Displays the next 20 locations",
			callback:	 commandMap,
		},
		"mapb": {
			name:		 "mapb",
			description: "Displays back the previous 20 locations",
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
				play = false
            	return nil
        	},
        },
	}
}

var play bool = true
var commands map[string]cliCommand
var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
var ch = make(chan string)
var pokiAPILocationAreaEndpoint string = "https://pokeapi.co/api/v2/location-area/"
var config *Config = &Config{
	scanner: scanner,
	ch: ch,
	next: &pokiAPILocationAreaEndpoint,
	previous: nil,
}

func main() {
	commands = getCommandsMap()
	for ; play == true; {
		fmt.Print("pokedex > ")
		go getInput(config)
		inputString := <- ch
		parsedString := parseString(inputString)
		command, ok := commands[parsedString]
		if !ok {
			fmt.Println("Invalid input. Try again. Enter 'help' to list available commands. Enter 'q' to quit.")
		} else {
			//fmt.Println("command - " + command.name + ": " + command.description)
			err := command.callback(config)
			if err != nil {
				log.Println(err)
			}
		}
	}
}