"use server";
import { FetchedPokemonDBEntry, PokemonDBEntry } from "@/app/types/PokemonDBEntry";
import { db } from '@vercel/postgres';
const client = await db.connect();

export async function getRandomPokemon(generation: number ): Promise<FetchedPokemonDBEntry> {
  const res = await client.sql`SELECT * FROM pokemon WHERE generation = ${generation} ORDER BY RANDOM() LIMIT 1`;

  return res.rows[0] as FetchedPokemonDBEntry;
}


export async function getAllPokemonDBEntry(): Promise<PokemonDBEntry[]> {
  const res = await fetch(`${process.env.POKEMON_BASE_URL}/${process.env.SUFFIX_GET_ALL_POKEMON}`, {
    cache: 'no-store',
  });

  if (!res.ok) throw new Error('Failed to fetch Pokémon votes');


  const data = await res.json();
  return data.pokemon;
}