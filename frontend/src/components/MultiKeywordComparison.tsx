import React, { useState, useEffect } from "react";
import keywordService, { Keyword } from "../services/keyword_service";
import trendService, { TrendRecord } from "../services/trend_service";
import { Input } from "./ui/Input";
import { Button } from "./ui/Button";
import { Calendar, TrendingUp } from "lucide-react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  BarChart,
  Bar,
} from "recharts";

interface KeywordComparisonData {
  keyword: Keyword;
  trendData: TrendRecord[];
  totalVolume: number;
  avgSentiment: number;
  growth: number;
}

interface ComparisonChartData {
  date: string;
  [key: string]: string | number; // Dynamic keys for each keyword
}

const MultiKeywordComparison: React.FC = () => {
  const [availableKeywords, setAvailableKeywords] = useState<Keyword[]>([]);
  const [selectedKeywords, setSelectedKeywords] = useState<number[]>([]);
  const [comparisonData, setComparisonData] = useState<KeywordComparisonData[]>(
    []
  );
  const [chartData, setChartData] = useState<ComparisonChartData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 日付範囲
  const [dateFrom, setDateFrom] = useState<string>(() => {
    const date = new Date();
    date.setDate(date.getDate() - 30);
    return date.toISOString().split("T")[0];
  });

  const [dateTo, setDateTo] = useState<string>(() => {
    const date = new Date();
    return date.toISOString().split("T")[0];
  });

  // ビューモード
  const [viewMode, setViewMode] = useState<"volume" | "sentiment" | "combined">(
    "volume"
  );

  // キーワード一覧の読み込み
  useEffect(() => {
    loadAvailableKeywords();
  }, []);

  // 選択されたキーワードの比較データ読み込み
  useEffect(() => {
    if (selectedKeywords.length > 0) {
      loadComparisonData();
    }
  }, [selectedKeywords, dateFrom, dateTo]);

  const loadAvailableKeywords = async () => {
    try {
      setLoading(true);
      const response = await keywordService.getKeywords();
      setAvailableKeywords(response.keywords);
    } catch (err) {
      setError("キーワードの読み込みに失敗しました");
      console.error("Failed to load keywords:", err);
    } finally {
      setLoading(false);
    }
  };

  const loadComparisonData = async () => {
    try {
      setLoading(true);
      setError(null);

      const comparisonPromises = selectedKeywords.map(async (keywordId) => {
        const keyword = availableKeywords.find((k) => k.id === keywordId);
        if (!keyword) return null;

        try {
          const response = await trendService.getTrendData({
            q: keywordId,
            from: dateFrom,
            to: dateTo,
          });
          const trendData = response.records;

          // メトリクス計算
          const totalVolume = trendData.reduce(
            (sum: number, record: TrendRecord) => sum + record.volume,
            0
          );
          const avgSentiment =
            trendData.length > 0
              ? trendData.reduce(
                  (sum: number, record: TrendRecord) => sum + record.sentiment,
                  0
                ) / trendData.length
              : 0;

          // 成長率計算（最初と最後のデータを比較）
          const growth =
            trendData.length >= 2
              ? ((trendData[trendData.length - 1].volume -
                  trendData[0].volume) /
                  trendData[0].volume) *
                100
              : 0;

          return {
            keyword,
            trendData,
            totalVolume,
            avgSentiment,
            growth,
          };
        } catch (err) {
          console.error(
            `Failed to load data for keyword ${keyword.keyword}:`,
            err
          );
          return {
            keyword,
            trendData: [],
            totalVolume: 0,
            avgSentiment: 0,
            growth: 0,
          };
        }
      });

      const results = await Promise.all(comparisonPromises);
      const validResults = results.filter(
        (result) => result !== null
      ) as KeywordComparisonData[];

      setComparisonData(validResults);
      generateChartData(validResults);
    } catch (err) {
      setError("比較データの読み込みに失敗しました");
      console.error("Failed to load comparison data:", err);
    } finally {
      setLoading(false);
    }
  };

  const generateChartData = (data: KeywordComparisonData[]) => {
    // 全ての日付を収集
    const allDates = new Set<string>();
    data.forEach((item) => {
      item.trendData.forEach((record) => {
        allDates.add(new Date(record.date).toISOString().split("T")[0]);
      });
    });

    // チャート用データを生成
    const chartData: ComparisonChartData[] = Array.from(allDates)
      .sort()
      .map((date) => {
        const chartPoint: ComparisonChartData = { date };

        data.forEach((item) => {
          const record = item.trendData.find(
            (r) => new Date(r.date).toISOString().split("T")[0] === date
          );

          chartPoint[`${item.keyword.keyword}_volume`] = record?.volume || 0;
          chartPoint[`${item.keyword.keyword}_sentiment`] =
            record?.sentiment || 0;
        });

        return chartPoint;
      });

    setChartData(chartData);
  };

  const handleKeywordToggle = (keywordId: number) => {
    setSelectedKeywords((prev) => {
      if (prev.includes(keywordId)) {
        return prev.filter((id) => id !== keywordId);
      } else if (prev.length < 5) {
        // 最大5キーワードまで
        return [...prev, keywordId];
      }
      return prev;
    });
  };

  const getKeywordColors = () => {
    const colors = ["#8884d8", "#82ca9d", "#ffc658", "#ff7c7c", "#8dd1e1"];
    return colors;
  };

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString("ja-JP", {
      month: "short",
      day: "numeric",
    });
  };

  if (loading && comparisonData.length === 0) {
    return <div className="loading">読み込み中...</div>;
  }

  return (
    <div className="multi-keyword-comparison bg-white dark:bg-slate-800 rounded-lg p-6 border border-gray-200 dark:border-slate-700">
      <div className="comparison-header mb-6">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white flex items-center space-x-2">
          <TrendingUp className="w-6 h-6" />
          <span>キーワード比較分析</span>
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mt-2">
          複数のキーワードの人気度を比較して、ファッション業界の動向を把握できます
        </p>
      </div>

      {/* コントロールパネル */}
      <div className="comparison-controls mb-8 space-y-6">
        <div className="date-controls grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center space-x-2">
              <Calendar className="w-4 h-4" />
              <span>開始日:</span>
            </label>
            <Input
              type="date"
              value={dateFrom}
              onChange={(e) => setDateFrom(e.target.value)}
              className="w-full"
            />
          </div>
          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center space-x-2">
              <Calendar className="w-4 h-4" />
              <span>終了日:</span>
            </label>
            <Input
              type="date"
              value={dateTo}
              onChange={(e) => setDateTo(e.target.value)}
              className="w-full"
            />
          </div>
        </div>

        <div className="view-controls space-y-2">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            表示モード:
          </label>
          <select
            value={viewMode}
            onChange={(e) =>
              setViewMode(e.target.value as "volume" | "sentiment" | "combined")
            }
            className="w-full md:w-auto px-3 py-2 text-sm bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 dark:bg-slate-700 dark:border-slate-600 dark:text-white"
          >
            <option value="volume">話題量</option>
            <option value="sentiment">評判</option>
            <option value="combined">両方表示</option>
          </select>
        </div>
      </div>

      {/* キーワード選択 */}
      <div className="keyword-selection mb-8">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          比較するキーワードを選択（最大5つ）
        </h3>
        <div className="keyword-list grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3">
          {availableKeywords.map((keyword) => (
            <label
              key={keyword.id}
              className="keyword-checkbox flex items-center space-x-2 p-3 bg-gray-50 dark:bg-slate-700 rounded-lg border border-gray-200 dark:border-slate-600 hover:bg-gray-100 dark:hover:bg-slate-600 cursor-pointer transition-colors"
            >
              <input
                type="checkbox"
                checked={selectedKeywords.includes(keyword.id)}
                onChange={() => handleKeywordToggle(keyword.id)}
                disabled={
                  !selectedKeywords.includes(keyword.id) &&
                  selectedKeywords.length >= 5
                }
                className="rounded border-gray-300 dark:border-slate-500 text-purple-600 focus:ring-purple-500 dark:bg-slate-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">
                {keyword.keyword}
              </span>
            </label>
          ))}
        </div>
        <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
          比較したいキーワードを選択してください
        </p>
      </div>

      {error && (
        <div className="error-message mb-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-700 dark:text-red-300">
          {error}
        </div>
      )}

      {/* 比較チャート */}
      {selectedKeywords.length > 0 && chartData.length > 0 && (
        <div className="comparison-charts space-y-8">
          {(viewMode === "volume" || viewMode === "combined") && (
            <div className="chart-container bg-white dark:bg-slate-700 p-6 rounded-lg border border-gray-200 dark:border-slate-600">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                話題量の比較
              </h3>
              <ResponsiveContainer width="100%" height={400}>
                <LineChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                  <XAxis
                    dataKey="date"
                    tickFormatter={formatDate}
                    className="text-gray-600 dark:text-gray-300"
                  />
                  <YAxis className="text-gray-600 dark:text-gray-300" />
                  <Tooltip
                    labelFormatter={(label) => formatDate(label as string)}
                    contentStyle={{
                      backgroundColor: "var(--background)",
                      border: "1px solid var(--border)",
                      borderRadius: "8px",
                      color: "var(--foreground)",
                    }}
                  />
                  <Legend />
                  {comparisonData.map((item, index) => (
                    <Line
                      key={`${item.keyword.keyword}_volume`}
                      type="monotone"
                      dataKey={`${item.keyword.keyword}_volume`}
                      stroke={getKeywordColors()[index]}
                      strokeWidth={2}
                      name={`${item.keyword.keyword} (話題量)`}
                    />
                  ))}
                </LineChart>
              </ResponsiveContainer>
            </div>
          )}

          {(viewMode === "sentiment" || viewMode === "combined") && (
            <div className="chart-container bg-white dark:bg-slate-700 p-6 rounded-lg border border-gray-200 dark:border-slate-600">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                評判の比較
              </h3>
              <ResponsiveContainer width="100%" height={400}>
                <LineChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                  <XAxis
                    dataKey="date"
                    tickFormatter={formatDate}
                    className="text-gray-600 dark:text-gray-300"
                  />
                  <YAxis
                    domain={[0, 1]}
                    className="text-gray-600 dark:text-gray-300"
                  />
                  <Tooltip
                    labelFormatter={(label) => formatDate(label as string)}
                    contentStyle={{
                      backgroundColor: "var(--background)",
                      border: "1px solid var(--border)",
                      borderRadius: "8px",
                      color: "var(--foreground)",
                    }}
                  />
                  <Legend />
                  {comparisonData.map((item, index) => (
                    <Line
                      key={`${item.keyword.keyword}_sentiment`}
                      type="monotone"
                      dataKey={`${item.keyword.keyword}_sentiment`}
                      stroke={getKeywordColors()[index]}
                      strokeWidth={2}
                      strokeDasharray="5 5"
                      name={`${item.keyword.keyword} (評判)`}
                    />
                  ))}
                </LineChart>
              </ResponsiveContainer>
            </div>
          )}
        </div>
      )}

      {/* 比較統計 */}
      {comparisonData.length > 0 && (
        <div className="comparison-stats mt-8 bg-white dark:bg-slate-700 p-6 rounded-lg border border-gray-200 dark:border-slate-600">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
            キーワード別統計
          </h3>
          <div className="stats-container mb-6">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={comparisonData}>
                <CartesianGrid strokeDasharray="3 3" className="opacity-30" />
                <XAxis
                  dataKey="keyword.keyword"
                  className="text-gray-600 dark:text-gray-300"
                />
                <YAxis className="text-gray-600 dark:text-gray-300" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: "var(--background)",
                    border: "1px solid var(--border)",
                    borderRadius: "8px",
                    color: "var(--foreground)",
                  }}
                />
                <Legend />
                <Bar dataKey="totalVolume" fill="#8884d8" name="総話題量" />
              </BarChart>
            </ResponsiveContainer>
          </div>

          <div className="stats-table overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200 dark:border-slate-600">
                  <th className="text-left py-3 px-4 font-medium text-gray-900 dark:text-white">
                    キーワード
                  </th>
                  <th className="text-left py-3 px-4 font-medium text-gray-900 dark:text-white">
                    総話題量
                  </th>
                  <th className="text-left py-3 px-4 font-medium text-gray-900 dark:text-white">
                    平均評判
                  </th>
                  <th className="text-left py-3 px-4 font-medium text-gray-900 dark:text-white">
                    伸び率
                  </th>
                  <th className="text-left py-3 px-4 font-medium text-gray-900 dark:text-white">
                    人気の変化
                  </th>
                </tr>
              </thead>
              <tbody>
                {comparisonData
                  .sort((a, b) => b.totalVolume - a.totalVolume)
                  .map((item) => (
                    <tr
                      key={item.keyword.id}
                      className="border-b border-gray-100 dark:border-slate-700 hover:bg-gray-50 dark:hover:bg-slate-600"
                    >
                      <td className="keyword-name py-3 px-4 font-medium text-gray-900 dark:text-white">
                        {item.keyword.keyword}
                      </td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-300">
                        {item.totalVolume.toLocaleString()}
                      </td>
                      <td className="py-3 px-4">
                        <span
                          className={`px-2 py-1 rounded-full text-xs font-medium ${
                            item.avgSentiment > 0.6
                              ? "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300"
                              : item.avgSentiment < 0.4
                              ? "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300"
                              : "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300"
                          }`}
                        >
                          {(item.avgSentiment * 100).toFixed(1)}%
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        <span
                          className={`px-2 py-1 rounded-full text-xs font-medium ${
                            item.growth > 0
                              ? "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300"
                              : "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300"
                          }`}
                        >
                          {item.growth > 0 ? "+" : ""}
                          {item.growth.toFixed(1)}%
                        </span>
                      </td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-300">
                        {item.growth > 10
                          ? "📈 急上昇"
                          : item.growth > 0
                          ? "📊 上昇"
                          : item.growth > -10
                          ? "📉 下降"
                          : "📉 急降下"}
                      </td>
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {selectedKeywords.length === 0 && (
        <div className="no-selection text-center py-12 text-gray-500 dark:text-gray-400">
          <TrendingUp className="w-12 h-12 mx-auto mb-4 opacity-50" />
          <p className="text-lg">比較したいキーワードを選択してください</p>
        </div>
      )}
    </div>
  );
};

export default MultiKeywordComparison;
