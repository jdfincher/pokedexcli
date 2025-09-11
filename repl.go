package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jdfincher/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Client, string) (*pokeapi.Client, error)
}

func getCommands() (commands map[string]cliCommand) {
	return map[string]cliCommand{
		"help":    {name: "help", description: "Displays a help message", callback: commandHelp},
		"exit":    {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"map":     {name: "map", description: "Shows the NEXT 20 locations", callback: commandMap},
		"mapb":    {name: "mapb", description: "Shows the PREVIOUS 20 locations", callback: commandMapBack},
		"explore": {name: "explore", description: "Shows Pokemon in a specific area", callback: commandExplore},
	}
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	trim := strings.TrimSpace(lower)
	clean := strings.Split(trim, " ")
	return clean
}

func commandExit(client *pokeapi.Client, arg string) (*pokeapi.Client, error) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return client, nil
}

func commandHelp(client *pokeapi.Client, arg string) (*pokeapi.Client, error) {
	fmt.Print(`
░█░█░█▀▀░█▀█░█▀▀░█▀▀
░█░█░▀▀█░█▀█░█░█░█▀▀
░▀▀▀░▀▀▀░▀░▀░▀▀▀░▀▀▀

`)
	fmt.Println(arg)
	commands := getCommands()
	for _, v := range commands {
		fmt.Println(v.name + ": " + v.description)
	}
	fmt.Print("\n")
	return client, nil
}

func commandMap(client *pokeapi.Client, arg string) (*pokeapi.Client, error) {
	var err error
	if client.Loc.Next == "" {
		client, err = client.Get(client.BaseURL + "/location-area/")
		if err != nil {
			return client, err
		}
	} else {
		client, err = client.Get(client.Loc.Next)
		if err != nil {
			return client, err
		}
	}
	fmt.Println("********************")
	for i := 0; i < len(client.Loc.Results); i++ {
		fmt.Println(client.Loc.Results[i].Name)
	}
	fmt.Println("********************")
	return client, nil
}

func commandMapBack(client *pokeapi.Client, arg string) (*pokeapi.Client, error) {
	var err error
	if client.Loc == nil {
		client, err = client.Get(client.BaseURL + "/location-area/")
		if err != nil {
			return client, err
		}
		fmt.Println("Already at the beginning of list, use 'map' command to advance")
		return client, nil
	} else if client.Loc.Previous == "" {
		fmt.Println("Already at the beginning of list, use 'map' command to advance")
		return client, nil
	} else {
		client, err = client.Get(client.Loc.Previous)
		if err != nil {
			return client, err
		}
		fmt.Println("********************")
		for i := 0; i < len(client.Loc.Results); i++ {
			fmt.Println(client.Loc.Results[i].Name)
		}
		fmt.Println("********************")
	}
	return client, nil
}

func commandExplore(client *pokeapi.Client, arg string) (*pokeapi.Client, error) {
	var err error
	if arg == "" {
		fmt.Println("Location name must be provided to explore area")
		return client, nil
	}
	client, err = client.GetPok(client.BaseURL + "/location-area/" + arg)
	if err != nil {
		return client, err
	}
	fmt.Println("********************")
	fmt.Printf("Exploring %v...\n", arg)
	for i := 0; i < len(client.Pok.PokemonEncounters); i++ {
		fmt.Printf(" - %v\n", client.Pok.PokemonEncounters[i].Pokemon.Name)
	}
	fmt.Println("********************")
	return client, nil
}

func repLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	client := pokeapi.NewClient(5 * time.Minute)
	var err error
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		word := scanner.Text()
		cleaned := cleanInput(word)
		com, ok := commands[cleaned[0]]
		if !ok {
			fmt.Println("Unknown Command")
			continue
		}
		var arg string
		if len(cleaned) > 1 {
			arg = cleaned[1]
			client, err = com.callback(client, arg)
			if err != nil {
				fmt.Print(err)
				continue
			}
		} else if len(cleaned) < 2 {
			arg = string("")
			client, err = com.callback(client, arg)
			if err != nil {
				fmt.Print(err)
			}
		}
	}
}
