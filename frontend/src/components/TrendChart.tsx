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
  TimeScale,
} from "chart.js";
import { Line } from "react-chartjs-2";
import { TrendRecord } from "../services/trend_service";

// Chart.jsのコンポーネントを登録
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale
);

interface TrendChartProps {
  data: TrendRecord[];
  loading: boolean;
}

const TrendChart: React.FC<TrendChartProps> = ({ data, loading }) => {
  // データが空の場合
  if (loading) {
    return (
      <div className="chart-loading">
        <div className="spinner"></div>
        <p>チャートを読み込み中...</p>
      </div>
    );
  }

  // dataがnullまたはundefinedの場合
  if (!data) {
    return (
      <div className="chart-empty">
        <p>データの読み込みに失敗しました</p>
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="chart-empty">
        <p>表示するデータがありません</p>
      </div>
    );
  }

  // データを日付順にソート
  const sortedData = [...data].sort(
    (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()
  );

  // チャート用のデータを準備
  const chartData = {
    labels: sortedData.map((record) => {
      const date = new Date(record.date);
      return date.toLocaleDateString("ja-JP", {
        month: "short",
        day: "numeric",
      });
    }),
    datasets: [
      {
        label: "ボリューム",
        data: sortedData.map((record) => record.volume),
        borderColor: "rgb(75, 192, 192)",
        backgroundColor: "rgba(75, 192, 192, 0.2)",
        yAxisID: "y",
        tension: 0.1,
      },
      {
        label: "センチメント",
        data: sortedData.map((record) => record.sentiment),
        borderColor: "rgb(255, 99, 132)",
        backgroundColor: "rgba(255, 99, 132, 0.2)",
        yAxisID: "y1",
        tension: 0.1,
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
          text: "日付",
        },
      },
      y: {
        type: "linear" as const,
        display: true,
        position: "left" as const,
        title: {
          display: true,
          text: "ボリューム",
        },
        beginAtZero: true,
      },
      y1: {
        type: "linear" as const,
        display: true,
        position: "right" as const,
        title: {
          display: true,
          text: "センチメント",
        },
        min: -1,
        max: 1,
        grid: {
          drawOnChartArea: false,
        },
      },
    },
    plugins: {
      legend: {
        position: "top" as const,
      },
      title: {
        display: true,
        text: "トレンドデータ推移",
      },
      tooltip: {
        callbacks: {
          afterLabel: function (context: any) {
            if (context.datasetIndex === 1) {
              // センチメントの場合、詳細な説明を追加
              const sentiment = context.parsed.y;
              if (sentiment > 0.3) return "ポジティブ";
              if (sentiment < -0.3) return "ネガティブ";
              return "ニュートラル";
            }
            return "";
          },
        },
      },
    },
  };

  return (
    <div className="trend-chart">
      <div className="chart-container">
        <Line data={chartData} options={options} />
      </div>

      {/* データサマリー */}
      {data && data.length > 0 && (
        <div className="chart-summary">
          <div className="summary-item">
            <span className="label">総データ数:</span>
            <span className="value">{data.length}件</span>
          </div>
          <div className="summary-item">
            <span className="label">平均ボリューム:</span>
            <span className="value">
              {Math.round(
                data.reduce((sum, record) => sum + record.volume, 0) /
                  data.length
              )}
            </span>
          </div>
          <div className="summary-item">
            <span className="label">平均センチメント:</span>
            <span className="value">
              {(
                data.reduce((sum, record) => sum + record.sentiment, 0) /
                data.length
              ).toFixed(2)}
            </span>
          </div>
        </div>
      )}
    </div>
  );
};

export default TrendChart;
