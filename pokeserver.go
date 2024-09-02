package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

var repo *Repository

func main() {
	viper.SetConfigFile("./.env")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	repo, err = NewRepository(context.Background(), viper.GetString("database.url"))
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	repo.createPokeVotesTable()
	http.HandleFunc("/getall", handleShowAllPokemon)
	http.HandleFunc("/vote", handleVote)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handlePokeStop)

	port := fmt.Sprintf(":%s", viper.GetString("server.port"))
	fmt.Printf("Started poke app on http://localhost%s", port)
	httperr := http.ListenAndServe(port, nil)
	if errors.Is(httperr, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	}
}

type PokeDBEntry struct {
	Id   int
	Name string
	Vote int
	Url  string
}
