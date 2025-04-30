export interface PokemonDBEntry {
  Name: string;
  Id: number;
  Vote: number;
  Url: string;
}

export interface FetchedPokemon {
  name: string;
  url: string;
}

export interface FetchedPokemonDBEntry {
  id: number;
  generation: number;
  name: string;
  image: string;
  vote: number;
}
