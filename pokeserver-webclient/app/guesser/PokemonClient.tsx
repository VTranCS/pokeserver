'use client';

import Image from "next/image";
import React, { useState, useTransition } from "react";
import { getRandomPokemon, type Pokemon } from "./actions/getPokemon";

export default function PokemonClient({ initialPokemon }: { initialPokemon: Pokemon }) {
  const [pokemon, setPokemon] = useState(initialPokemon);
  const [imageVisible, setImageVisible] = useState(false);
  const [disableInput, setDisableInput] = useState(false);
  const [failed, setFailed] = useState(false);
  const [failedCount, setFailedCount] = useState(0);
  const [score, setScore] = useState(0);
  const [textInput, setTextInput] = useState('');
  const [isPending, startTransition] = useTransition();

  function handleKeyDown(e: { key: string }) {
    if (e.key === 'Enter') {
      if (textInput.trim().toLowerCase() === pokemon.name.toLowerCase()) {
        setDisableInput(true);
        setImageVisible(true);
        setScore(score + 1);
      } else {
        const newFailed = failedCount + 1;
        setFailedCount(newFailed);
        if (newFailed >= 3) {
          setTextInput(pokemon.name);
          setDisableInput(true);
          setImageVisible(true);
          setFailed(false);
          setScore(score - 1);
        } else {
          setTextInput('');
          setFailed(true);
        }
      }
    } else if (failed) {
      setFailed(false);
      setTextInput('');
    }
  }

  function getAnotherPokemon() {
    startTransition(async () => {
      const newPokemon = await getRandomPokemon();
      setPokemon(newPokemon);
      setDisableInput(false);
      setImageVisible(false);
      setFailed(false);
      setFailedCount(0);
      setTextInput('');
    });
  }

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        <Image
          src={pokemon?.sprites.front_default}
          alt="Pokemon"
          width={180}
          height={180}
          className={`${imageVisible ? "" : "opacity-50 blur-md"}`}
        />
        <textarea
          disabled={disableInput}
          value={textInput}
          rows={1}
          onChange={(e) => setTextInput(e.target.value)}
          onKeyDown={handleKeyDown}
          className={`rounded-lg focus:outline-none focus:ring-0 resize-none ${failed ? "bg-red-700" : "bg-gray-700"} text-white text-center text-base`}
        />
        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <button
            disabled={isPending}
            className="rounded-full border transition-colors hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5"
            onClick={getAnotherPokemon}
          >
            {isPending ? "Loading..." : "Get Another"}
          </button>
        </div>
        <h2>Remaining Tries: {3 - failedCount}</h2>
        <h2>Score: {score}</h2>
      </main>
    </div>
  );
}
