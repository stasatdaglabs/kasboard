import {Bar} from "react-chartjs-2";
import {BlueScoreOverTimeData} from "../scripts/database";
import {timestampToString} from "../scripts/date";

export type BlueScoreOverTimeChartProps = {
    blueScoreOverTime: BlueScoreOverTimeData,
}

const BlueScoreOverTime = ({blueScoreOverTime}: BlueScoreOverTimeChartProps) => {
    const blueScoreOverTimeData = {
        labels: blueScoreOverTime.map(item => timestampToString(item.timestamp)),
        datasets: [
            {
                data: blueScoreOverTime.map(item => item.blueScore),
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
        <Bar data={blueScoreOverTimeData} options={blueScoreOverTimeOptions}/>
    );
};

export default BlueScoreOverTime;
