import React, { useState } from "react";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Doughnut } from "react-chartjs-2";
import trendService, { SentimentResponse } from "../services/trend_service";
import { Keyword } from "../services/keyword_service";
import { Input } from "./ui/Input";
import { Button } from "./ui/Button";
import { Calendar } from "lucide-react";

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
        labels: ["好意的", "普通", "否定的"],
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
        text: `${keyword.keyword} の評判・印象分析`,
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
    <div className="sentiment-analysis bg-white dark:bg-slate-800 rounded-lg p-6 border border-gray-200 dark:border-slate-700">
      <div className="sentiment-header mb-6">
        <h3 className="text-xl font-semibold text-gray-900 dark:text-white flex items-center space-x-2">
          <Calendar className="w-5 h-5" />
          <span>評判・印象分析</span>
        </h3>
      </div>

      {/* 日付選択とロードボタン */}
      <div className="sentiment-controls mb-6 flex flex-col sm:flex-row gap-4 items-end">
        <div className="flex-1 space-y-2">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            分析日:
          </label>
          <Input
            type="date"
            value={selectedDate}
            onChange={(e) => setSelectedDate(e.target.value)}
            max={new Date().toISOString().split("T")[0]}
            className="w-full"
          />
        </div>
        <Button
          onClick={loadSentimentData}
          disabled={loading || !selectedDate}
          className="bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white"
        >
          {loading ? "分析中..." : "分析実行"}
        </Button>
      </div>

      {/* エラー表示 */}
      {error && (
        <div className="error-message mb-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-700 dark:text-red-300">
          {error}
        </div>
      )}

      {/* ローディング表示 */}
      {loading && (
        <div className="sentiment-loading flex flex-col items-center justify-center py-8">
          <div className="spinner w-8 h-8 border-4 border-purple-200 border-t-purple-600 rounded-full animate-spin mb-4"></div>
          <p className="text-gray-600 dark:text-gray-400">
            評判・印象を分析中...
          </p>
        </div>
      )}

      {/* センチメントデータ表示 */}
      {sentimentData && !loading && (
        <div className="sentiment-results space-y-6">
          {/* ドーナツチャート */}
          <div className="sentiment-chart bg-white dark:bg-slate-700 p-6 rounded-lg border border-gray-200 dark:border-slate-600">
            {chartData && <Doughnut data={chartData} options={chartOptions} />}
          </div>

          {/* 詳細データ */}
          <div className="sentiment-details grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="sentiment-item bg-green-50 dark:bg-green-900/20 p-4 rounded-lg border border-green-200 dark:border-green-800">
              <div className="sentiment-label flex items-center space-x-2 mb-2">
                <span className="sentiment-icon text-2xl">😊</span>
                <span className="font-medium text-green-800 dark:text-green-300">
                  好意的
                </span>
              </div>
              <div className="sentiment-value text-2xl font-bold text-green-600 dark:text-green-400">
                {(sentimentData.positive * 100).toFixed(1)}%
              </div>
            </div>

            <div className="sentiment-item bg-gray-50 dark:bg-gray-700/50 p-4 rounded-lg border border-gray-200 dark:border-gray-600">
              <div className="sentiment-label flex items-center space-x-2 mb-2">
                <span className="sentiment-icon text-2xl">😐</span>
                <span className="font-medium text-gray-800 dark:text-gray-300">
                  普通
                </span>
              </div>
              <div className="sentiment-value text-2xl font-bold text-gray-600 dark:text-gray-400">
                {(sentimentData.neutral * 100).toFixed(1)}%
              </div>
            </div>

            <div className="sentiment-item bg-red-50 dark:bg-red-900/20 p-4 rounded-lg border border-red-200 dark:border-red-800">
              <div className="sentiment-label flex items-center space-x-2 mb-2">
                <span className="sentiment-icon text-2xl">😔</span>
                <span className="font-medium text-red-800 dark:text-red-300">
                  否定的
                </span>
              </div>
              <div className="sentiment-value text-2xl font-bold text-red-600 dark:text-red-400">
                {(sentimentData.negative * 100).toFixed(1)}%
              </div>
            </div>
          </div>

          {/* 分析サマリー */}
          <div className="sentiment-summary bg-blue-50 dark:bg-blue-900/20 p-6 rounded-lg border border-blue-200 dark:border-blue-800">
            <h4 className="text-lg font-semibold text-blue-900 dark:text-blue-300 mb-3">
              分析結果
            </h4>
            <p className="text-gray-700 dark:text-gray-300 leading-relaxed">
              {selectedDate}における「
              <strong className="text-purple-600 dark:text-purple-400">
                {keyword.keyword}
              </strong>
              」の評判・印象分析結果です。
              {sentimentData.positive > 0.5 ? (
                <span className="summary-positive text-green-600 dark:text-green-400 font-medium">
                  全体的に好意的な反応が多く見られます。
                </span>
              ) : sentimentData.negative > 0.5 ? (
                <span className="summary-negative text-red-600 dark:text-red-400 font-medium">
                  否定的な反応が多く見られます。
                </span>
              ) : (
                <span className="summary-neutral text-gray-600 dark:text-gray-400 font-medium">
                  普通の反応が中心的です。
                </span>
              )}
            </p>
          </div>
        </div>
      )}

      {/* データがない場合の表示 */}
      {!sentimentData && !loading && !error && (
        <div className="sentiment-empty text-center py-12 text-gray-500 dark:text-gray-400">
          <Calendar className="w-12 h-12 mx-auto mb-4 opacity-50" />
          <p className="text-lg">
            日付を選択して「分析実行」ボタンをクリックしてください。
          </p>
        </div>
      )}
    </div>
  );
};

export default SentimentAnalysis;
