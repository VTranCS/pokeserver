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
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

var conn *pgx.Conn

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handlePokeStop)
	http.HandleFunc("/getall", handleShowAllPokemon)
	http.HandleFunc("/vote", handleVote)

	port := fmt.Sprintf(":%s", viper.GetString("server.port"))
	fmt.Printf("Started poke app on http://localhost%s", port)
	httperr := http.ListenAndServe(port, nil)
	if errors.Is(httperr, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	}
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
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon?limit=%d", pokemonRange)
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

// get the number of votes for a pokemon
func getPokemonDBEntry(pokemon Pokemon) int {
	rows, _ := conn.Query(context.Background(), "SELECT * FROM pokevotes WHERE name = $1", pokemon.Name)
	pokemonDBEntry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PokeDBEntry])
	if err != nil {
		log.Print(err)
	}

	if rows.CommandTag().RowsAffected() < 1 {
		createPokemonVote(pokemon)
	}
	return pokemonDBEntry.Vote
}

func getPokemonDBEntryById(id int) PokeDBEntry {
	rows, _ := conn.Query(context.Background(), "SELECT * FROM pokevotes WHERE id = $1", id)
	aPokeDBEntry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PokeDBEntry])
	if err != nil {
		log.Print(err)
	}

	return aPokeDBEntry
}

func getAllPokemonDBEntry() []PokeDBEntry {
	rows, _ := conn.Query(context.Background(), "SELECT * FROM pokevotes ORDER BY id ASC")
	pokemonDBEntries, err := pgx.CollectRows(rows, pgx.RowToStructByName[PokeDBEntry])
	if err != nil {
		log.Print(err)
	}

	return pokemonDBEntries
}

// Create the entry in the pokevotes tables
func createPokemonVote(pokemon Pokemon) bool {
	_, err := conn.Exec(context.Background(), "insert into pokevotes values($1,$2,$3,$4)",
		pokemon.Name, 0, pokemon.Sprites.FrontDefault, pokemon.ID)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}

func updatePokemonVote(id int, vote int) bool {
	_, err := conn.Exec(context.Background(), "UPDATE pokevotes SET vote= vote + $1 WHERE id=$2", vote, id)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
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

type PokeDBEntry struct {
	Id   int
	Name string
	Vote int
	Url  string
}
