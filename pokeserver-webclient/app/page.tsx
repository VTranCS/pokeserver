"use client";

import React from "react";
import Link from "next/link";
export default function PokemonVoteDisplay({}) {


  const buttonBase =
  "rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center bg-purple-500 hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:min-w-44";

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center">
        <img src="https://upload.wikimedia.org/wikipedia/commons/5/53/Pok%C3%A9_Ball_icon.svg" alt="Example SVG" width={360} height={360} />
        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <Link href="/guesser" className={buttonBase}>
            Guesser
          </Link>
          <Link href="/voter" className={buttonBase}>
            Voter
          </Link>
          <Link href="/collection" className={buttonBase}>
            Collection
          </Link>
        </div>
      </main>
    </div>
  );
}
