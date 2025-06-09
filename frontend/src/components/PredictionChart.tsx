import React from "react";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import { Line } from "react-chartjs-2";
import { PredictionData } from "../services/trend_service";

// Chart.jsのコンポーネントを登録
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

interface PredictionChartProps {
  data: PredictionData[];
  loading: boolean;
}

const PredictionChart: React.FC<PredictionChartProps> = ({ data, loading }) => {
  // データが空の場合
  if (loading) {
    return (
      <div className="chart-loading">
        <div className="spinner"></div>
        <p>予測チャートを読み込み中...</p>
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="chart-empty">
        <p>予測データがありません</p>
      </div>
    );
  }

  // データを日付順にソート
  const sortedData = [...data].sort(
    (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()
  );

  // チャート用のデータを準備
  const chartData = {
    labels: sortedData.map((prediction) => {
      const date = new Date(prediction.date);
      return date.toLocaleDateString("ja-JP", {
        month: "short",
        day: "numeric",
      });
    }),
    datasets: [
      {
        label: "予測ボリューム",
        data: sortedData.map((prediction) => prediction.volume),
        borderColor: "rgb(153, 102, 255)",
        backgroundColor: "rgba(153, 102, 255, 0.2)",
        borderDash: [5, 5], // 点線で予測データを表現
        tension: 0.1,
        pointStyle: "triangle",
      },
    ],
  };

  // チャートのオプション
  const options = {
    responsive: true,
    interaction: {
      mode: "index" as const,
      intersect: false,
    },
    scales: {
      x: {
        display: true,
        title: {
          display: true,
          text: "予測日",
        },
      },
      y: {
        display: true,
        title: {
          display: true,
          text: "予測ボリューム",
        },
        beginAtZero: true,
      },
    },
    plugins: {
      legend: {
        position: "top" as const,
      },
      title: {
        display: true,
        text: "トレンド予測",
      },
      tooltip: {
        callbacks: {
          label: function (context: any) {
            return `予測ボリューム: ${context.parsed.y}`;
          },
          afterLabel: function () {
            return "※ これは予測値です";
          },
        },
      },
    },
  };

  return (
    <div className="prediction-chart">
      <div className="chart-container">
        <Line data={chartData} options={options} />
      </div>

      {/* 予測サマリー */}
      <div className="chart-summary">
        <div className="summary-item">
          <span className="label">予測期間:</span>
          <span className="value">{data.length}日間</span>
        </div>
        <div className="summary-item">
          <span className="label">平均予測値:</span>
          <span className="value">
            {Math.round(
              data.reduce((sum, prediction) => sum + prediction.volume, 0) /
                data.length
            )}
          </span>
        </div>
        <div className="summary-item">
          <span className="label">最大予測値:</span>
          <span className="value">
            {Math.max(...data.map((p) => p.volume))}
          </span>
        </div>
        <div className="summary-item">
          <span className="label">最小予測値:</span>
          <span className="value">
            {Math.min(...data.map((p) => p.volume))}
          </span>
        </div>
      </div>

      {/* 注意書き */}
      <div className="prediction-note">
        <p>
          ⚠️
          予測データは過去のトレンドに基づく推定値です。実際の結果とは異なる場合があります。
        </p>
      </div>
    </div>
  );
};

export default PredictionChart;
