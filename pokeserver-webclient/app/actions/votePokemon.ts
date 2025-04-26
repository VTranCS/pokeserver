"use server"

export async function votePokemon(pokemonId: number, action: "up" | "down"): Promise<number> {
    const url = `${process.env.POKEMON_BASE_URL}/${process.env.SUFFIX_VOTE_POKEMON}?id=${pokemonId}&vote=${action}`;
    const response = await fetch(url);
    const data = await response.json();
    return data.Vote;
}