"use client";

import Image from "next/image";
import React, { startTransition, useState } from "react";
import { getRandomPokemon } from "../actions/getPokemon";
import { votePokemon } from "../actions/votePokemon";
import { FetchedPokemonDBEntry } from "../types/PokemonDBEntry";
export default function PokemonVoteDisplay({}) {
  const defaultVoteText = "???";
  const [pokemon, setPokemon] = useState<FetchedPokemonDBEntry | null>(null);
  const [score, setScore] = useState(defaultVoteText);
  const [disableInput, setDisableInput] = useState(false);

  async function sendVote(
    pokemonId: number,
    action: "up" | "down"
  ): Promise<number> {
    const results = await votePokemon(pokemonId, action);
    setScore(String(results));
    setDisableInput(true);
    return 0;
  }

  const getAnotherPokemon = React.useCallback(() => {
    startTransition(async () => {
      const newPokemon = await getRandomPokemon();
      setPokemon(newPokemon);
      setDisableInput(false);
    });
  }, []);

  React.useEffect(() => {
    getAnotherPokemon();
  }, [getAnotherPokemon]);

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        <div className="flex justify-center items-center h-[360px] ">
          {pokemon === null ? (
            <div className="w-[370px] h-[370px] border-4 border-purple-600 border-t-transparent rounded-full animate-spin" />
          ) : (
            <Image
              src={pokemon ? pokemon?.image : "/file.svg"}
              alt="Next.js logo"
              width={360}
              height={370}
            />
          )}
        </div>
        <div className="flex justify-center items-center w-full">
          <span className="text-2xl">
            {pokemon === null ? "Loading" : pokemon.name.toUpperCase()}
          </span>
        </div>
        <div className="flex justify-center items-center w-full ">
          <span className="text-2xl">Score: {score} </span>
        </div>
        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <div className="flex gap-4 items-center flex-col sm:flex-row">
            <button
              disabled={disableInput}
              className={`${
                disableInput ? "bg-gray-500" : "bg-green-500"
              } rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44`}
              onClick={() => {
                if (pokemon?.id !== undefined) sendVote(pokemon.id, "up");
              }}
            >
              Upvote
            </button>
            <button
              disabled={disableInput}
              className={`${
                disableInput ? "bg-gray-500" : "bg-red-500"
              } rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44`}
              onClick={() => {
                if (pokemon?.id !== undefined) sendVote(pokemon.id, "down");
              }}
            >
              Downvote
            </button>
          </div>
        </div>
        <div className="flex justify-center items-center w-full">
          <button
            className="bg-purple-600 rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44"
            onClick={() => {
              setScore(defaultVoteText);
              getAnotherPokemon();
            }}
          >
            Next
          </button>
        </div>
      </main>
    </div>
  );
}
