import Head from 'next/head'
import styles from '../styles/Home.module.css'
import {GetServerSideProps} from "next";

export default ({prop}) => {
    return (
        <div className={styles.container}>
            <Head>
                <title>Kaspa Dashboard</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            Hello world!
            <br />
            {prop}
        </div>
    )
}

export const getServerSideProps: GetServerSideProps = async () => {
    return {
        props: {
            prop: "Hello from getServerSideProps"
        }
    }
}
