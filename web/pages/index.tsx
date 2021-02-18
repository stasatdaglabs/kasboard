import Head from 'next/head'
import styles from '../styles/Home.module.css'
import {GetServerSideProps} from "next";
import {Client} from "pg"

const Home = ({now}) => {
    return (
        <div className={styles.container}>
            <Head>
                <title>Kaspa Dashboard</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            Hello world!
            <br/>
            Blue score: {now}
        </div>
    )
}

export default Home;

export const getServerSideProps: GetServerSideProps = async () => {
    const client = new Client();
    await client.connect();
    const result = await client.query("SELECT MAX(blue_score) AS blue_score FROM blocks");
    await client.end();

    const now = result.rows[0].blue_score;

    return {
        props: {
            now: now
        }
    }
}
