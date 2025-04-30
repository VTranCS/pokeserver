'use server';

import React from "react";
import { getRandomPokemon } from "./actions/getPokemon";
import PokemonVoteDisplay from "./voter/PokemonVoteDisplay";
import { seedPokemon } from "./actions/dbActions";
export default async function Home() {
    const generation = process.env.POKEMON_GENERATION ? parseInt(process.env.POKEMON_GENERATION, 10) : 1;
    seedPokemon(generation);
    const pokemon = await getRandomPokemon(generation);

    return (
      <div className="p-4">
        <PokemonVoteDisplay
          initialPokemon={pokemon}
          generation={generation}/>
      </div>
    );
}
