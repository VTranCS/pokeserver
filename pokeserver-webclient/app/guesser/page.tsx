"use client";

import Image from "next/image";
import React, { useEffect, useRef, useState, useTransition } from "react";
import { FetchedPokemonDBEntry } from "../types/PokemonDBEntry";
import { getRandomPokemon } from "../actions/getPokemon";

export default function PokemonClient({}) {
  const [pokemon, setPokemon] = useState<FetchedPokemonDBEntry | null>(null);
  const [imageVisible, setImageVisible] = useState(false);
  const [disableInput, setDisableInput] = useState(false);
  const [failed, setFailed] = useState(false);
  const [failedCount, setFailedCount] = useState(0);
  const [score, setScore] = useState(0);
  const [textInput, setTextInput] = useState("");
  const [isPending, startTransition] = useTransition();
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const getAnotherPokemon = React.useCallback(() => {
    startTransition(async () => {
      const newPokemon = await getRandomPokemon();
      setPokemon(newPokemon);
      setDisableInput(false);
      setImageVisible(false);
      setFailed(false);
      setFailedCount(0);
      setTextInput("");
    });
  }, []);

  React.useEffect(() => {
    getAnotherPokemon();
  }, [getAnotherPokemon]);

  useEffect(() => {
    const activeTag = document.activeElement?.tagName;
    const handleKeyDown = (e: KeyboardEvent) => {
      const isLetter = /^[a-zA-Z]$/.test(e.key);
      if (isLetter && activeTag !== "INPUT") {
        inputRef.current?.focus();
      }

      if (e.key === "Enter" && document.activeElement?.tagName !== "INPUT") {
        if (!imageVisible) {
          setScore(score - 1);
        }

        getAnotherPokemon();
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  });

  function handleKeyDown(e: { key: string }) {
    if (e.key === "Enter") {
      // Correct answer
      if (textInput.trim().toLowerCase() === pokemon?.name.toLowerCase()) {
        setDisableInput(true);
        setImageVisible(true);
        setScore(score + 1);
      } else {
        // Incorrect answer
        const newFailed = failedCount + 1;

        setFailedCount(newFailed);
        // Failed too many times
        if (newFailed >= 3) {
          setTextInput(pokemon?.name || "");
          setDisableInput(true);
          setImageVisible(true);
          setFailed(false);
          setScore(score - 1);
          // Reset for next guess
        } else {
          setTextInput("");
          setFailed(true);
        }
      }
    } else if (failed) {
      setFailed(false);
      setTextInput("");
    }
  }

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
      {pokemon === null || isPending ? (
        <div className="flex justify-center items-center h-40">
          <div className="w-12 h-12 border-4 border-purple-600 border-t-transparent rounded-full animate-spin" />
        </div>
      ) : (
        <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
          <Image
            src={pokemon?.image || "/file.svg"}
            alt="Pokemon"
            width={360}
            height={360}
            className={`${imageVisible ? "" : "filter brightness-0 invert"}`}
          />
          <textarea
            ref={inputRef}
            disabled={disableInput}
            value={textInput}
            rows={1}
            onChange={(e) => setTextInput(e.target.value)}
            onKeyDown={handleKeyDown}
            className={`overflow-hidden rounded-lg focus:outline-none focus:ring-0 resize-none ${
              failed ? "bg-red-700" : "bg-gray-700"
            } text-white text-center text-base`}
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
      )}
    </div>
  );
}
