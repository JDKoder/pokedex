package main
import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand {
		"help": {
			name: "help",
			description: "Show the help for pokedex",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
	}
}

func main() {
	ioscanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex > ")
		ioscanner.Scan()
		input_word := cleanInput(ioscanner.Text())
		if len(input_word) < 1 {
			continue
		}
		command := input_word[0]
		if val, ok := getCommands()[command]; ok {
			err := val.callback()
			if err != nil {
				fmt.Printf("Command [%s] failed: %s", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
		
	}
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, val := range getCommands() {
		fmt.Printf("%s: %s\n",  val.name, val.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Call to os.Exit failed")

}

func cleanInput(text string) []string {
	slices := strings.Fields(strings.ToLower(text))
	//fmt.Println(slices)
	return slices
	
}
