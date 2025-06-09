import React, { useState } from "react";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Doughnut } from "react-chartjs-2";
import trendService, { SentimentResponse } from "../services/trend_service";
import { Keyword } from "../services/keyword_service";

// Chart.jsのコンポーネントを登録
ChartJS.register(ArcElement, Tooltip, Legend);

interface SentimentAnalysisProps {
  keyword: Keyword;
}

const SentimentAnalysis: React.FC<SentimentAnalysisProps> = ({ keyword }) => {
  const [selectedDate, setSelectedDate] = useState<string>(() => {
    const today = new Date();
    return today.toISOString().split("T")[0];
  });
  const [sentimentData, setSentimentData] = useState<SentimentResponse | null>(
    null
  );
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadSentimentData = async () => {
    setLoading(true);
    setError(null);

    try {
      // 過去30日間のセンチメント分析を実行
      const response = await trendService.getSentimentDetail({
        keyword_id: keyword.id,
        period: 30,
      });

      // レスポンス形式を変換
      const convertedResponse = {
        positive:
          response.positive_count /
          (response.positive_count +
            response.negative_count +
            response.neutral_count),
        neutral:
          response.neutral_count /
          (response.positive_count +
            response.negative_count +
            response.neutral_count),
        negative:
          response.negative_count /
          (response.positive_count +
            response.negative_count +
            response.neutral_count),
      };

      setSentimentData(convertedResponse);
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error ||
        "センチメントデータの読み込みに失敗しました";
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // チャート用のデータを準備
  const chartData = sentimentData
    ? {
        labels: ["ポジティブ", "ニュートラル", "ネガティブ"],
        datasets: [
          {
            data: [
              Math.round(sentimentData.positive * 100),
              Math.round(sentimentData.neutral * 100),
              Math.round(sentimentData.negative * 100),
            ],
            backgroundColor: [
              "#10B981", // Green for positive
              "#6B7280", // Gray for neutral
              "#EF4444", // Red for negative
            ],
            borderColor: ["#059669", "#4B5563", "#DC2626"],
            borderWidth: 2,
          },
        ],
      }
    : null;

  const chartOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: "bottom" as const,
      },
      title: {
        display: true,
        text: `${keyword.keyword} のセンチメント分析`,
      },
      tooltip: {
        callbacks: {
          label: function (context: any) {
            return `${context.label}: ${context.parsed}%`;
          },
        },
      },
    },
  };

  return (
    <div className="sentiment-analysis">
      <div className="sentiment-header">
        <h3>センチメント分析</h3>
      </div>

      {/* 日付選択とロードボタン */}
      <div className="sentiment-controls">
        <label>
          分析日:
          <input
            type="date"
            value={selectedDate}
            onChange={(e) => setSelectedDate(e.target.value)}
            max={new Date().toISOString().split("T")[0]}
          />
        </label>
        <button onClick={loadSentimentData} disabled={loading || !selectedDate}>
          {loading ? "分析中..." : "分析実行"}
        </button>
      </div>

      {/* エラー表示 */}
      {error && <div className="error-message">{error}</div>}

      {/* ローディング表示 */}
      {loading && (
        <div className="sentiment-loading">
          <div className="spinner"></div>
          <p>センチメント分析を実行中...</p>
        </div>
      )}

      {/* センチメントデータ表示 */}
      {sentimentData && !loading && (
        <div className="sentiment-results">
          {/* ドーナツチャート */}
          <div className="sentiment-chart">
            {chartData && <Doughnut data={chartData} options={chartOptions} />}
          </div>

          {/* 詳細データ */}
          <div className="sentiment-details">
            <div className="sentiment-item positive">
              <div className="sentiment-label">
                <span className="sentiment-icon">😊</span>
                <span>ポジティブ</span>
              </div>
              <div className="sentiment-value">
                {(sentimentData.positive * 100).toFixed(1)}%
              </div>
            </div>

            <div className="sentiment-item neutral">
              <div className="sentiment-label">
                <span className="sentiment-icon">😐</span>
                <span>ニュートラル</span>
              </div>
              <div className="sentiment-value">
                {(sentimentData.neutral * 100).toFixed(1)}%
              </div>
            </div>

            <div className="sentiment-item negative">
              <div className="sentiment-label">
                <span className="sentiment-icon">😔</span>
                <span>ネガティブ</span>
              </div>
              <div className="sentiment-value">
                {(sentimentData.negative * 100).toFixed(1)}%
              </div>
            </div>
          </div>

          {/* 分析サマリー */}
          <div className="sentiment-summary">
            <h4>分析サマリー</h4>
            <p>
              {selectedDate}における「{keyword.keyword}
              」のセンチメント分析結果です。
              {sentimentData.positive > 0.5 ? (
                <span className="summary-positive">
                  全体的にポジティブな反応が多く見られます。
                </span>
              ) : sentimentData.negative > 0.5 ? (
                <span className="summary-negative">
                  ネガティブな反応が多く見られます。
                </span>
              ) : (
                <span className="summary-neutral">
                  ニュートラルな反応が中心的です。
                </span>
              )}
            </p>
          </div>
        </div>
      )}

      {/* データがない場合の表示 */}
      {!sentimentData && !loading && !error && (
        <div className="sentiment-empty">
          <p>日付を選択して「分析実行」ボタンをクリックしてください。</p>
        </div>
      )}
    </div>
  );
};

export default SentimentAnalysis;
