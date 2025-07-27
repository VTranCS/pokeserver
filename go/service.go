package main

import (
	"context"
	"fmt"
	"math/rand/v2"
)

// PokemonService handles Pokemon business logic
type PokemonService struct {
	repo   PokemonRepository
	client *PokeClient
	logger *Logger
}

// NewPokemonService creates a new Pokemon service
func NewPokemonService(repo PokemonRepository, client *PokeClient, logger *Logger) *PokemonService {
	return &PokemonService{
		repo:   repo,
		client: client,
		logger: logger,
	}
}

// GetRandomPokemon gets a random Pokemon and updates its vote
func (s *PokemonService) GetRandomPokemon(ctx context.Context) (Pokemon, error) {
	pokemon, err := s.client.GetRandomPokemon()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get random Pokemon from API")
		return Pokemon{}, fmt.Errorf("failed to get random Pokemon: %w", err)
	}

	// Get or create the Pokemon in database
	_, err = s.repo.GetPokemonDBEntry(ctx, pokemon)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get Pokemon DB entry")
		return Pokemon{}, fmt.Errorf("failed to get Pokemon DB entry: %w", err)
	}

	// Add a random vote
	voteValue := rand.IntN(20)
	_, err = s.repo.UpdatePokemonVote(ctx, pokemon.ID, voteValue)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update Pokemon vote")
		return Pokemon{}, fmt.Errorf("failed to update Pokemon vote: %w", err)
	}

	s.logger.WithFields(map[string]any{
		"pokemon_name": pokemon.Name,
		"pokemon_id":   pokemon.ID,
		"vote_added":   voteValue,
	}).Info("Successfully processed random Pokemon")

	return pokemon, nil
}

// GetAllPokemon retrieves all Pokemon from the database
func (s *PokemonService) GetAllPokemon(ctx context.Context) ([]PokeDBEntry, error) {
	pokemon, err := s.repo.GetAllPokemonDBEntry(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get all Pokemon")
		return nil, fmt.Errorf("failed to get all Pokemon: %w", err)
	}
	return pokemon, nil
}

// VotePokemon adds a vote to a Pokemon
func (s *PokemonService) VotePokemon(ctx context.Context, pokemonID int, direction string) (PokeDBEntry, error) {
	vote := 0
	switch direction {
	case "up":
		vote = 1
	case "down":
		vote = -1
	default:
		return PokeDBEntry{}, fmt.Errorf("invalid vote direction: %s", direction)
	}

	_, err := s.repo.UpdatePokemonVote(ctx, pokemonID, vote)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update Pokemon vote")
		return PokeDBEntry{}, fmt.Errorf("failed to update Pokemon vote: %w", err)
	}

	pokemonEntry, err := s.repo.GetPokemonDBEntryById(ctx, pokemonID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get updated Pokemon entry")
		return PokeDBEntry{}, fmt.Errorf("failed to get updated Pokemon entry: %w", err)
	}

	s.logger.WithFields(map[string]any{
		"pokemon_id": pokemonID,
		"vote":       vote,
		"new_total":  pokemonEntry.Vote,
	}).Info("Successfully voted for Pokemon")

	return pokemonEntry, nil
}

// InitializeDatabase creates the necessary database tables
func (s *PokemonService) InitializeDatabase(ctx context.Context) error {
	_, err := s.repo.CreatePokeVotesTable(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create Pokemon votes table")
		return fmt.Errorf("failed to create Pokemon votes table: %w", err)
	}
	return nil
}