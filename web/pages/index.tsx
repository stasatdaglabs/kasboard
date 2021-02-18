import Head from 'next/head';
import styles from '../styles/Home.module.css';
import {GetServerSideProps} from "next";
import {BlueScoreOverTimeData, getBlueScoreOverTime} from "../scripts/database";
import BlueScoreOverTime from "../components/BlueScoreOverTime";

type HomeProps = {
    blueScoreOverTime: BlueScoreOverTimeData,
};

const Home = ({blueScoreOverTime}: HomeProps) => {
    return (
        <div className={styles.container}>
            <Head>
                <title>Kaspa Dashboard</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            <div className={styles.blueScoreOverTime}>
                <BlueScoreOverTime blueScoreOverTime={blueScoreOverTime}/>
            </div>
        </div>
    )
};

export default Home;

export const getServerSideProps: GetServerSideProps = async () => {
    const blueScoreOverTime = await getBlueScoreOverTime();

    return {
        props: {
            blueScoreOverTime,
        }
    }
};
