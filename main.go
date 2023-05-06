package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	//"errors"
	//"sort"
)

func getInput(scanner *bufio.Scanner, ch chan string) {
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func parseString(toParse string) string {
	toParse = strings.TrimSpace(toParse)
	toParse = strings.ToLower(toParse)
	return toParse
}

func commandHelp() error {
	fmt.Println("\nPokedex commands:\n")
	for _, command := range commands {
		fmt.Fprintln(os.Stderr, command.name, ":", command.description)
	}
	fmt.Print("\n")
	return nil
}

func commandExit() error {
	fmt.Print("Exit? (y/n) ")
	go getInput(scanner, ch)
	input := <- ch
	input = parseString(input)
	if input == "y" || input == "yes" {
		play = false
	}
	return nil
}

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

func getCommandsMap() (map[string]cliCommand) {
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
	}
}

var play bool = true
var commands map[string]cliCommand
var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
var ch = make(chan string)

func main() {
	commands = getCommandsMap()

	for ; play == true; {
		fmt.Print("pokedex > ")
		go getInput(scanner, ch)

		inputString := <- ch
		parsedString := parseString(inputString)
		command, ok := commands[parsedString]
		if !ok {
			fmt.Println("invalid input, try again")
		} else {
			//fmt.Println("command - " + command.name + ": " + command.description)
			command.callback()
		}
	}
}