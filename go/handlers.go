package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"text/template"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	pokemonService *PokemonService
	logger         *Logger
}

// NewHandlers creates a new handlers instance
func NewHandlers(pokemonService *PokemonService, logger *Logger) *Handlers {
	return &Handlers{
		pokemonService: pokemonService,
		logger:         logger,
	}
}

// HandleHealth provides a health check endpoint
func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// HandlePokeStop handles the root endpoint
func (h *Handlers) HandlePokeStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET method is allowed", nil)
		return
	}

	ctx := context.Background()
	
	pokemon, err := h.pokemonService.GetRandomPokemon(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get random Pokemon")
		WriteErrorResponse(w, http.StatusInternalServerError, "pokemon_fetch_failed", "Failed to get Pokemon", nil)
		return
	}

	pageData := IndexPageData{
		Title: "PokeServer",
		Name:  pokemon.Name,
		Image: pokemon.Sprites.FrontDefault,
		Id:    strconv.Itoa(pokemon.ID),
	}

	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	if err := tmpl.Execute(w, pageData); err != nil {
		h.logger.WithError(err).Error("Failed to execute template")
		WriteErrorResponse(w, http.StatusInternalServerError, "template_error", "Failed to render page", nil)
		return
	}
}

// HandleGetPokemon handles the /getpokemon endpoint
func (h *Handlers) HandleGetPokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET method is allowed", nil)
		return
	}

	ctx := context.Background()
	
	pokemon, err := h.pokemonService.GetRandomPokemon(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get random Pokemon")
		WriteErrorResponse(w, http.StatusInternalServerError, "pokemon_fetch_failed", "Failed to get Pokemon", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pokemon)
}

// HandleGetAllPokemon handles the /getall endpoint
func (h *Handlers) HandleGetAllPokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET method is allowed", nil)
		return
	}

	ctx := context.Background()
	
	allPokemon, err := h.pokemonService.GetAllPokemon(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get all Pokemon")
		WriteErrorResponse(w, http.StatusInternalServerError, "pokemon_list_failed", "Failed to get Pokemon list", nil)
		return
	}

	response := map[string]interface{}{
		"pokemon": allPokemon,
		"count":   len(allPokemon),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleVote handles the /vote endpoint
func (h *Handlers) HandleVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET method is allowed", nil)
		return
	}

	ctx := context.Background()
	
	// Parse and validate query parameters
	params := r.URL.Query()
	direction := params.Get("vote")
	pokeIDStr := params.Get("id")
	
	// Validate required parameters
	validationErrors := ValidateQueryParams(r, []string{"vote", "id"})
	if len(validationErrors) > 0 {
		WriteErrorResponse(w, http.StatusBadRequest, "validation_error", "Invalid request parameters", validationErrors)
		return
	}
	
	// Validate vote direction
	if err := ValidateVoteDirection(direction); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "invalid_vote_direction", err.Error(), nil)
		return
	}
	
	// Validate Pokemon ID
	pokeID, err := ValidatePokemonID(pokeIDStr)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "invalid_pokemon_id", err.Error(), nil)
		return
	}

	// Vote for Pokemon
	pokemonEntry, err := h.pokemonService.VotePokemon(ctx, pokeID, direction)
	if err != nil {
		h.logger.WithError(err).Error("Failed to vote for Pokemon")
		// Check if it's a "not found" error
		if fmt.Sprintf("%v", err) == fmt.Sprintf("Pokemon with ID %d not found", pokeID) {
			WriteErrorResponse(w, http.StatusNotFound, "pokemon_not_found", err.Error(), nil)
			return
		}
		WriteErrorResponse(w, http.StatusInternalServerError, "vote_failed", "Failed to vote for Pokemon", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pokemonEntry)
}

// Legacy handlers for backwards compatibility

func handlePokeStop(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	myPokemon := getPokemon(1025) // Default max value
	pageData := IndexPageData{
		Title: "PokeServer",
		Name:  myPokemon.Name,
		Image: myPokemon.Sprites.FrontDefault,
		Id:    strconv.Itoa(myPokemon.ID),
	}

	repo.getPokemonDBEntry(ctx, myPokemon)
	repo.updatePokemonVote(ctx, myPokemon.ID, rand.IntN(20))
	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	tmpl.Execute(w, pageData)
}

func handleGetPokemon(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	myPokemon := getPokemon(1025) // Default max value
	repo.getPokemonDBEntry(ctx, myPokemon)
	repo.updatePokemonVote(ctx, myPokemon.ID, rand.IntN(20))
	json.NewEncoder(w).Encode(myPokemon)
}

func handleGetAllPokemon(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	allPokemon, err := repo.getAllPokemonDBEntry(ctx)
	if err != nil {
		h := &Handlers{logger: NewLogger()}
		h.logger.WithError(err).Error("Failed to get all Pokemon")
	}

	response := map[string]interface{}{
		"pokemon": allPokemon,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := r.URL.Query()
	direction := params.Get("vote")
	pokeId, _ := strconv.Atoi(params.Get("id"))
	vote := 0
	if direction == "down" {
		vote = -1
	} else if direction == "up" {
		vote = 1
	}
	repo.updatePokemonVote(ctx, pokeId, 1*vote)
	aPokeDBEntry, err := repo.getPokemonDBEntryById(ctx, pokeId)
	if err != nil {
		h := &Handlers{logger: NewLogger()}
		h.logger.WithError(err).Error("Failed to get Pokemon by ID")
	}
	json.NewEncoder(w).Encode(aPokeDBEntry)
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
