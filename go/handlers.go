package main

import (
	"context"
	"encoding/json"
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
	ctx := context.Background()
	
	pokemon, err := h.pokemonService.GetRandomPokemon(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get random Pokemon")
		http.Error(w, "Failed to get Pokemon", http.StatusInternalServerError)
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
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}

// HandleGetPokemon handles the /getpokemon endpoint
func (h *Handlers) HandleGetPokemon(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	
	pokemon, err := h.pokemonService.GetRandomPokemon(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get random Pokemon")
		http.Error(w, "Failed to get Pokemon", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pokemon)
}

// HandleGetAllPokemon handles the /getall endpoint
func (h *Handlers) HandleGetAllPokemon(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	
	allPokemon, err := h.pokemonService.GetAllPokemon(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get all Pokemon")
		http.Error(w, "Failed to get Pokemon list", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"pokemon": allPokemon,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleVote handles the /vote endpoint
func (h *Handlers) HandleVote(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	
	// Parse query parameters
	params := r.URL.Query()
	direction := params.Get("vote")
	pokeIDStr := params.Get("id")
	
	if direction == "" || pokeIDStr == "" {
		http.Error(w, "Missing required parameters: vote and id", http.StatusBadRequest)
		return
	}
	
	pokeID, err := strconv.Atoi(pokeIDStr)
	if err != nil {
		http.Error(w, "Invalid Pokemon ID", http.StatusBadRequest)
		return
	}

	pokemonEntry, err := h.pokemonService.VotePokemon(ctx, pokeID, direction)
	if err != nil {
		h.logger.WithError(err).Error("Failed to vote for Pokemon")
		http.Error(w, "Failed to vote for Pokemon", http.StatusInternalServerError)
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
