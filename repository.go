package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	return &Repository{
		conn: conn,
	}, nil
}

// get the number of votes for a pokemon
func (r Repository) getPokemonDBEntry(pokemon Pokemon) int {
	rows, _ := r.conn.Query(context.Background(), "SELECT * FROM pokevotes WHERE name = $1", pokemon.Name)
	pokemonDBEntry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PokeDBEntry])
	if err != nil {
		log.Print(err)
	}

	if rows.CommandTag().RowsAffected() < 1 {
		r.createPokemonVote(pokemon)
	}
	return pokemonDBEntry.Vote
}

func (r Repository) getPokemonDBEntryById(id int) PokeDBEntry {
	rows, _ := r.conn.Query(context.Background(), "SELECT * FROM pokevotes WHERE id = $1", id)
	aPokeDBEntry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PokeDBEntry])
	if err != nil {
		log.Print(err)
	}

	return aPokeDBEntry
}

func (r Repository) getAllPokemonDBEntry() []PokeDBEntry {
	rows, _ := r.conn.Query(context.Background(), "SELECT * FROM pokevotes ORDER BY id ASC")
	pokemonDBEntries, err := pgx.CollectRows(rows, pgx.RowToStructByName[PokeDBEntry])
	if err != nil {
		log.Print(err)
	}

	return pokemonDBEntries
}

// Create the entry in the pokevotes tables
func (r Repository) createPokemonVote(pokemon Pokemon) bool {
	_, err := r.conn.Exec(context.Background(), "insert into pokevotes values($1,$2,$3,$4)",
		pokemon.Name, 0, pokemon.Sprites.FrontDefault, pokemon.ID)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}

func (r Repository) createPokeVotesTable() bool {
	_, err := r.conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS pokevotes ( NAME VARCHAR(100),"+
		"vote INT, Url VARCHAR(100), Id INT);")
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}

func (r Repository) updatePokemonVote(id int, vote int) bool {
	_, err := r.conn.Exec(context.Background(), "UPDATE pokevotes SET vote= vote + $1 WHERE id=$2",
		vote, id)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}
