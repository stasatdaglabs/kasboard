import '../styles/globals.css';
import {AppProps} from "next/app";

const Kashboard = ({Component, pageProps}: AppProps) => {
    return <Component {...pageProps} />
};

export default Kashboard;
