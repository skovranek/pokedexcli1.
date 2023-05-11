package main

import (
	"fmt"
	"internal/pokeapi"
	"bufio"
	"os"
	"strings"
	"errors"
	"math"
	"math/rand"
	"log"
)

var pokeapiURL string = "https://pokeapi.co/api/v2/"

var pokeapiAreaEndpoint string = fmt.Sprint(pokeapiURL + "location-area/")
var limitParam string = "10"
var locationAreasURL string = fmt.Sprint(pokeapiAreaEndpoint+"?offset=0&limit="+limitParam)
var nextAreasURL *string = &locationAreasURL
var previousAreasURL *string = nil

var pokeapiPokemonEndpoint string = fmt.Sprint(pokeapiURL + "pokemon/")

type Pokemon struct {
	pokeapi.Pokemon
	Caught bool `json:"caught"`
}

var pokedex map[string]Pokemon = map[string]Pokemon{}

var play bool = true
var commands map[string]cliCommand
var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
var ch = make(chan string)
var secondInput string = "";
var invalidMsg string = "Error: invalid input, please try again.\n"

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
	fmt.Println(">> 1) "+first+", 2) "+second)
	return first, second
}

func commandMap() error {
	if secondInput != "" {
		err := errors.New(invalidMsg+"(Command 'map' does not take an input. Enter only 'map')")
		return err
	}
	line := "-----------------"
	fmt.Println(line+"\nThe Next "+limitParam+" Areas\n"+line)
	if nextAreasURL == nil {
		err := errors.New("Error: cannot go past the end of the map.")
		return err
	}
	return printAreas(nextAreasURL)
}

func commandMapBack() error {
	if secondInput != "" {
		err := errors.New(invalidMsg+"(Command 'mapb' does not take an input. Enter only 'mapb')")
		return err
	}
	line := "---------------------"
	fmt.Println(line+"\nThe Previous "+limitParam+" Areas\n"+line)
	if previousAreasURL == nil {
		err := errors.New("Error: cannot go back before the start of the map.")
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
		err := errors.New("(Must include an area to explore. Enter 'explore area')")
		return err
	}
	areaName := fmt.Sprint(strings.Title(strings.ReplaceAll(secondInput, "-", " ")))
	line := strings.Repeat("-", len(areaName) + 13)
	fmt.Println(line+"\nExploring "+areaName+"...\n"+line)

	areaInput := strings.ReplaceAll(secondInput, " ", "-")
	areaURL := pokeapiAreaEndpoint + areaInput
	exploredArea, err := pokeapi.ExploreArea(&areaURL)
	if err != nil {
		msg := "(Must include the correct name of an area to explore. Enter 'explore area')"
		err := errors.New(invalidMsg+msg)
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range exploredArea.PokemonEncounters {
		fmt.Println(" > "+pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch() error {
	if secondInput == "" {
		err := errors.New(invalidMsg+"(Must include a pokemon to catch. Enter 'catch pokemon')")
		return err
	}
	pokemonName := secondInput
	pokemonCap := strings.Title(pokemonName)
	line := strings.Repeat("-", len(pokemonName) + 26)
	fmt.Println(line+"\nThrowing a Pokeball at "+pokemonCap+"...\n"+line)

	if _, ok := pokedex[pokemonName]; !ok {
		pokemonURL := pokeapiPokemonEndpoint + pokemonName
		newPokemon, err := pokeapi.GetPokemon(&pokemonURL)
		if err != nil {
			msg := "(Must include the correct name of an pokemon to catch. Enter 'catch pokemon')"
			err := errors.New(invalidMsg+msg)
			return err
		}
		pokedex[pokemonName] = Pokemon{
			*newPokemon, //Pokemon struct field:values
			false, //Caught
		}
	}

	pokemon := pokedex[pokemonName]
	exp := math.Pow(math.Log(float64(pokemon.BaseExperience+200)), 2)
	throw := float64(rand.Intn(50))
	fmt.Println("throw: ",throw,", exp: ",exp)

	if throw > exp {
		fmt.Println("\nGotcha! "+pokemonCap+" was caught!")
		if !pokemon.Caught {
			pokemon.Caught = true
			pokedex[pokemonName] = pokemon
			fmt.Println("\n"+pokemonCap+"'s data was newly added to the Pokedex!")

			// TODO: print inspection
		} else {
			fmt.Println("\n"+pokemonCap+" is already in the Pokedex.")
		}
	} else {
		fmt.Println("\n"+pokemonCap+" escaped!")
	}
	return nil
}

func getOrder(pokemonInt int) string {
	orderString := fmt.Sprint(pokemonInt)
	if pokemonInt < 10 {
		orderString = "00" + orderString
	} else if pokemonInt < 100 {
		orderString = "0" + orderString
	}
	return orderString
}

func commandInspect() error {
	if secondInput == "" {
		err := errors.New(invalidMsg+"(Must include a pokemon to inspect. Enter 'inspect pokemon')")
		return err
	}
	pokemonName := secondInput
	if pokemon, ok := pokedex[pokemonName]; ok && pokemon.Caught {
		pokemonOrder := getOrder(pokemon.Order)
		line := strings.Repeat("-", len(pokemonName) + 16)+"\n"
		fmt.Print(line+"|Pokedex: No."+pokemonOrder+" "+strings.Title(pokemonName)+"|\n"+line)
		fmt.Println("Name:",pokemonName,"\nHeight:",pokemon.Height,"\nWeight:",pokemon.Weight,"\nWeight:",pokemon.Weight,"\nStats:")
		for _, stat := range pokemon.Stats {
			statName := strings.Title(strings.ReplaceAll(stat.Stat.Name, "-", " "))
			fmt.Println(" > "+statName+":",stat.BaseStat)
		}
		fmt.Println("Type:")
		for _, Type := range pokemon.Types {
			fmt.Println(" > "+strings.Title(Type.Type.Name))
		}
	} else {
		err := errors.New("Error: "+pokemonName+" is not in the Pokedex.\n(Must catch a pokemon before its data can be added to the Pokedex)")
		return err
	}
	return nil
}

func commandHelp() error {
	if secondInput != "" {
		err := errors.New(invalidMsg+"(Command 'help' does not take an input. Enter only 'help')")
		return err
	}
	line := "----------------"
	fmt.Println(line+"\nPokedex commands\n"+line)
	for _, command := range commands {
		fmt.Println("> "+command.name+" - "+command.description)
	}
	return nil
}

func commandExit() error {
	if secondInput == "y" || secondInput == "yes" {
		play = false
		fmt.Println("Goodbye!")
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

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

func newCommandsMap() map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:		 "map",
			description: fmt.Sprint("Displays the next "+limitParam+" areas."),
			callback:	 commandMap,
		},
		"mapb": {
			name:		 "mapb",
			description: fmt.Sprint("Displays back the previous "+limitParam+" areas."),
			callback:	 commandMapBack,
		},
		"explore": {
			name:		 "explore <area>",
			description: "Explore an area to search for pokemon.",
			callback:	 commandExplore,
		},
		"catch": {
			name: 		 "catch <pokemon>",
			description: "Catch a pokemon and add it to the Pokedex.",
			callback:	 commandCatch,
		},
		"inspect": {
			name:		 "inspect <pokemon>",
			description: "Check the Pokedex for data about a pokemon.",
			callback:	 commandInspect,
		},
		"help": {
			name:		 "help",
			description: "Lists available commands.",
			callback:	 commandHelp,
		},
		"exit": {
			name:		 "exit",
			description: "Exit the Pokedex.",
			callback:	 commandExit,
		},
		"q": {
			name:		 "q",
			description: "Also exit.",
			callback: func() error {
				fmt.Println("Goodbye!")
				play = false
				return nil
			},
		},
	}
}

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
			fmt.Println(invalidMsg+"(Enter 'help' to list available commands)")
		} else {
			err := command.callback()
			if err != nil {
				log.Println(err)
			}
		}
		fmt.Print("\n")
	}
}