package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

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
