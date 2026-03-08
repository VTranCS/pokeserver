package main

// Pokemon represents the structure of Pokemon data from PokeAPI.
type Pokemon struct {
	ID      int    `json:"id"`   // Unique Pokemon identifier
	Name    string `json:"name"` // Pokemon name
	Sprites struct {
		FrontDefault string `json:"front_default"` // URL to default front sprite
		Other        struct {
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"` // URL to official artwork
				FrontShiny   string `json:"front_shiny"`   // URL to shiny artwork
			} `json:"official-artwork"`
		} `json:"other"`
	} `json:"sprites"`
}

// Pokeapi represents the response structure from PokeAPI list endpoint.
type Pokeapi struct {
	Count    int `json:"count"`    // Total number of Pokemon
	Next     any `json:"next"`     // URL to next page (if any)
	Previous any `json:"previous"` // URL to previous page (if any)
	Results  []struct {
		Name string `json:"name"` // Pokemon name
		URL  string `json:"url"`  // URL to detailed Pokemon data
	} `json:"results"`
}
