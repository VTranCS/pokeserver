package main

type Pokemon struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
		Other        struct {
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
		} `json:"other"`
	} `json:"sprites"`
}

type Pokeapi struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
