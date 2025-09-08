package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() (commands map[string]cliCommand) {
	return map[string]cliCommand{
		"help": {name: "help", description: "Displays a help message", callback: commandHelp},
		"exit": {name: "exit", description: "Exit the Pokedex", callback: commandExit},
	}
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	trim := strings.TrimSpace(lower)
	clean := strings.Split(trim, " ")
	return clean
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	commands := getCommands()
	for _, v := range commands {
		fmt.Println(v.name + ": " + v.description)
	}
	fmt.Printf("\n")
	return nil
}

func repLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		word := scanner.Text()
		cleaned := cleanInput(word)
		com, ok := commands[cleaned[0]]
		if !ok {
			fmt.Println("Unknown Command")
			continue
		}
		err := com.callback()
		if err != nil {
			fmt.Printf("error:%v", err)
		}
	}
}
