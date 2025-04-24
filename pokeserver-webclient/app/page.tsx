'use client';

import Image from "next/image";
import React, { useState, useEffect} from "react";
export default function Home() {

    interface Pokemon {
        name: string;
        id: number;

        sprites: {
            front_default: string;
        };
    }
    const defaultVoteText = "???";
    const [pokemon, setPokemon] = useState<Pokemon | null>(null)
    const [score, setScore] = useState(defaultVoteText);
    const [disableInput, setDisableInput] = useState(false);
    const [fetchTrigger, setFetchTrigger] = useState(true); // State to trigger useEffect
 
    async function votePokemon(pokemonId: number, action: "up" | "down"): Promise<number> {
        const url = `http://localhost:9091/vote?id=${pokemonId}&vote=${action}`;
        const response = await fetch(url);
        console.log("RESPONSE", response);
        const data = await response.json();
        setScore(data.Vote);
        setDisableInput(true);
        return 0;
    }
    useEffect(() => {
        const fetchRandomPokemon = async () => {
            console.log('url', process.env);
            fetch((process.env.NEXT_PUBLIC_POKEMON_BASE_URL || 'localhost:3000')
                + (process.env.NEXT_PUBLIC_SUFFIX_GET_POKEMON || '/pokemon'))
                .then(response => response.json())
                .then(data => {
                    console.log(data);
                    setPokemon(data)
                })
            setFetchTrigger(false); // Reset trigger
            setDisableInput(false);
            setScore(defaultVoteText);
        };

        fetchRandomPokemon();
        setFetchTrigger(false);
    }, [fetchTrigger]); // Only runs when fetchTrigger changes

    return (
        <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
            <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
                <div className="flex justify-center items-center w-full">
                    <Image
                        src={pokemon ? pokemon.sprites.front_default : '/nextjs.svg'}
                        alt="Next.js logo"
                        width={360}
                        height={100}
                        priority
                    />
                </div>
                <div className="flex justify-center items-center w-full">
                    <span className="text-2xl">{pokemon?.name.toUpperCase()}</span>
                </div>
                <div className="flex justify-center items-center w-full ">
                    <span className="text-2xl">Score: {score} </span>
                </div>
                <div
                    className="flex gap-4 items-center flex-col sm:flex-row">
                    <div className="flex gap-4 items-center flex-col sm:flex-row">
                        <button
                            disabled={disableInput}
                            className={`${disableInput ? 'bg-gray-500' : 'bg-green-500'} rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44`}
                            onClick={() => { if (pokemon?.id !== undefined) votePokemon(pokemon.id, "up"); }}
                        >
                            Upvote
                        </button>
                        <button
                            disabled={disableInput}
                            className={`${disableInput ? 'bg-gray-500' : 'bg-red-500'} rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44`}
                            onClick={() => { if (pokemon?.id !== undefined) votePokemon(pokemon.id, "down"); }}
                        >
                            Downvote
                        </button>
                    </div>

                </div>
                <div className="flex justify-center items-center w-full">
                    <button
                        className="bg-purple-600 rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44"
                        onClick={() => { console.log("click"); setFetchTrigger(true); }}
                    >
                        Next
                    </button>
                </div>
            </main>
        </div>
    );
}
