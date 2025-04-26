"use server";
import { Pokemon } from "@/app/types/Pokemon";
import { PokemonDBEntry } from "@/app/types/PokemonDBEntry";
export async function getRandomPokemon(): Promise<Pokemon> {
  const res = await fetch(`${process.env.POKEMON_BASE_URL}/${process.env.SUFFIX_GET_POKEMON}`, {
    cache: 'no-store',
  });

  if (!res.ok) throw new Error('Failed to fetch Pokémon');

  return res.json();
}


export async function getAllPokemonDBEntry(): Promise<PokemonDBEntry[]> {
  const res = await fetch(`${process.env.POKEMON_BASE_URL}/${process.env.SUFFIX_GET_ALL_POKEMON}`, {
    cache: 'no-store',
  });

  if (!res.ok) throw new Error('Failed to fetch Pokémon votes');


  const data = await res.json();
  return data.pokemon;
}