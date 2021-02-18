import Head from 'next/head';
import styles from '../styles/Home.module.css';
import {GetServerSideProps} from "next";
import {getGreatestBlueScore} from "../scripts/database";

const Home = ({greatestBlueScore}: HomeProps) => {
    return (
        <div className={styles.container}>
            <Head>
                <title>Kaspa Dashboard</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            Greatest blue score: {greatestBlueScore}
        </div>
    )
};

export default Home;

type HomeProps = {
    greatestBlueScore: Number,
};

export const getServerSideProps: GetServerSideProps = async () => {
    const greatestBlueScore = await getGreatestBlueScore();

    return {
        props: {
            greatestBlueScore: greatestBlueScore
        }
    }
};
