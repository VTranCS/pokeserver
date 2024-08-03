package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

type PageData struct {
	Title string
	Name  string
	Image string
}

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
	return getPokemonByURL(pokeSum.Results[i].URL)
}

// Handler for root endpoint
func handlePokeStop(w http.ResponseWriter, r *http.Request) {
	myPokemon := getPokemon()
	pageData := PageData{
		Title: "PokeServer",
		Name:  myPokemon.Name,
		Image: myPokemon.Sprites.FrontDefault,
	}
	getPokemonVote(myPokemon.Name)
	updatePokemonVote(myPokemon.Name, rand.IntN(20))
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, pageData)
}

// get the number of votes for a pokemon
func getPokemonVote(pokename string) int {
	rows, _ := conn.Query(context.Background(), "SELECT * FROM pokevotes WHERE name = $1", pokename)
	var vote int
	var name string

	for rows.Next() {
		err := rows.Scan(&name, &vote)
		if err != nil {
			log.Print(err.Error())
		}
	}

	if rows.CommandTag().RowsAffected() < 1 {
		createPokemonVote(pokename)
	}
	return vote
}

// Create the entry in the pokevotes tables
func createPokemonVote(pokename string) bool {
	_, err := conn.Exec(context.Background(), "insert into pokevotes values($1,$2)", pokename, 0)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}

func updatePokemonVote(pokename string, vote int) bool {
	_, err := conn.Exec(context.Background(), "UPDATE pokevotes SET vote= vote + $1 WHERE name=$2", vote, pokename)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}

var conn *pgx.Conn

func main() {

	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handlePokeStop)

	fmt.Printf("Started poke app")
	httperr := http.ListenAndServe(":9091", nil)
	if errors.Is(httperr, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	}
}
