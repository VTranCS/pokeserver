package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil, err
	}
	return &Repository{
		pool: pool,
	}, nil
}

// get the number of votes for a pokemon
func (r Repository) getPokemonDBEntry(ctx context.Context, pokemon Pokemon) (int, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM pokevotes WHERE name = $1", pokemon.Name)

	if err != nil {
		log.Print(err)
	}

	defer rows.Close()

	pokemonDBEntry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PokeDBEntry])

	if err != nil {
		log.Print(err)
	}

	if rows.CommandTag().RowsAffected() < 1 {
		r.createPokemonVote(context.Background(), pokemon)
	}
	return pokemonDBEntry.Vote, err
}

func (r Repository) getPokemonDBEntryById(ctx context.Context, id int) (PokeDBEntry, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM pokevotes WHERE id = $1", id)

	if err != nil {
		log.Print(err)
	}
	defer rows.Close()

	aPokeDBEntry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PokeDBEntry])

	if err != nil {
		log.Print(err)
	}

	return aPokeDBEntry, err
}

func (r Repository) getAllPokemonDBEntry(ctx context.Context) ([]PokeDBEntry, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM pokevotes ORDER BY id ASC")

	if err != nil {
		log.Print(err)
	}

	defer rows.Close()

	pokemonDBEntries, err := pgx.CollectRows(rows, pgx.RowToStructByName[PokeDBEntry])

	if err != nil {
		log.Print(err)
	}

	return pokemonDBEntries, nil
}

// Create the entry in the pokevotes tables
func (r Repository) createPokemonVote(ctx context.Context, pokemon Pokemon) (bool, error) {
	_, err := r.pool.Exec(context.Background(), "insert into pokevotes values($1,$2,$3,$4)",
		pokemon.Name, 0, pokemon.Sprites.FrontDefault, pokemon.ID)
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	return true, nil
}

func (r Repository) createPokeVotesTable(ctx context.Context) (bool, error) {
	_, err := r.pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS pokevotes ( NAME VARCHAR(100),"+
		"vote INT, Url VARCHAR(100), Id INT);")
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	return true, nil
}

func (r Repository) updatePokemonVote(ctx context.Context, id int, vote int) (bool, error) {
	_, err := r.pool.Exec(context.Background(), "UPDATE pokevotes SET vote= vote + $1 WHERE id=$2",
		vote, id)
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	return true, nil
}

func (r Repository) resetPokeVotes(ctx context.Context) (bool, error) {
	_, err := r.pool.Exec(context.Background(), "TRUNCATE pokevotes")
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	return true, nil
}
