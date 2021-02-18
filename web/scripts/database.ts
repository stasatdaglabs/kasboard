import {Client} from "pg"

export async function getGreatestBlueScore() {
    const client = new Client();
    await client.connect();
    const result = await client.query("SELECT MAX(blue_score) AS blue_score FROM blocks");
    await client.end();

    return result.rows[0].blue_score;
}