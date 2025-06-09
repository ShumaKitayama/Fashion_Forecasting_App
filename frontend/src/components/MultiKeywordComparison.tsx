import React, { useState, useEffect } from "react";
import keywordService, { Keyword } from "../services/keyword_service";
import trendService, { TrendRecord } from "../services/trend_service";
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
    <div className="multi-keyword-comparison">
      <div className="comparison-header">
        <h2>キーワード比較分析</h2>
        <p>
          複数のキーワードのトレンドを比較して、ファッション業界の動向を把握できます
        </p>
      </div>

      {/* コントロールパネル */}
      <div className="comparison-controls">
        <div className="date-controls">
          <label>
            開始日:
            <input
              type="date"
              value={dateFrom}
              onChange={(e) => setDateFrom(e.target.value)}
            />
          </label>
          <label>
            終了日:
            <input
              type="date"
              value={dateTo}
              onChange={(e) => setDateTo(e.target.value)}
            />
          </label>
        </div>

        <div className="view-controls">
          <label>表示モード:</label>
          <select
            value={viewMode}
            onChange={(e) =>
              setViewMode(e.target.value as "volume" | "sentiment" | "combined")
            }
          >
            <option value="volume">ボリューム</option>
            <option value="sentiment">センチメント</option>
            <option value="combined">複合ビュー</option>
          </select>
        </div>
      </div>

      {/* キーワード選択 */}
      <div className="keyword-selection">
        <h3>比較するキーワードを選択（最大5つ）</h3>
        <div className="keyword-list">
          {availableKeywords.map((keyword) => (
            <label key={keyword.id} className="keyword-checkbox">
              <input
                type="checkbox"
                checked={selectedKeywords.includes(keyword.id)}
                onChange={() => handleKeywordToggle(keyword.id)}
                disabled={
                  !selectedKeywords.includes(keyword.id) &&
                  selectedKeywords.length >= 5
                }
              />
              <span>{keyword.keyword}</span>
            </label>
          ))}
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}

      {/* 比較チャート */}
      {selectedKeywords.length > 0 && chartData.length > 0 && (
        <div className="comparison-charts">
          {(viewMode === "volume" || viewMode === "combined") && (
            <div className="chart-container">
              <h3>ボリューム比較</h3>
              <ResponsiveContainer width="100%" height={400}>
                <LineChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" tickFormatter={formatDate} />
                  <YAxis />
                  <Tooltip
                    labelFormatter={(label) => formatDate(label as string)}
                  />
                  <Legend />
                  {comparisonData.map((item, index) => (
                    <Line
                      key={`${item.keyword.keyword}_volume`}
                      type="monotone"
                      dataKey={`${item.keyword.keyword}_volume`}
                      stroke={getKeywordColors()[index]}
                      strokeWidth={2}
                      name={`${item.keyword.keyword} (ボリューム)`}
                    />
                  ))}
                </LineChart>
              </ResponsiveContainer>
            </div>
          )}

          {(viewMode === "sentiment" || viewMode === "combined") && (
            <div className="chart-container">
              <h3>センチメント比較</h3>
              <ResponsiveContainer width="100%" height={400}>
                <LineChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" tickFormatter={formatDate} />
                  <YAxis domain={[0, 1]} />
                  <Tooltip
                    labelFormatter={(label) => formatDate(label as string)}
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
                      name={`${item.keyword.keyword} (センチメント)`}
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
        <div className="comparison-stats">
          <h3>キーワード別統計</h3>
          <div className="stats-container">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={comparisonData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="keyword.keyword" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="totalVolume" fill="#8884d8" name="総ボリューム" />
              </BarChart>
            </ResponsiveContainer>
          </div>

          <div className="stats-table">
            <table>
              <thead>
                <tr>
                  <th>キーワード</th>
                  <th>総ボリューム</th>
                  <th>平均センチメント</th>
                  <th>成長率</th>
                  <th>トレンド</th>
                </tr>
              </thead>
              <tbody>
                {comparisonData
                  .sort((a, b) => b.totalVolume - a.totalVolume)
                  .map((item) => (
                    <tr key={item.keyword.id}>
                      <td className="keyword-name">{item.keyword.keyword}</td>
                      <td>{item.totalVolume.toLocaleString()}</td>
                      <td>
                        <span
                          className={`sentiment ${
                            item.avgSentiment > 0.6
                              ? "positive"
                              : item.avgSentiment < 0.4
                              ? "negative"
                              : "neutral"
                          }`}
                        >
                          {(item.avgSentiment * 100).toFixed(1)}%
                        </span>
                      </td>
                      <td>
                        <span
                          className={`growth ${
                            item.growth > 0 ? "positive" : "negative"
                          }`}
                        >
                          {item.growth > 0 ? "+" : ""}
                          {item.growth.toFixed(1)}%
                        </span>
                      </td>
                      <td>
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
        <div className="no-selection">
          比較したいキーワードを選択してください
        </div>
      )}
    </div>
  );
};

export default MultiKeywordComparison;
