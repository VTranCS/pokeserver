package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

var conn *pgx.Conn

func main() {
	viper.SetConfigFile("./.env")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	conn, err = pgx.Connect(context.Background(), viper.GetString("database.url"))
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

type PokeDBEntry struct {
	Id   int
	Name string
	Vote int
	Url  string
}
