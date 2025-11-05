package main
import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"math"
	"math/rand"
)

type cliCommand struct {
	name string
	description string
	callback func(conf *config) error
}

type config struct {
	next string
	previous string
	args []string
}

type Locations struct {
	Count    int		`json:"count"`
	Next     *string 	`json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string	`json:"name"`
		URL  string	`json:"url"`
	} 					`json:"results"`
}

const (
	mapURL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	areaURL = "https://pokeapi.co/api/v2/location-area/%s"
	pokemonURL = "https://pokeapi.co/api/v2/pokemon/%s"
)

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
			description: "Displays the names of last 20 Locations areas in the Pokemon world.  Each subsequent call to mapb should display the previous 20 locations.",
			callback: commandMapBack,
		},
		"explore": {
			name: "explore",
			description: "takes a area name (see help for commands: map, mapb ) as an argument and returns a list of the pokemon in that area",
			callback: commandExploreArea,
		},
		"catch": {
			name: "catch",
			description: "Attempt to catch a pokemon given by name in the first argument.  A message is returned whether or not the attempt to catch the pokemon was succesful or not.  If it was, the pokemon is added to the pokedex map",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "Given a captured pokemon as the first argument, prints information about that pokemon.",
			callback: commandInspect,
		},
	}
}

var (
	pokedex = make(map[string]Pokemon)
)

func main() {
	ioscanner := bufio.NewScanner(os.Stdin)
	conf := config{next: "", previous: ""}
	for true {
		fmt.Print("Pokedex > ")
		ioscanner.Scan()
		input_word := cleanInput(ioscanner.Text())
		var command string
		var args []string
		conf.args = []string{}
		if len(input_word) < 1 {
			continue
		} else {
			command = input_word[0]
			for i, _ := range input_word {
				if i == 0 {
					continue
				}
				args = append(args, input_word[i])
			}
			//fmt.Println("args is %w", args)
			conf.args = args
		}
		
		if val, ok := getCommands()[command]; ok {
			err := val.callback(&conf)
			if err != nil {
				fmt.Errorf("Command [%s] failed: %w", command, err)
			}
		} else {
			fmt.Println("Unknown command")
		}
		
	}
}

func commandHelp(conf *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
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
		conf.next = mapURL
	}
	locations, err := getLocations(conf.next)
	stringReferenceAssignment(&conf.next, locations.Next)
	stringReferenceAssignment(&conf.previous, locations.Previous)
	return err
}

/**
pass config by reference because each call to map should update the config's 
next url and previous url.  We first check if next is nil, if so, we use the default
url
**/
func commandMapBack(conf *config) error {
	if conf.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	locations, err := getLocations(conf.previous)
	stringReferenceAssignment(&conf.next, locations.Next)
	stringReferenceAssignment(&conf.previous, locations.Previous)
	return err
}

func commandExploreArea(conf *config) error {
	if (len(conf.args) < 1) {
		return fmt.Errorf("comandExploreArea takes an area as an argument, but args length was 0")
	}
	url := fmt.Sprintf(areaURL, conf.args[0])
	var area Area
	err := makeGetRequest(url, &area)
	if err != nil {
		return fmt.Errorf("explore area request to url failed: %s\n%w",url,err)
	}
	if area.PokemonEncounters == nil || len(area.PokemonEncounters) < 1 {
		return nil
	}
	for i := 0; i< len(area.PokemonEncounters); i++ {
		fmt.Println(area.PokemonEncounters[i].Pokemon.Name)
	}

	return nil
}

func commandCatch(conf *config) error {
	if (len(conf.args) < 1) {
		return fmt.Errorf("commandCatch takes a pokemon name as an argument, but args length was 0")
	}
	url := fmt.Sprintf(pokemonURL, conf.args[0])
	var pokemon Pokemon 
	err := makeGetRequest(url, &pokemon)
	if err != nil {
		return fmt.Errorf("commandCatch Error %w", err)
	}
	// base experience values range between 35 to 635.  I'd like to not have to attempt catching a pokemon more than 7 times and no less than 1 or 2.
	// pull the pokemon base experience, divide it by 100 and round up to nearest decimal to produce the max catch value. If value is 0 or 1, increase it to 2.
	// set a random number between 1 and the max catch value
	maxDifficulty := float64(pokemon.BaseExperience) / 100.0
	//fmt.Printf("maxDifficulty = %.1f\n", maxDifficulty)
	maxDifficultyR := math.Ceil(maxDifficulty)
	//fmt.Printf("maxDifficultyR = %v\n", maxDifficultyR)

	catchChance := rand.Intn(int(maxDifficultyR + 1))
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if catchChance == 1 {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}

func commandInspect(conf *config) error {
	if (len(conf.args) < 1) {
		return fmt.Errorf("commandInspect takes a pokemon name as an argument, but args length was 0")
	}
	var pokemon Pokemon
	pokemon, ok := pokedex[conf.args[0]]
	if !ok {
		fmt.Println("Please capture the given pokemon to view its data")
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, Stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", Stat.Stat.Name, Stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, Type := range pokemon.Types {
		fmt.Printf("  - %s\n", Type.Type.Name)
	}
	return nil
}
func getLocations(url string) (Locations, error) {
	var locResults Locations
	err := makeGetRequest(url, &locResults)
	if err != nil {
		return locResults, fmt.Errorf("getLocations request to url failed: %s\n%w",url,err)
	}
	for i := 0; i< len(locResults.Results); i++ {
		fmt.Println(locResults.Results[i].Name)
	}
	return locResults, nil
}

/** Because JSON might return nullable strings in the response, 
*** assignment to a string type is taken care of gracefully.
**/
func stringReferenceAssignment(dest *string, src *string) {
	if src != nil {
		*dest = *src
	} else {
		*dest = ""
	}
}

func cleanInput(text string) []string {
	slices := strings.Fields(strings.ToLower(text))
	return slices
}
