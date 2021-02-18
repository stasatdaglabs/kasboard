import {Client} from "pg"

export async function getBlueScoreOverTime() {
    const client = new Client();
    await client.connect();
    const result = await client.query("SELECT blue_score, timestamp FROM blocks ORDER BY timestamp DESC LIMIT 100");
    await client.end();

    return result.rows.reverse();
}
