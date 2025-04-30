import { getRandomPokemon} from "../actions/getPokemon";import PokemonClient from './PokemonClient';

export default async function HomePage() {
  const generation = process.env.POKEMON_GENERATION ? parseInt(process.env.POKEMON_GENERATION, 10) : 1
  const pokemon = await getRandomPokemon(generation)

  return (
    <div className="p-4">
      <PokemonClient
        initialPokemon={pokemon}
        generation={generation} />
    </div>
  );
}
