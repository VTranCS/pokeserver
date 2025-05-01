"use server";
import { FetchedPokemonDBEntry } from "@/app/types/PokemonDBEntry";
import { db } from '@vercel/postgres';


export async function getRandomPokemon(): Promise<FetchedPokemonDBEntry> {
  const generation = process.env.POKEMON_GENERATION ? parseInt(process.env.POKEMON_GENERATION, 10) : 1;
  const client = await db.connect();
  const res = await client.sql`SELECT * FROM pokemon WHERE generation = ${generation} ORDER BY RANDOM() LIMIT 1`;
  client.release();
  return res.rows[0] as FetchedPokemonDBEntry;
}


export async function getAllPokemonDBEntry(generation: number): Promise<FetchedPokemonDBEntry[]> {
  const client = await db.connect();
  const res = await client.sql`SELECT * FROM pokemon WHERE generation = ${generation}`;

  const data = res.rows as FetchedPokemonDBEntry[];
  client.release();
  return data;
}