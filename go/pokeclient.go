package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/spf13/viper"
)

// Retrieve specific Pokemon Data for PokeApi
func getPokemonByURL(url string) Pokemon {
	resp, getErr := http.Get(url)
	if getErr != nil || resp.StatusCode != 200 {
		log.Fatal(getErr)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	pokemon := Pokemon{}
	jsonErr := json.Unmarshal(body, &pokemon)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	fmt.Printf("The pokemon is %v!\n", pokemon.Name)
	return pokemon
}

func getPokemon(pokemonRange int) Pokemon {
	url := fmt.Sprintf("%s%d", viper.GetString("pokeapi.url"), pokemonRange)
	resp, getErr := http.Get(url)
	if getErr != nil || resp.StatusCode != 200 {
		log.Fatal(resp.StatusCode)
		log.Fatal(getErr)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	pokeSum := Pokeapi{}
	jsonErr := json.Unmarshal(body, &pokeSum)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	var i = rand.IntN(len(pokeSum.Results))
	return getPokemonByURL(pokeSum.Results[i].URL)
}
