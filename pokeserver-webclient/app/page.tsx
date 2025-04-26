'use server';

import React from "react";
import { getRandomPokemon } from "./actions/getPokemon";
import PokemonVoteDisplay from "./voter/PokemonVoteDisplay";
export default async function Home() {
    const pokemon = await getRandomPokemon();

    return (
      <div className="p-4">
        <PokemonVoteDisplay initialPokemon={pokemon} />
      </div>
    );
}
