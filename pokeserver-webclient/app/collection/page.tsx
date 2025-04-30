
import { getAllPokemonDBEntry } from "../actions/getPokemon";
import PokemonDBList from "./PokemonDBList";
export default async function HomePage() {
  const generation = process.env.POKEMON_GENERATION ? parseInt(process.env.POKEMON_GENERATION, 10) : 1
  const pokemon = await getAllPokemonDBEntry(generation);

  return (
    <div className="p-4">
      <PokemonDBList pokemonList={pokemon} />
    </div>
  );
}
