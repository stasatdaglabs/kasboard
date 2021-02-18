import {Bar} from "react-chartjs-2";
import {BlueScoreOverTimeData} from "../scripts/database";

export type BlueScoreOverTimeChartProps = {
    blueScoreOverTime: BlueScoreOverTimeData,
}

const BlueScoreOverTime = ({blueScoreOverTime}: BlueScoreOverTimeChartProps) => {
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
        tooltips: {
            enabled: false,
        }
    };

    return (
        <Bar data={blueScoreOverTimeData} options={blueScoreOverTimeOptions}/>
    );
};

export default BlueScoreOverTime;
