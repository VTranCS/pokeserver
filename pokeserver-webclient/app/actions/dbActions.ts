import { db } from '@vercel/postgres';
import { FetchedPokemon } from '../types/PokemonDBEntry';


export async function seedPokemon() {
    const generation = process.env.POKEMON_GENERATION ? parseInt(process.env.POKEMON_GENERATION, 10) : 1;
    const client = await db.connect();
    await client.sql`
      CREATE TABLE IF NOT EXISTS pokemon (
        id INT PRIMARY KEY,
        generation INT NOT NULL,
        name VARCHAR(255) NOT NULL,
        image TEXT NOT NULL,
        vote INT NOT NULL
      );
    `;

    const generationExists = await client.sql`
        SELECT EXISTS(SELECT 1 FROM pokemon WHERE generation = ${generation});
    `;

    if (generationExists.rows[0]?.exists) {
        console.log(`Generation ${generation} already exists in the database.`);
        return;
    }

    const pokemonGenerationURL = `https://pokeapi.co/api/v2/generation/${generation}`;
    const response = await fetch(pokemonGenerationURL);
    const data = await response.json();

    console.log(`Pulling Pokémon for generation ${generation}...`);

    const values: string[] = data.pokemon_species.map((pokemon: FetchedPokemon) => {
        console.log(pokemon);
        const urlParts = pokemon.url.split('/').filter(Boolean);
        const id = Number(urlParts[urlParts.length - 1]);
        const sprite = `https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/${id}.png`;

        // Escape quotes for name and sprite
        const safeName = pokemon.name.replace(/'/g, "''");
        const safeSprite = sprite.replace(/'/g, "''");

        return `(${id}, '${safeName}', ${generation}, '${safeSprite}', 0)`;
    });

    if (values.length > 0) {
        const insertQuery = `
            INSERT INTO pokemon (id, name, generation, image, vote)
            VALUES ${values.join(', ')}
            ON CONFLICT (id) DO NOTHING;
        `;
        console.log(insertQuery)
        await client.query(insertQuery);
    }

    console.log(`✅ Successfully seeded generation ${generation}!`);
}

