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

func GetDocTokens(pokemons []pokemon) *doc_tokens {
	doc_tokens := doc_tokens{}
	for _, pokemon := range pokemons {
		name := pokemon.Name["english"]
		tokens := Tokenize(pokemon.Description)
		for _, t := range pokemon.Types {
			tokens = append(tokens, t)
		}
		doc_tokens[name] = tokens
	}
	return &doc_tokens
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("ERROR: missing arguements, please input the text you want to search!")
	}
	input_str := strings.Join(os.Args[1:], " ")

	// Get pokemon descriptios
	pokemons := GetPokedex()
	doc_tokens := GetDocTokens(pokemons)

	// Search and Sort the result descending by cosine_similarity
	search_service := GetSearchService(*doc_tokens)
	search_results := search_service.Search(input_str, true)

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
