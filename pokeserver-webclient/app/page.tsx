'use client';

import Image from "next/image";
import React, { useState, useEffect } from "react";
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
  const [nameVisible, setNameVisible] = useState(false);

  useEffect(() => {

    const fetchRandomPokemon = async () => {
      console.log('url', process.env);
      fetch((process.env.NEXT_PUBLIC_POKEMON_BASE_URL || 'localhost:3000')
        + (process.env.NEXT_PUBLIC_SUFFIX_GET_POKEMON || '/pokemon'))
        .then(response => response.json())
        // 4. Setting *dogImage* to the image url that we received from the response above
        .then(data => {
          console.log(data);
          setPokemon(data)
        })
      setFetchTrigger(false); // Reset trigger
    };

    fetchRandomPokemon();
  }, [fetchTrigger]); // Only runs when fetchTrigger changes
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        <Image
          src={pokemon ? pokemon.sprites.front_default : '/nextjs.svg'}
          alt="Next.js logo"
          width={180}
          height={38}
          priority
          className="opacity-50 blur-md hover:opacity-100 hover:blur-none"
          onMouseEnter={() => setNameVisible(true)}
        />
        <div className="mt-2 h-6">
          {nameVisible &&         
          <h2 className="text-center font-[family-name:var(--font-geist-mono)]">
            {pokemon ? pokemon.name : 'Loading...'}
          </h2>}

        </div>

  

        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <button
            className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44"
            onClick={() => {setFetchTrigger(true); setNameVisible(false);}}
          >
            Get Another
          </button>
          <Link
            className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44"
            onClick={() => {setFetchTrigger(true); setNameVisible(false);}}
            href="/collection"
          >
            Collection
          </Link>
        </div>
      </main>
      <footer className="row-start-3 flex gap-6 flex-wrap items-center justify-center">
      </footer>
    </div>
  );
}
