'use client';

import Image from "next/image";
import React, { useState, useEffect, useRef } from "react";
import { revalidatePath } from 'next/cache';
import Link from "next/link";
export default function Home() {

  interface Pokemon {
    name: string;

    sprites: {
      front_default: string;
    };
  }

  let [pokemon, setPokemon] = useState<Pokemon | null>(null)
  const [fetchTrigger, setFetchTrigger] = useState(true); // State to trigger useEffect
  const [imageVisible, setImageVisible] = useState(false);
  const [disableInput, setDisableInput] = useState(false);
  const [failed, setFailed] = useState(false);
  const [failedCount, setFailedCount] = useState(0);
  const buttonRef = useRef<HTMLButtonElement | null>(null);
  const [score, setScore] = useState(0);

  useEffect(() => {

    const fetchRandomPokemon = async () => {
      setDisableInput(false);
      setImageVisible(false);
      setFailed(false);
      setFailedCount(0);
      setTextInput('');
      console.log('url', process.env);
      fetch((process.env.NEXT_PUBLIC_POKEMON_BASE_URL || 'localhost:3000')
        + (process.env.NEXT_PUBLIC_SUFFIX_GET_POKEMON || '/pokemon'))
        .then(response => response.json())
        .then(data => {
          console.log(data);
          setPokemon(data)
        })
      setFetchTrigger(false); // Reset trigger
    };

    fetchRandomPokemon();
  }, [fetchTrigger]); // Only runs when fetchTrigger changes

  const [textInput, setTextInput] = useState('');

  function handleKeyDown(e: { key: string; }) {
    if (e.key === 'Enter') {
      console.log('Enter key pressed');
      if (textInput === pokemon?.name) {
        setDisableInput(true);
        setImageVisible(true);
        setScore(score + 1);
      }
      else {
        setTextInput('');
        setFailedCount(failedCount + 1);
        if (failedCount >= 2) {
          setTextInput(pokemon?.name || '');
          setDisableInput(true);
          setImageVisible(true);
          setFailed(false);
          setScore(score - 1);
        }
        else {
          setFailed(true);
          setFailedCount(failedCount + 1);
        }
      }

    }

    else {
      if (failed) {
        setFailed(false);
        setTextInput('');
      }

    }
  }

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        <Image
          src={pokemon ? pokemon.sprites.front_default : '/nextjs.svg'}
          alt="Next.js logo"
          width={180}
          height={38}
          priority
          className={`${imageVisible ? "" : "opacity-50 blur-md"}`}
        />
        <div className="mt-2 h-6">
          <textarea
            disabled={disableInput}
            value={textInput}
            rows={1}
            onChange={(e) => setTextInput(e.target.value)}
            onKeyDown={(e) => { handleKeyDown(e) }}
            className={`rounded-lg focus:outline-none focus:ring-0 resize-none ${failed ? "bg-red-700" : "bg-gray-700"}  text-white text-center text-sm sm:text-base font-[family-name:var(--font-geist-mono)] flex items-center justify-center`} />
        </div>
        <div
          className="flex gap-4 items-center flex-col sm:flex-row">
        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <button
            className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44"
            ref={buttonRef}
            onMouseDown={(e) => e.preventDefault()}
            onClick={() => { setFetchTrigger(true); setImageVisible(false); }}
          >
            Get Another
          </button>
          <Link
            className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44"
            href="/collection"
          >
            Collection
          </Link>
        </div>
        </div>

        <h2>Remaining Tries: {3-failedCount}</h2>
        <h2>Score: {score}</h2>
      </main>
      <footer className="row-start-3 flex gap-6 flex-wrap items-center justify-center">
      </footer>
    </div>
  );
}
