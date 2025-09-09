package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jdfincher/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *pokeapi.Config) (*pokeapi.Config, error)
}

func getCommands() (commands map[string]cliCommand) {
	return map[string]cliCommand{
		"help": {name: "help", description: "Displays a help message", callback: commandHelp},
		"exit": {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"map":  {name: "map", description: "Shows the NEXT 20 locations", callback: commandMap},
		"mapb": {name: "mapb", description: "Shows the PREVIOUS 20 locations", callback: commandMapBack},
	}
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	trim := strings.TrimSpace(lower)
	clean := strings.Split(trim, " ")
	return clean
}

func commandExit(c *pokeapi.Config) (*pokeapi.Config, error) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return c, nil
}

func commandHelp(c *pokeapi.Config) (*pokeapi.Config, error) {
	fmt.Print(`
░█░█░█▀▀░█▀█░█▀▀░█▀▀
░█░█░▀▀█░█▀█░█░█░█▀▀
░▀▀▀░▀▀▀░▀░▀░▀▀▀░▀▀▀

`)
	commands := getCommands()
	for _, v := range commands {
		fmt.Println(v.name + ": " + v.description)
	}
	fmt.Printf("\n")
	return c, nil
}

func commandMap(c *pokeapi.Config) (*pokeapi.Config, error) {
	if c == nil {
		config, err := pokeapi.GetConfig("https://pokeapi.co/api/v2/location-area/")
		if err != nil {
			return config, err
		}
		fmt.Println("********************")
		for i := 0; i < len(config.Results); i++ {
			fmt.Println(config.Results[i].Name)
		}
		fmt.Println("********************")
		return config, nil
	} else {
		config, err := pokeapi.GetConfig(c.Next)
		if err != nil {
			return config, err
		}
		fmt.Println("********************")
		for i := 0; i < len(config.Results); i++ {
			fmt.Println(config.Results[i].Name)
		}
		fmt.Println("********************")
		return config, nil
	}
}

func commandMapBack(c *pokeapi.Config) (*pokeapi.Config, error) {
	var err error
	if c == nil {
		fmt.Println("Use command 'map' to advance")
		return c, fmt.Errorf("error: cannot go backwards, already at beggining")
	}
	if c.Previous != "" {
		c, err = pokeapi.GetConfig(c.Previous)
	} else {
		fmt.Println("Use Command 'map to advance")
		return c, fmt.Errorf("error: cannot go backwards, already at beggining")
	}
	if err != nil {
		return c, fmt.Errorf("error retrieving previous: %w", err)
	}
	fmt.Println("********************")
	for i := 0; i < len(c.Results); i++ {
		fmt.Println(c.Results[i].Name)
	}
	fmt.Println("********************")
	return c, nil
}

func repLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	var config *pokeapi.Config
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
		config, err = com.callback(config)
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
	}
}
