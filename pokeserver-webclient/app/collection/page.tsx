
import { getAllPokemonDBEntry } from "../actions/getPokemon";
import PokemonDBList from "./PokemonDBList";
export default async function HomePage() {
  const pokemon = await getAllPokemonDBEntry();

  return (
    <div className="p-4">
      <PokemonDBList pokemonList={pokemon} />
    </div>
  );
}
