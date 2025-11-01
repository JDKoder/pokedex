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
	callback func(conf *config) error
}

type config struct {
	next string
	previous string
}

type Locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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
		"map": {
			name: "map",
			description: "Displays the names of 20 Locations areas in the Pokemon world.  Each subsequent call to map should display the next 20 locations.",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Displays the names of last 20 Locations areas in the Pokemon world.  Each subsequent call to map should display the previous 20 locations.",
			callback: commandMapBack,
		},
	}
}

func main() {
	ioscanner := bufio.NewScanner(os.Stdin)
	conf := config{next: "", previous: ""}
	for true {
		fmt.Print("Pokedex > ")
		ioscanner.Scan()
		input_word := cleanInput(ioscanner.Text())
		if len(input_word) < 1 {
			continue
		}
		command := input_word[0]
		if val, ok := getCommands()[command]; ok {
			err := val.callback(&conf)
			if err != nil {
				fmt.Printf("Command [%s] failed: %s", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
		
	}
}

func commandHelp(conf *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, val := range getCommands() {
		fmt.Printf("%s: %s\n",  val.name, val.description)
	}
	return nil
}

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Call to os.Exit failed")

}

/**
pass config by reference because each call to map should update the config's 
next url and previous url.  We first check if next is nil, if so, we use the default
url
**/
func commandMap(conf *config) error {
	if conf.next == "" {
		conf.next = "https://pokeapi.co/api/v2/location/"
	}
	var locResults Locations
	err := makeGetRequest(conf.next, &locResults)
	if err != nil {
		return err
	}
	fmt.Printf("locResults: %d\n", locResults.Count)
	for i := 0; i< len(locResults.Results); i++ {
		fmt.Println(locResults.Results[i].Name)
	}
	if str, ok := locResults.Previous.(string); ok {
		conf.previous = str
	} else {
		conf.previous = ""
	}
	conf.next = locResults.Next
	return nil
}

/**
pass config by reference because each call to map should update the config's 
next url and previous url.  We first check if next is nil, if so, we use the default
url
**/
func commandMapBack(conf *config) error {
	if conf.previous == "" {
		conf.previous = "https://pokeapi.co/api/v2/location/"
	}
	var locResults Locations
	err := makeGetRequest(conf.previous, &locResults)
	if err != nil {
		return err
	}
	fmt.Printf("locResults: %d\n", locResults.Count)
	for i := 0; i< len(locResults.Results); i++ {
		fmt.Println(locResults.Results[i].Name)
	}
	if str, ok := locResults.Previous.(string); ok {
		conf.previous = str
	} else {
		conf.previous = ""
	}
	conf.next = locResults.Next
	return nil
}

func cleanInput(text string) []string {
	slices := strings.Fields(strings.ToLower(text))
	return slices
}
