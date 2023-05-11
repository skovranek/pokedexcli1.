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
	"log"
)

/* 
TODO:
[x] list pokemon abilities
[ ] verify if this should be outside func in pokeAPI package: NO
var locationAreas *LocationAreas = &LocationAreas{}
^No!
[ ] add arrows in input
[-] better variable names
[ ] remove/comment out print statements
[ ] add comments
*/

type Pokemon struct {
	pokeapi.Pokemon
	Caught bool `json:"caught"`
}

var pokedex map[string]Pokemon = map[string]Pokemon{}

var pokeapiURL string = "https://pokeapi.co/api/v2/"
var areaEndpoint string = fmt.Sprint(pokeapiURL + "location-area/")
var limitParam string = "10"
var areasURL string = fmt.Sprint(areaEndpoint+"?offset=0&limit="+limitParam)
var nextAreasURL *string = &areasURL
var previousAreasURL *string = nil
var pokemonEndpoint string = fmt.Sprint(pokeapiURL + "pokemon/")

var difficulty int = 100 // easy: 0, medium: 100, hard: 200

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
	return getAreas(nextAreasURL)
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
	return getAreas(previousAreasURL)
}

func getAreas(URL *string) error {
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
	areaURL := areaEndpoint + areaInput
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
	name := secondInput
	nameCap := strings.Title(name)
	line := strings.Repeat("-", len(name) + 26)
	fmt.Println(line+"\nThrowing a pokeball at "+nameCap+"...\n"+line)

	if _, ok := pokedex[name]; !ok {
		URL := pokemonEndpoint + name
		newPokemon, err := pokeapi.GetPokemon(&URL)
		if err != nil {
			msg := "(Must include the correct name of an pokemon to catch. Enter 'catch pokemon')"
			err := errors.New(invalidMsg+msg)
			return err
		}
		pokedex[name] = Pokemon{
			*newPokemon, //Pokemon struct field:values
			false, //Caught
		}
	}

	pokemon := pokedex[name]
	exp := math.Pow(math.Log(float64(pokemon.BaseExperience + difficulty)), 2)
	throw := float64(rand.Intn(50))
	fmt.Println("throw: ",throw,", exp: ",exp)

	if throw > exp {
		fmt.Println("Gotcha! "+nameCap+" was caught!")
		if !pokemon.Caught {
			pokemon.Caught = true
			pokedex[name] = pokemon
			fmt.Println(nameCap+"'s data was newly added to the Pokedex!")
			fmt.Println("pokedex > inspect "+name+"")
			commandInspect()
		} else {
			fmt.Println(""+nameCap+"'s data is already in the Pokedex.")
		}
	} else {
		fmt.Println(nameCap+" escaped!")
	}
	return nil
}

func getIndex(order int) string {
	index := fmt.Sprint(order)
	if order < 10 {
		index = "00" + index
	} else if order < 100 {
		index = "0" + index
	}
	return index
}

func commandInspect() error {
	if secondInput == "" {
		err := errors.New(invalidMsg+"(Must include a pokemon to inspect. Enter 'inspect pokemon')")
		return err
	}
	name := secondInput
	if pokemon, ok := pokedex[name]; ok && pokemon.Caught {
		pokedexIndex := getIndex(pokemon.Order)
		line := strings.Repeat("-", len(name) + 16)+"\n"
		fmt.Print(line+"Pokedex: No."+pokedexIndex+" "+strings.Title(name)+"\n"+line)

		fmt.Println("Height:",pokemon.Height)
		fmt.Println("Weight:",pokemon.Weight)

		fmt.Println("Statistics:")
		for _, s := range pokemon.Stats {
			statName := strings.Title(strings.ReplaceAll(s.Stat.Name, "-", " "))
			fmt.Println(" > "+statName+":",s.BaseStat)
		}

		t := "Type: " + strings.Title(pokemon.Types[0].Type.Name)
		if len(pokemon.Types) > 1 {
			t = "Dual-" + t + "/" + strings.Title(pokemon.Types[1].Type.Name)
		}
		fmt.Println(t)

		fmt.Println("Abilities:")
		for _, a := range pokemon.Abilities {
			fmt.Println(" > "+strings.Title(strings.ReplaceAll(a.Ability.Name, "-", " ")))
		}
	} else {
		msg := "(Must catch a pokemon before its data can be added to the Pokedex)"
		err := errors.New("Error: "+name+" is not in the Pokedex.\n"+msg)
		return err
	}
	return nil
}

func commandPokedex() error {
	if secondInput != "" {
		err := errors.New(invalidMsg+"(Command 'pokedex' does not take an input. Enter only 'pokedex')")
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
		fmt.Println("> No."+index+" "+strings.Title(pokemon.Name))
	}
	return nil
}

func commandHelp() error {
	if secondInput != "" {
		err := errors.New(invalidMsg+"(Command 'help' does not take an input. Enter only 'help')")
		return err
	}
	line := "----------------"
	fmt.Println(line+"\nPokedex Commands\n"+line)
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
		input := <- ch
		input, secondInput = parse(input)
		command, ok := commands[input]
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