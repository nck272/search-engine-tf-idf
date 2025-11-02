package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

const PATH_POKEDEX string = "../data/pokemon.json"

type pokemon struct {
	Id          int               `json:"id"`
	Name        map[string]string `json:"name"`
	Types       []string          `json:"type"`
	Description string            `json:"description"`
}

func GetPokedex() []pokemon {
	f, err := os.ReadFile(PATH_POKEDEX)
	if err != nil {
		log.Fatalf("ERROR: Could not read data from %v: %v", PATH_POKEDEX, err)
	}

	var pokemons []pokemon
	json.Unmarshal(f, &pokemons)

	return pokemons
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("ERROR: missing arguements, please input the text you want to search!")
	}
	input_str := strings.Join(os.Args[1:], " ")

	// Get pokemon descriptios
	pokemons := GetPokedex()
	search_results := []string{}
	input_tokens := Tokenize(input_str)
	if len(input_tokens) == 1 {
		for _, pokemon := range pokemons {
			if strings.Contains(Standardlize(pokemon.Description), Standardlize(input_str)) {
				search_results = append(search_results, pokemon.Name["english"])
			}
		}
	} else {
		// Search and Sort the result descending by cosine_similarity
		doc_tokens := GetDocTokens(pokemons)
		search_service := GetSearchService(doc_tokens)
		search_results = search_service.Search(input_tokens, true)
	}

	// Return result which is a list of description that match with input text
	if len(search_results) > 0 {
		log.Println("--- BEST MATCHES: ")
		for i, result := range search_results {
			log.Printf("[%3d] %v", i, result)
		}
	} else {
		log.Print("FOUND NOTHING!")
	}
}
