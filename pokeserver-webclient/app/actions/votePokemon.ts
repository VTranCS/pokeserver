"use server"
import { db } from '@vercel/postgres';

export async function votePokemon(pokemonId: number, action: "up" | "down"): Promise<number> {
    const client = await db.connect()
    const res = await client.sql`UPDATE pokemon SET vote = vote + ${action === "up" ? 1 : -1} WHERE id = ${pokemonId} RETURNING vote`;
    if (res.rows.length === 0) {
        throw new Error('Failed to update vote');
    }
    client.release();
    return res.rows[0].vote;
}