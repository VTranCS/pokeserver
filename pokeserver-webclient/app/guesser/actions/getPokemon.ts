"use server";

export interface Pokemon {
  name: string;
  sprites: {
    front_default: string;
    other: {
      "official-artwork": {
        front_default: string;
      };
    };
  };
}

export async function getRandomPokemon(): Promise<Pokemon> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_POKEMON_BASE_URL}/getpokemon`, {
    cache: 'no-store',
  });

  if (!res.ok) throw new Error('Failed to fetch Pokémon');

  return res.json();
}
