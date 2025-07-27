package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PokemonRepository defines the interface for pokemon data operations
type PokemonRepository interface {
	GetPokemonDBEntry(ctx context.Context, pokemon Pokemon) (int, error)
	GetPokemonDBEntryById(ctx context.Context, id int) (PokeDBEntry, error)
	GetAllPokemonDBEntry(ctx context.Context) ([]PokeDBEntry, error)
	CreatePokemonVote(ctx context.Context, pokemon Pokemon) (bool, error)
	UpdatePokemonVote(ctx context.Context, id int, vote int) (bool, error)
	CreatePokeVotesTable(ctx context.Context) (bool, error)
	ResetPokeVotes(ctx context.Context) (bool, error)
}

// Repository implements PokemonRepository interface
type Repository struct {
	pool   *pgxpool.Pool
	logger *Logger
}

// NewRepository creates a new repository instance
func NewRepository(ctx context.Context, connStr string, logger *Logger) (*Repository, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		logger.WithError(err).Error("Unable to create connection pool")
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	return &Repository{
		pool:   pool,
		logger: logger,
	}, nil
}

// Close closes the database connection pool
func (r *Repository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

// GetPokemonDBEntry gets the number of votes for a pokemon
func (r *Repository) GetPokemonDBEntry(ctx context.Context, pokemon Pokemon) (int, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM pokevotes WHERE name = $1", pokemon.Name)
	if err != nil {
		r.logger.WithError(err).Error("Failed to query pokemon by name")
		return 0, fmt.Errorf("failed to query pokemon by name: %w", err)
	}
	defer rows.Close()

	pokemonDBEntry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PokeDBEntry])
	if err != nil {
		if err == pgx.ErrNoRows {
			// Pokemon doesn't exist, create it
			_, createErr := r.CreatePokemonVote(ctx, pokemon)
			if createErr != nil {
				return 0, createErr
			}
			return 0, nil
		}
		r.logger.WithError(err).Error("Failed to collect pokemon row")
		return 0, fmt.Errorf("failed to collect pokemon row: %w", err)
	}

	return pokemonDBEntry.Vote, nil
}

func (r *Repository) GetPokemonDBEntryById(ctx context.Context, id int) (PokeDBEntry, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM pokevotes WHERE id = $1", id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to query pokemon by ID")
		return PokeDBEntry{}, fmt.Errorf("failed to query pokemon by ID: %w", err)
	}
	defer rows.Close()

	aPokeDBEntry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PokeDBEntry])
	if err != nil {
		r.logger.WithError(err).Error("Failed to collect pokemon row by ID")
		return PokeDBEntry{}, fmt.Errorf("failed to collect pokemon row by ID: %w", err)
	}

	return aPokeDBEntry, nil
}

func (r *Repository) GetAllPokemonDBEntry(ctx context.Context) ([]PokeDBEntry, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM pokevotes ORDER BY id ASC")
	if err != nil {
		r.logger.WithError(err).Error("Failed to query all pokemon")
		return nil, fmt.Errorf("failed to query all pokemon: %w", err)
	}
	defer rows.Close()

	pokemonDBEntries, err := pgx.CollectRows(rows, pgx.RowToStructByName[PokeDBEntry])
	if err != nil {
		r.logger.WithError(err).Error("Failed to collect all pokemon rows")
		return nil, fmt.Errorf("failed to collect all pokemon rows: %w", err)
	}

	return pokemonDBEntries, nil
}

// CreatePokemonVote creates an entry in the pokevotes table
func (r *Repository) CreatePokemonVote(ctx context.Context, pokemon Pokemon) (bool, error) {
	_, err := r.pool.Exec(ctx, "INSERT INTO pokevotes (name, vote, url, id) VALUES ($1, $2, $3, $4)",
		pokemon.Name, 0, pokemon.Sprites.FrontDefault, pokemon.ID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create pokemon vote")
		return false, fmt.Errorf("failed to create pokemon vote: %w", err)
	}
	return true, nil
}

func (r *Repository) CreatePokeVotesTable(ctx context.Context) (bool, error) {
	_, err := r.pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS pokevotes (name VARCHAR(100), vote INT, url VARCHAR(100), id INT)")
	if err != nil {
		r.logger.WithError(err).Error("Failed to create pokevotes table")
		return false, fmt.Errorf("failed to create pokevotes table: %w", err)
	}
	return true, nil
}

func (r *Repository) UpdatePokemonVote(ctx context.Context, id int, vote int) (bool, error) {
	_, err := r.pool.Exec(ctx, "UPDATE pokevotes SET vote = vote + $1 WHERE id = $2", vote, id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to update pokemon vote")
		return false, fmt.Errorf("failed to update pokemon vote: %w", err)
	}
	return true, nil
}

func (r *Repository) ResetPokeVotes(ctx context.Context) (bool, error) {
	_, err := r.pool.Exec(ctx, "TRUNCATE pokevotes")
	if err != nil {
		r.logger.WithError(err).Error("Failed to reset pokevotes")
		return false, fmt.Errorf("failed to reset pokevotes: %w", err)
	}
	return true, nil
}

// Legacy methods for backwards compatibility - these will be removed later
func (r Repository) getPokemonDBEntry(ctx context.Context, pokemon Pokemon) (int, error) {
	return r.GetPokemonDBEntry(ctx, pokemon)
}

func (r Repository) getPokemonDBEntryById(ctx context.Context, id int) (PokeDBEntry, error) {
	return r.GetPokemonDBEntryById(ctx, id)
}

func (r Repository) getAllPokemonDBEntry(ctx context.Context) ([]PokeDBEntry, error) {
	return r.GetAllPokemonDBEntry(ctx)
}

func (r Repository) createPokemonVote(ctx context.Context, pokemon Pokemon) (bool, error) {
	return r.CreatePokemonVote(ctx, pokemon)
}

func (r Repository) createPokeVotesTable(ctx context.Context) (bool, error) {
	return r.CreatePokeVotesTable(ctx)
}

func (r Repository) updatePokemonVote(ctx context.Context, id int, vote int) (bool, error) {
	return r.UpdatePokemonVote(ctx, id, vote)
}

func (r Repository) resetPokeVotes(ctx context.Context) (bool, error) {
	return r.ResetPokeVotes(ctx)
}
