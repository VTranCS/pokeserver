package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
)

// PokeClient handles communication with the PokeAPI
type PokeClient struct {
	baseURL string
	maxID   int
	logger  *Logger
}

// NewPokeClient creates a new PokeAPI client
func NewPokeClient(baseURL string, maxID int, logger *Logger) *PokeClient {
	return &PokeClient{
		baseURL: baseURL,
		maxID:   maxID,
		logger:  logger,
	}
}

// GetPokemonByURL retrieves specific Pokemon data from PokeAPI by URL
func (c *PokeClient) GetPokemonByURL(url string) (Pokemon, error) {
	resp, err := http.Get(url)
	if err != nil {
		c.logger.WithError(err).Error("Failed to make HTTP request to PokeAPI")
		return Pokemon{}, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.WithFields(map[string]any{
			"status_code": resp.StatusCode,
			"url":         url,
		}).Error("PokeAPI returned non-200 status code")
		return Pokemon{}, fmt.Errorf("PokeAPI returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.WithError(err).Error("Failed to read response body")
		return Pokemon{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		c.logger.WithError(err).Error("Failed to unmarshal Pokemon JSON")
		return Pokemon{}, fmt.Errorf("failed to unmarshal Pokemon JSON: %w", err)
	}

	c.logger.WithFields(map[string]any{
		"pokemon_name": pokemon.Name,
		"pokemon_id":   pokemon.ID,
	}).Info("Successfully retrieved Pokemon")

	return pokemon, nil
}

// GetRandomPokemon retrieves a random Pokemon from PokeAPI
func (c *PokeClient) GetRandomPokemon() (Pokemon, error) {
	url := fmt.Sprintf("%s%d", c.baseURL, c.maxID)
	
	resp, err := http.Get(url)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get Pokemon list from PokeAPI")
		return Pokemon{}, fmt.Errorf("failed to get Pokemon list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.WithFields(map[string]any{
			"status_code": resp.StatusCode,
			"url":         url,
		}).Error("PokeAPI returned non-200 status code for Pokemon list")
		return Pokemon{}, fmt.Errorf("PokeAPI returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.WithError(err).Error("Failed to read Pokemon list response body")
		return Pokemon{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var pokeSum Pokeapi
	if err := json.Unmarshal(body, &pokeSum); err != nil {
		c.logger.WithError(err).Error("Failed to unmarshal Pokemon list JSON")
		return Pokemon{}, fmt.Errorf("failed to unmarshal Pokemon list JSON: %w", err)
	}

	if len(pokeSum.Results) == 0 {
		c.logger.Error("No Pokemon found in API response")
		return Pokemon{}, fmt.Errorf("no Pokemon found in API response")
	}

	// Select a random Pokemon
	randomIndex := rand.IntN(len(pokeSum.Results))
	return c.GetPokemonByURL(pokeSum.Results[randomIndex].URL)
}

// Legacy functions for backwards compatibility - these will be removed later
func getPokemonByURL(url string) Pokemon {
	// This should not be used in new code, but keeping for compatibility
	client := &PokeClient{logger: NewLogger()}
	pokemon, err := client.GetPokemonByURL(url)
	if err != nil {
		// Keep the old behavior for now
		panic(err)
	}
	return pokemon
}

func getPokemon(pokemonRange int) Pokemon {
	// This should not be used in new code, but keeping for compatibility
	client := &PokeClient{
		baseURL: "https://pokeapi.co/api/v2/pokemon?limit=",
		maxID:   pokemonRange,
		logger:  NewLogger(),
	}
	pokemon, err := client.GetRandomPokemon()
	if err != nil {
		// Keep the old behavior for now
		panic(err)
	}
	return pokemon
}
