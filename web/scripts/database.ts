import {Client} from "pg"

export type BlueScoreOverTimeData = [{
    blue_score: number,
    timestamp: number,
}];

export async function getBlueScoreOverTime(): Promise<BlueScoreOverTimeData> {
    const client = new Client();
    await client.connect();
    const result = await client.query("SELECT blue_score, timestamp FROM blocks ORDER BY timestamp DESC LIMIT 100");
    await client.end();

    return result.rows.reverse();
}
