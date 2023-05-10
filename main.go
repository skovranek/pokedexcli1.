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

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

func newCommandsMap() map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:		 "map",
			description: fmt.Sprint("Displays the next "+limitParam+" areas"),
			callback:	 commandMap,
		},
		"mapb": {
			name:		 "mapb",
			description: fmt.Sprint("Displays back the previous "+limitParam+" areas"),
			callback:	 commandMapBack,
		},
		"explore": {
			name:		 "explore <area>",
			description: "Explore an area to search for pokemon",
			callback:	 commandExplore,
		},
		"catch": {
			name: 		 "catch",
			description: "Catch a pokemon and add it to the Pokedex",
			callback:	 commandCatch,
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
			callback: func() error {
				fmt.Println("Goodbye!")
				play = false
				return nil
			},
		},
	}
}

func commandMap() error {
	line := "---------------------"
	fmt.Println(line+"\nThe Next "+limitParam+" Areas\n"+line)
	if nextAreasURL == nil {
		err := errors.New("Error: end of map")
		return err
	}
	return printAreas(nextAreasURL)
}

func commandMapBack() error {
	line := "-------------------------"
	fmt.Println(line+"\nThe Previous "+limitParam+" Areas\n"+line)
	if previousAreasURL == nil {
		err := errors.New("Error: cannot map back from start")
		return err
	}
	return printAreas(previousAreasURL)
}

func printAreas(URL *string) error {
	areas, err := pokeapi.GetLocationAreas(URL)
	if err != nil {
		return err
	}
	nextAreasURL = areas.Next
	previousAreasURL = areas.Previous
	for _, area := range areas.Results {
		fmt.Println(strings.Title(strings.ReplaceAll(area.Name, "-", " ")))
	}
	return nil
}

func commandExplore() error {
	if secondInput == "" {
		err := errors.New("Error: must include an area to explore.\n(Enter 'explore area')")
		return err
	}
	areaInput := strings.ReplaceAll(secondInput, " ", "-")
	areaURL := pokeapiAreaEndpoint + areaInput
	exploredArea, err := pokeapi.ExploreArea(&areaURL)
	if err != nil {
		return err
	}
	areaName := fmt.Sprint(strings.Title(strings.ReplaceAll(secondInput, "-", " ")))
	line := strings.Repeat("-", len(areaName) + 13)
	fmt.Println(line+"\nExploring "+areaName+"...\n"+line+"\nFound Pokemon:")
	for _, pokemon := range exploredArea.PokemonEncounters {
		fmt.Println(" > "+pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch() error {
	if secondInput == "" {
		err := errors.New("Error: must include a pokemon to catch.\n(Enter 'catch pokemon')")
		return err
	}
	pokemonName := secondInput
	pokemonURL := pokeapiPokemonEndpoint + pokemonName
	_, err := pokeapi.GetPokemon(&pokemonURL)
	if err != nil {
		return err
	}

	line := strings.Repeat("-", len(pokemonName) + 26)
	fmt.Println(line+"\nThrowing a Pokeball at "+pokemonName+"...\n"+line)
	fmt.Println("\n"+pokemonName+" escaped!")
	return nil
}

/*
Give the user a chance to catch the Pokemon 
using the math/rand package.
You can use the pokemon's "base experience" 
to determine the chance of catching it. 
	pokemon.BaseExperience int 
The higher the base experience, the harder it should be to catch.

Once the Pokemon is caught, add it to the user's Pokedex. 
I used a map[string]Pokemon to keep track of caught Pokemon.
You'll want to store the Pokemon's data 
so that in the next step we can use it.

Example usage
Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu escaped!
Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!
*/

func commandHelp() error {
	line := "----------------"
	fmt.Println(line+"\nPokedex commands\n"+line)
	for _, command := range commands {
		fmt.Println(command.name+": "+command.description)
	}
	return nil
}

func commandExit() error {
	if secondInput == "y" || secondInput == "yes" {
		play = false
		fmt.Println("\nGoodbye!")
		return nil
	}
	fmt.Print("Exit? (y/n) ")
	go getInput()
	input := <- ch
	input, _ = parse(input)
	if input == "y" || input == "yes" {
		play = false
		fmt.Println("\nGoodbye!")
	}
	return nil
}

func getInput() error {
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func parse(toParse string) (string, string) {
	toParse = strings.TrimSpace(toParse)
	toParse = strings.ToLower(toParse)
	paramSlice := strings.SplitN(toParse, " ", 2)
	first := paramSlice[0]
	second := ""
	if len(paramSlice) == 2 {
		second = paramSlice[1]
	}
	return first, second
}

var pokeapiURL string = "https://pokeapi.co/api/v2/"

var pokeapiAreaEndpoint string = fmt.Sprint(pokeapiURL + "location-area/")
var limitParam string = "10"
var locationAreasURL string = fmt.Sprint(pokeapiAreaEndpoint+"?offset=0&limit="+limitParam)
var nextAreasURL *string = &locationAreasURL
var previousAreasURL *string = nil

var pokeapiPokemonEndpoint string = fmt.Sprint(pokeapiURL + "pokemon/")

var play bool = true
var commands map[string]cliCommand
var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
var ch = make(chan string)
var secondInput string = "";

func main() {
	commands = newCommandsMap()
	line := "\n----------------------\n"
	fmt.Print(line+"Welcome to PokedexCLI!"+line+"\n")
	for ; play == true; {
		fmt.Print("pokedex > ")
		go getInput()
		inputString := <- ch
		inputString, secondInput = parse(inputString)
		command, ok := commands[inputString]
		fmt.Print("\n")
		if !ok {
			fmt.Println("Invalid input, please try again.\n(Enter 'help' to list available commands)\n(Enter 'q' to quit)")
		} else {
			err := command.callback()
			if err != nil {
				log.Println(err)
			}
		}
		fmt.Print("\n")
	}
}
