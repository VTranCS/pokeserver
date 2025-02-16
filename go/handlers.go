package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"text/template"

	"github.com/spf13/viper"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// Handler for root endpoint
func handlePokeStop(w http.ResponseWriter, r *http.Request) {
	myPokemon := getPokemon(viper.GetInt("pokeapi.max"))
	pageData := IndexPageData{
		Title: "PokeServer",
		Name:  myPokemon.Name,
		Image: myPokemon.Sprites.FrontDefault,
		Id:    strconv.Itoa(myPokemon.ID),
	}

	repo.getPokemonDBEntry(myPokemon)
	repo.updatePokemonVote(myPokemon.ID, rand.IntN(20))
	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	tmpl.Execute(w, pageData)
}

// Handler for root endpoint
func handleGetPokemon(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	myPokemon := getPokemon(viper.GetInt("pokeapi.max"))
	json.NewEncoder(w).Encode(myPokemon)
}

// Handler for root endpoint
func handleShowAllPokemon(w http.ResponseWriter, r *http.Request) {
	allPokemon, err := repo.getAllPokemonDBEntry()
	if err != nil {
		log.Print(err.Error())

	}
	pageData := ShowAllPageData{
		Title:   "All Pokemon",
		Pokemon: allPokemon,
	}
	tmpl := template.Must(template.ParseFiles("static/templates/getallpokemon.html"))
	tmpl.Execute(w, pageData)
}

type PokeDBList struct {
	Data []PokeDBEntry
}

func handleGetAllPokemon(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	allPokemon, err := repo.getAllPokemonDBEntry()
	if err != nil {
		log.Print(err.Error())

	}

	response := map[string]interface{}{
		"pokemon": allPokemon,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

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
	repo.updatePokemonVote(pokeId, 1*vote)
	aPokeDBEntry, err := repo.getPokemonDBEntryById(pokeId)
	if err != nil {
		log.Print(err.Error())
	}
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
