package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"text/template"

	"github.com/spf13/viper"
)

// Handler for root endpoint
func handlePokeStop(w http.ResponseWriter, r *http.Request) {
	myPokemon := getPokemon(viper.GetInt("pokeapi.max"))
	pageData := IndexPageData{
		Title: "PokeServer",
		Name:  myPokemon.Name,
		Image: myPokemon.Sprites.FrontDefault,
		Id:    strconv.Itoa(myPokemon.ID),
	}
	getPokemonDBEntry(myPokemon)
	updatePokemonVote(myPokemon.ID, rand.IntN(20))
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, pageData)
}

// Handler for root endpoint
func handleShowAllPokemon(w http.ResponseWriter, r *http.Request) {
	allPokemon := getAllPokemonDBEntry()
	pageData := ShowAllPageData{
		Title:   "All Pokemon",
		Pokemon: allPokemon,
	}
	tmpl := template.Must(template.ParseFiles("templates/getallpokemon.html"))
	tmpl.Execute(w, pageData)
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	paramters := r.URL.Query()
	direction := paramters.Get("vote")
	pokeId, _ := strconv.Atoi(paramters.Get("id"))
	vote := 0
	if direction == "down" {
		vote = -1
	} else if direction == "up" {
		vote = 1
	}
	updatePokemonVote(pokeId, 1*vote)
	aPokeDBEntry := getPokemonDBEntryById(pokeId)
	fmt.Fprint(w, aPokeDBEntry.Vote)
}

type IndexPageData struct {
	Title string
	Name  string
	Image string
	Id    string
}

type ShowAllPageData struct {
	Title   string
	Pokemon []PokeDBEntry
}
