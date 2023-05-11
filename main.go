package main

import (
	"fmt"
	"sort"
	"internal/pokeapi"
	"bufio"
	"os"
	"strings"
	"errors"
	"math"
	"math/rand"
)

/* 
TODO:
[x] list pokemon abilities
[x] verify if this should be outside func in pokeAPI package: NO
var locationAreas *LocationAreas = &LocationAreas{}
^No!
[ ] use f strings
[ ] split up code into files
[ ] better variable names
[ ] rename secondInput
[ ] add arrows in input
[ ] remove/comment out print statements
[ ] add comments
*/

type Pokemon struct {
	pokeapi.Pokemon
	Caught bool `json:"caught"`
}

var pokedex map[string]Pokemon = map[string]Pokemon{}

var pokeapiURL string = "https://pokeapi.co/api/v2/"
var areaEndpoint string = pokeapiURL+"location-area/"
var limitParam string = "10"
var areasURL string = areaEndpoint+"?offset=0&limit="+limitParam
var nextAreasURL *string = &areasURL
var previousAreasURL *string = nil
var pokemonEndpoint string = pokeapiURL+"pokemon/"

var difficulty int = 0 // easy: 0, medium: 100, hard: 200

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
	args := strings.SplitN(toParse, " ", 2)

	first := args[0]
	second := ""
	if len(args) == 2 {
		second = args[1]
	}
	fmt.Printf("******msg: 1) %s, 2) %s\n", first, second)

	return first, second
}

func commandMap() error {
	if secondInput != "" {
		msg := "(Command 'map' does not take an input. Enter only 'map')"
		err := errors.New(invalidMsg+msg)
		return err
	}

	line := "-----------------"
	fmt.Printf("%s\nThe Next %s Areas\n%s\n", line, limitParam, line)

	if nextAreasURL == nil {
		err := errors.New("Error: cannot go past the end of the map.")
		return err
	}
	return getAreas(*nextAreasURL)
}

func commandMapBack() error {
	if secondInput != "" {
		msg := "(Command 'mapb' does not take an input. Enter only 'mapb')"
		err := errors.New(invalidMsg+msg)
		return err
	}

	line := "---------------------"
	fmt.Printf("%s\nThe Previous %s Areas\n%s\n", line, limitParam, line)

	if previousAreasURL == nil {
		err := errors.New("Error: cannot go back before the start of the map.")
		return err
	}

	return getAreas(*previousAreasURL)
}

func getAreas(URL string) error {
	areas, err := pokeapi.GetLocationAreas(URL)
	if err != nil {
		return err
	}

	nextAreasURL = areas.Next
	previousAreasURL = areas.Previous

	for _, a := range areas.Results {
		area := strings.ReplaceAll(a.Name, "-", " ")
		area = strings.Title(area)
		fmt.Println(area)
	}
	return nil
}

func commandExplore() error {
	if secondInput == "" {
		err := errors.New("(Must include an area to explore. Enter 'explore area')")
		return err
	}

	areaName := strings.ReplaceAll(secondInput, "-", " ")
	areaName = strings.Title(areaName)

	line := strings.Repeat("-", len(areaName) + 13)
	fmt.Printf("%s\nExploring %s...\n%s\n", line, areaName, line)

	areaInput := strings.ReplaceAll(secondInput, " ", "-")
	areaURL := areaEndpoint + areaInput

	exploredArea, err := pokeapi.ExploreArea(areaURL)
	if err != nil {
		msg := "(Must include the correct name of an area to explore. Enter 'explore area')"
		err := errors.New(invalidMsg+msg)
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range exploredArea.PokemonEncounters {
		fmt.Printf(" > %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch() error {
	if secondInput == "" {
		msg := "(Must include a pokemon to catch. Enter 'catch pokemon')"
		err := errors.New(invalidMsg+msg)
		return err
	}

	name := secondInput
	nameCap := strings.Title(name)
	line := strings.Repeat("-", len(name) + 26)
	fmt.Printf("%s\nThrowing a pokeball at %s...\n%s\n", line, nameCap, line)

	if _, ok := pokedex[name]; !ok {
		URL := pokemonEndpoint + name
		newPokemon, err := pokeapi.GetPokemon(URL)
		if err != nil {
			msg := "(Must include the correct name of an pokemon to catch. Enter 'catch pokemon')"
			err := errors.New(invalidMsg+msg)
			return err
		}
		pokedex[name] = Pokemon{
			newPokemon, //Pokemon struct field:values
			false,      //Caught bool
		}
	}
	pokemon := pokedex[name]

	exp := float64(pokemon.BaseExperience + difficulty) // diff > easy: 0, medium: 100, hard: 200
	exp = math.Log(exp)
	exp = math.Pow(exp, 2)
	
	throw := float64(rand.Intn(50))
	fmt.Printf("******msg: throw: %v, exp: %v\n", throw, exp)

	if throw > exp {
		fmt.Printf("Gotcha! %s was caught!\n", nameCap)
		if !pokemon.Caught {
			pokemon.Caught = true
			pokedex[name] = pokemon

			fmt.Printf("%s's data was newly added to the Pokedex!\n", nameCap)
			fmt.Printf("pokedex > inspect %s\n", name)
			commandInspect()
		} else {
			fmt.Printf("%s's data is already in the Pokedex.\n", nameCap)
		}
	} else {
		fmt.Printf("%s escaped!\n", nameCap)
	}
	return nil
}

func getIndex(order int) string {
	index := fmt.Sprint(order)
	zeros := ""

	if order < 10 {
		zeros = "00"
	} else if order < 100 {
		zeros = "0"
	}
	return (zeros + index)
}

func commandInspect() error {
	if secondInput == "" {
		msg := "(Must include a pokemon to inspect. Enter 'inspect pokemon')"
		err := errors.New(invalidMsg+msg)
		return err
	}
	name := secondInput

	if pokemon, ok := pokedex[name]; ok && pokemon.Caught {
		pokedexIndex := getIndex(pokemon.Order)
		line := strings.Repeat("-", len(name) + 16)
		fmt.Printf("%s\nPokedex: No.%s %s\n%s\n", line, pokedexIndex, strings.Title(name), line)

		fmt.Printf("Height: %s\n", pokemon.Height)
		fmt.Printf("Weight: %s\n", pokemon.Weight)

		fmt.Println("Stats:")
		for _, s := range pokemon.Stats {
			stat := strings.ReplaceAll(s.Stat.Name, "-", " ")
			stat = strings.Title(stat)
			fmt.Printf(" > %s: %v\n", stat, s.BaseStat)
		}

		type1 := strings.Title(pokemon.Types[0].Type.Name)
		dual := ""
		type2 := ""
		if len(pokemon.Types) > 1 {
			dual = "Dual-"
			type2 = "/" + strings.Title(pokemon.Types[1].Type.Name)
		}
		fmt.Printf("%sType: %s%s\n", dual, type1, type2)

		fmt.Println("Abilities:")
		for _, a := range pokemon.Abilities {
			ability := strings.ReplaceAll(a.Ability.Name, "-", " ")
			ability = strings.Title(ability)
			fmt.Printf(" > %s\n", ability)
		}
	} else {
		msg := fmt.Sprintf("Error: %s is not in the Pokedex.", name)
		msg = msg+"\n(Must catch a pokemon before its data can be added to the Pokedex)"
		err := errors.New(msg)
		return err
	}
	return nil
}

func commandPokedex() error {
	if secondInput != "" {
		msg := "(Command 'pokedex' does not take an input. Enter only 'pokedex')"
		err := errors.New(invalidMsg+msg)
		return err
	}

	fmt.Println("------------\nPokedex List\n------------")

	list := []Pokemon{}
	for _, pokemon := range pokedex {
		if pokemon.Caught {
			list = append(list, pokemon)
		}
	}

	if len(list) == 0 {
		fmt.Println(" > ...\n(Catch a pokemon to add its data to the Pokedex)")
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Order < list[j].Order
	})

	for _, pokemon := range list {
		index := getIndex(pokemon.Order)
		fmt.Printf("> No.%s %s\n", index, strings.Title(pokemon.Name))
	}
	return nil
}

func commandHelp() error {
	if secondInput != "" {
		msg := "(Command 'help' does not take an input. Enter only 'help')"
		err := errors.New(invalidMsg+msg)
		return err
	}

	fmt.Println("----------------\nPokedex Commands\n----------------")

	for _, command := range commands {
		fmt.Printf("> %s - %s\n", command.name, command.description)
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
		"pokedex": {
			name:		 "pokedex",
			description: "Display the list of pokemon in the Pokedex.",
			callback:	 commandPokedex,
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
		/*
		"q": {
			name:		 "q",
			description: "Quit.",
			callback: func() error {
				fmt.Println("Goodbye!")
				play = false
				return nil
			},
		},
		*/
	}
}

func main() {
	line := "----------------------"
	fmt.Printf("\n%s\nWelcome to PokedexCLI!\n%s\n\n", line, line)

	commands = newCommandsMap()

	for ; play == true; {
		fmt.Print("pokedex > ")

		go getInput()
		input := <- ch
		input, secondInput = parse(input)

		fmt.Println("")

		command, ok := commands[input]

		if !ok {
			msg := "(Enter 'help' to list available commands)"
			fmt.Println(invalidMsg+msg)
		} else {

			err := command.callback()

			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println("")
	}
}