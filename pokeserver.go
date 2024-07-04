package main

import (
	// "context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	// "github.com/jackc/pgx/v4"
)

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

func getPokemon() Pokemon {
	url := "https://pokeapi.co/api/v2/pokemon?limit=151"
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
	fmt.Printf("Pokemon id: %d\n", i)
	return getPokemonByURL(pokeSum.Results[i].URL)
}

func handlePokeStop(w http.ResponseWriter, r *http.Request) {
	var myPokemon = getPokemon()
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>%v</h1> "+
		"<img src=\"%v\"></body></html>", myPokemon.Name, myPokemon.Sprites.FrontDefault)
}

func main() {
	http.HandleFunc("/", handlePokeStop)

	fmt.Printf("Started poke app")
	httperr := http.ListenAndServe(":9091", nil)
	if errors.Is(httperr, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	}
}
