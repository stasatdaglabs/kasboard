import Head from 'next/head';
import styles from '../styles/Home.module.css';
import {GetServerSideProps} from "next";
import {getBlueScoreOverTime} from "../scripts/database";
import {Bar} from "react-chartjs-2";

const Home = ({blueScoreOverTime}: HomeProps) => {
    const blueScoreOverTimeData = {
        labels: blueScoreOverTime.map(item => new Date(item.timestamp / 1000)),
        datasets: [
            {
                data: blueScoreOverTime.map(item => item.blue_score),
            },
        ],
    };

    const blueScoreOverTimeOptions = {
        legend: {
            display: false,
        },
        scales: {
            xAxes: [
                {
                    display: false,
                },
            ],
        },
    };

    return (
        <div className={styles.container}>
            <Head>
                <title>Kaspa Dashboard</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            <Bar data={blueScoreOverTimeData} options={blueScoreOverTimeOptions}/>
        </div>
    )
};

export default Home;

type HomeProps = {
    blueScoreOverTime: [{
        blue_score: number,
        timestamp: number,
    }],
};

export const getServerSideProps: GetServerSideProps = async () => {
    const blueScoreOverTime = await getBlueScoreOverTime();

    return {
        props: {
            blueScoreOverTime,
        }
    }
};
