import { getRandomPokemon} from "./actions/getPokemon";import PokemonClient from './PokemonClient';

export default async function HomePage() {
  const pokemon = await getRandomPokemon();

  return (
    <div className="p-4">
      <PokemonClient initialPokemon={pokemon} />
    </div>
  );
}
