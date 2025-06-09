import React, { useState, useEffect } from "react";
import { useAuth } from "../contexts/AuthContext";
import keywordService, { Keyword } from "../services/keyword_service";
import trendService, {
  TrendRecord,
  PredictionData,
} from "../services/trend_service";
import TrendChart from "../components/TrendChart";
import KeywordManager from "../components/KeywordManager";
import PredictionChart from "../components/PredictionChart";
import SentimentAnalysis from "../components/SentimentAnalysis";
import MultiKeywordComparison from "../components/MultiKeywordComparison";
import { Button } from "../components/ui/Button";
import { Input } from "../components/ui/Input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../components/ui/Card";
import { ThemeToggle } from "../components/ThemeToggle";
import { Spinner, LoadingScreen } from "../components/ui/Spinner";
import {
  TrendingUp,
  BarChart3,
  Calendar,
  RefreshCw,
  Settings,
  ChevronDown,
  Activity,
  PieChart,
  LogOut,
  AlertCircle,
  Plus,
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

const DashboardPage: React.FC = () => {
  const { logout } = useAuth();
  const [keywords, setKeywords] = useState<Keyword[]>([]);
  const [selectedKeyword, setSelectedKeyword] = useState<Keyword | null>(null);
  const [trendData, setTrendData] = useState<TrendRecord[]>([]);
  const [predictions, setPredictions] = useState<PredictionData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 日付範囲のstate
  const [dateFrom, setDateFrom] = useState<string>(() => {
    const date = new Date();
    date.setDate(date.getDate() - 90); // 90日前（3ヶ月）
    return date.toISOString().split("T")[0];
  });

  const [dateTo, setDateTo] = useState<string>(() => {
    const date = new Date();
    return date.toISOString().split("T")[0];
  });

  const [predictionHorizon, setPredictionHorizon] = useState(7);
  const [activeTab, setActiveTab] = useState<"analysis" | "comparison">(
    "analysis"
  );
  const [sidebarOpen, setSidebarOpen] = useState(true);

  // 初期データ読み込み
  useEffect(() => {
    loadKeywords();
  }, []);

  // 選択されたキーワードが変更されたときにトレンドデータを読み込む
  useEffect(() => {
    if (selectedKeyword) {
      loadTrendData();
    }
  }, [selectedKeyword, dateFrom, dateTo]);

  const loadKeywords = async () => {
    try {
      const response = await keywordService.getKeywords();
      setKeywords(response.keywords);

      // 最初のキーワードを自動選択
      if (response.keywords.length > 0) {
        setSelectedKeyword(response.keywords[0]);
      }
    } catch (err: any) {
      setError("キーワードの読み込みに失敗しました");
      console.error("Load keywords error:", err);
    }
  };

  const loadTrendData = async () => {
    if (!selectedKeyword) return;

    setLoading(true);
    setError(null);

    try {
      const response = await trendService.getTrendData({
        q: selectedKeyword.id,
        from: dateFrom,
        to: dateTo,
      });

      setTrendData(response.records);
    } catch (err: any) {
      setError("トレンドデータの読み込みに失敗しました");
      console.error("Load trend data error:", err);
    } finally {
      setLoading(false);
    }
  };

  const loadPredictions = async () => {
    if (!selectedKeyword) return;

    setLoading(true);
    setError(null);

    try {
      const response = await trendService.predictTrend({
        keyword_id: selectedKeyword.id,
        days: predictionHorizon,
      });

      setPredictions(response.predictions);
    } catch (err: any) {
      setError("予測データの読み込みに失敗しました");
      console.error("Load predictions error:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleKeywordSelect = (keyword: Keyword) => {
    setSelectedKeyword(keyword);
    setPredictions([]); // 予測データをクリア
  };

  const handleKeywordUpdate = () => {
    loadKeywords(); // キーワード一覧を再読み込み
  };

  const handleLogout = async () => {
    try {
      await logout();
    } catch (err) {
      console.error("Logout error:", err);
    }
  };

  const stats = [
    {
      name: "総キーワード数",
      value: keywords.length,
      icon: TrendingUp,
      change: "+12%",
      changeType: "positive" as const,
    },
    {
      name: "アクティブ分析",
      value: selectedKeyword ? 1 : 0,
      icon: Activity,
      change: "",
      changeType: "neutral" as const,
    },
    {
      name: "データポイント",
      value: trendData.length,
      icon: BarChart3,
      change: "+5%",
      changeType: "positive" as const,
    },
    {
      name: "予測精度",
      value: "94%",
      icon: PieChart,
      change: "+2%",
      changeType: "positive" as const,
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50 dark:from-slate-900 dark:via-slate-800 dark:to-slate-900">
      {/* Modern Header */}
      <motion.header
        initial={{ y: -100 }}
        animate={{ y: 0 }}
        className="glass backdrop-blur-xl border-b border-white/20 dark:border-white/10 sticky top-0 z-50"
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center space-x-4">
              <motion.div
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.95 }}
                className="p-2 rounded-lg bg-primary/10 dark:bg-primary/20"
              >
                <TrendingUp className="w-8 h-8 text-primary" />
              </motion.div>
              <div>
                <h1 className="text-2xl font-bold bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
                  TrendScout
                </h1>
                <p className="text-sm text-muted-foreground">
                  Fashion Trend Analytics
                </p>
              </div>
            </div>

            <div className="flex items-center space-x-4">
              <ThemeToggle />
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setSidebarOpen(!sidebarOpen)}
                className="lg:hidden"
              >
                <Settings className="w-4 h-4" />
              </Button>
              <Button
                variant="outline"
                onClick={handleLogout}
                className="bg-white/50 dark:bg-black/50 border-white/20 dark:border-white/10 hover:bg-red-50 dark:hover:bg-red-900/20"
              >
                <LogOut className="w-4 h-4 mr-2" />
                ログアウト
              </Button>
            </div>
          </div>
        </div>
      </motion.header>

      {/* Error Display */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -50 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -50 }}
            className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-4"
          >
            <div className="bg-destructive/10 border border-destructive/20 text-destructive px-4 py-3 rounded-lg flex items-center space-x-2">
              <AlertCircle className="w-5 h-5" />
              <span>{error}</span>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Stats Overview */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8"
        >
          {stats.map((stat, index) => (
            <motion.div
              key={stat.name}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 + index * 0.1 }}
              whileHover={{ scale: 1.02 }}
            >
              <Card className="glass backdrop-blur-md border-white/20 dark:border-white/10 hover:shadow-lg transition-all duration-300">
                <CardContent className="p-6">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium text-muted-foreground">
                        {stat.name}
                      </p>
                      <div className="flex items-center space-x-2">
                        <p className="text-2xl font-bold">{stat.value}</p>
                        {stat.change && (
                          <span
                            className={`text-xs font-medium ${
                              stat.changeType === "positive"
                                ? "text-green-600 dark:text-green-400"
                                : "text-red-600 dark:text-red-400"
                            }`}
                          >
                            {stat.change}
                          </span>
                        )}
                      </div>
                    </div>
                    <div className="p-3 rounded-full bg-primary/10 dark:bg-primary/20">
                      <stat.icon className="w-6 h-6 text-primary" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </motion.div>
          ))}
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          {/* Sidebar - Keyword Management */}
          <motion.aside
            initial={{ x: -300 }}
            animate={{ x: sidebarOpen ? 0 : -300 }}
            className={`lg:col-span-1 ${
              sidebarOpen ? "block" : "hidden lg:block"
            }`}
          >
            <Card className="glass backdrop-blur-md border-white/20 dark:border-white/10 sticky top-24">
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Plus className="w-5 h-5" />
                  <span>キーワード管理</span>
                </CardTitle>
                <CardDescription>
                  トレンドを分析するキーワードを管理
                </CardDescription>
              </CardHeader>
              <CardContent>
                <KeywordManager
                  keywords={keywords}
                  selectedKeyword={selectedKeyword}
                  onKeywordSelect={handleKeywordSelect}
                  onKeywordUpdate={handleKeywordUpdate}
                />
              </CardContent>
            </Card>
          </motion.aside>

          {/* Main Content */}
          <motion.main
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.3 }}
            className="lg:col-span-3 space-y-8"
          >
            {/* Tab Navigation */}
            <Card className="glass backdrop-blur-md border-white/20 dark:border-white/10">
              <CardContent className="p-6">
                <div className="flex space-x-1 bg-muted/50 rounded-lg p-1">
                  {[
                    { id: "analysis", label: "個別分析", icon: BarChart3 },
                    { id: "comparison", label: "比較分析", icon: Activity },
                  ].map((tab) => (
                    <Button
                      key={tab.id}
                      variant={activeTab === tab.id ? "default" : "ghost"}
                      size="sm"
                      onClick={() =>
                        setActiveTab(tab.id as "analysis" | "comparison")
                      }
                      className={`flex-1 ${
                        activeTab === tab.id
                          ? "bg-white dark:bg-slate-800 shadow-sm"
                          : "hover:bg-white/50 dark:hover:bg-slate-700/50"
                      }`}
                    >
                      <tab.icon className="w-4 h-4 mr-2" />
                      {tab.label}
                    </Button>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Content Area */}
            <AnimatePresence mode="wait">
              {activeTab === "comparison" ? (
                <motion.div
                  key="comparison"
                  initial={{ opacity: 0, x: 20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -20 }}
                >
                  <Card className="glass backdrop-blur-md border-white/20 dark:border-white/10">
                    <CardHeader>
                      <CardTitle className="flex items-center space-x-2">
                        <Activity className="w-5 h-5" />
                        <span>キーワード比較分析</span>
                      </CardTitle>
                    </CardHeader>
                    <CardContent>
                      <MultiKeywordComparison />
                    </CardContent>
                  </Card>
                </motion.div>
              ) : selectedKeyword ? (
                <motion.div
                  key="analysis"
                  initial={{ opacity: 0, x: 20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -20 }}
                  className="space-y-8"
                >
                  {/* Control Panel */}
                  <Card className="glass backdrop-blur-md border-white/20 dark:border-white/10">
                    <CardHeader>
                      <CardTitle className="flex items-center space-x-2">
                        <Settings className="w-5 h-5" />
                        <span>{selectedKeyword.keyword} の分析設定</span>
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-6">
                      {/* Date Controls */}
                      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                        <div className="space-y-2">
                          <label className="text-sm font-medium flex items-center space-x-2">
                            <Calendar className="w-4 h-4" />
                            <span>開始日</span>
                          </label>
                          <Input
                            type="date"
                            value={dateFrom}
                            onChange={(e) => setDateFrom(e.target.value)}
                          />
                        </div>
                        <div className="space-y-2">
                          <label className="text-sm font-medium flex items-center space-x-2">
                            <Calendar className="w-4 h-4" />
                            <span>終了日</span>
                          </label>
                          <Input
                            type="date"
                            value={dateTo}
                            onChange={(e) => setDateTo(e.target.value)}
                          />
                        </div>
                        <div className="space-y-2">
                          <label className="text-sm font-medium">
                            予測期間: {predictionHorizon}日
                          </label>
                          <Input
                            type="range"
                            min="1"
                            max="60"
                            value={predictionHorizon}
                            onChange={(e) =>
                              setPredictionHorizon(Number(e.target.value))
                            }
                            className="w-full"
                          />
                        </div>
                        <div className="flex flex-col space-y-2">
                          <Button
                            onClick={loadTrendData}
                            disabled={loading}
                            className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700"
                          >
                            {loading ? (
                              <Spinner size="sm" className="mr-2" />
                            ) : (
                              <RefreshCw className="w-4 h-4 mr-2" />
                            )}
                            データ更新
                          </Button>
                          <Button
                            onClick={loadPredictions}
                            disabled={loading}
                            variant="outline"
                            className="bg-white/50 dark:bg-black/50"
                          >
                            {loading ? (
                              <Spinner size="sm" className="mr-2" />
                            ) : (
                              <TrendingUp className="w-4 h-4 mr-2" />
                            )}
                            予測生成
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>

                  {/* Charts Container */}
                  <div className="grid grid-cols-1 gap-8">
                    {/* Trend Chart */}
                    <Card className="glass backdrop-blur-md border-white/20 dark:border-white/10">
                      <CardHeader>
                        <CardTitle className="flex items-center space-x-2">
                          <BarChart3 className="w-5 h-5" />
                          <span>トレンドデータ</span>
                        </CardTitle>
                      </CardHeader>
                      <CardContent>
                        {loading ? (
                          <div className="flex items-center justify-center h-64">
                            <div className="text-center">
                              <Spinner size="lg" className="mb-4" />
                              <p className="text-muted-foreground">
                                データを読み込み中...
                              </p>
                            </div>
                          </div>
                        ) : (
                          <TrendChart data={trendData} loading={loading} />
                        )}
                      </CardContent>
                    </Card>

                    {/* Prediction Chart */}
                    {predictions.length > 0 && (
                      <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                      >
                        <Card className="glass backdrop-blur-md border-white/20 dark:border-white/10">
                          <CardHeader>
                            <CardTitle className="flex items-center space-x-2">
                              <TrendingUp className="w-5 h-5" />
                              <span>予測データ</span>
                            </CardTitle>
                          </CardHeader>
                          <CardContent>
                            <PredictionChart
                              data={predictions}
                              loading={loading}
                            />
                          </CardContent>
                        </Card>
                      </motion.div>
                    )}

                    {/* Sentiment Analysis */}
                    <motion.div
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: 0.2 }}
                    >
                      <Card className="glass backdrop-blur-md border-white/20 dark:border-white/10">
                        <CardHeader>
                          <CardTitle className="flex items-center space-x-2">
                            <PieChart className="w-5 h-5" />
                            <span>センチメント分析</span>
                          </CardTitle>
                        </CardHeader>
                        <CardContent>
                          <SentimentAnalysis keyword={selectedKeyword} />
                        </CardContent>
                      </Card>
                    </motion.div>
                  </div>
                </motion.div>
              ) : (
                <motion.div
                  key="no-selection"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                >
                  <Card className="glass backdrop-blur-md border-white/20 dark:border-white/10">
                    <CardContent className="p-12 text-center">
                      <div className="mx-auto w-24 h-24 bg-primary/10 dark:bg-primary/20 rounded-full flex items-center justify-center mb-6">
                        <TrendingUp className="w-12 h-12 text-primary" />
                      </div>
                      <h2 className="text-2xl font-bold mb-4">
                        キーワードを選択してください
                      </h2>
                      <p className="text-muted-foreground max-w-md mx-auto">
                        左側のパネルからキーワードを選択するか、新しいキーワードを追加してトレンド分析を開始しましょう。
                      </p>
                    </CardContent>
                  </Card>
                </motion.div>
              )}
            </AnimatePresence>
          </motion.main>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
