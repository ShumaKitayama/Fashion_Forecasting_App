import React, { useState } from "react";
import keywordService, { Keyword } from "../services/keyword_service";
import dataService from "../services/data_service";
import { Button } from "./ui/Button";
import { Input } from "./ui/Input";
import { Spinner } from "./ui/Spinner";
import {
  Plus,
  Edit3,
  Trash2,
  Database,
  Check,
  X,
  AlertCircle,
  Search,
  TrendingUp,
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

interface KeywordManagerProps {
  keywords: Keyword[];
  selectedKeyword: Keyword | null;
  onKeywordSelect: (keyword: Keyword) => void;
  onKeywordUpdate: () => void;
}

const KeywordManager: React.FC<KeywordManagerProps> = ({
  keywords,
  selectedKeyword,
  onKeywordSelect,
  onKeywordUpdate,
}) => {
  const [newKeyword, setNewKeyword] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editingValue, setEditingValue] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [collectingId, setCollectingId] = useState<number | null>(null);

  // 新しいキーワードを追加
  const handleAddKeyword = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newKeyword.trim()) return;

    setLoading(true);
    setError(null);

    try {
      await keywordService.createKeyword({ keyword: newKeyword.trim() });
      setNewKeyword("");
      onKeywordUpdate();
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error || "キーワードの追加に失敗しました";
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // キーワードの編集を開始
  const handleStartEdit = (keyword: Keyword) => {
    setEditingId(keyword.id);
    setEditingValue(keyword.keyword);
  };

  // キーワードの編集をキャンセル
  const handleCancelEdit = () => {
    setEditingId(null);
    setEditingValue("");
  };

  // キーワードの編集を保存
  const handleSaveEdit = async (id: number) => {
    if (!editingValue.trim()) return;

    setLoading(true);
    setError(null);

    try {
      await keywordService.updateKeyword(id, { keyword: editingValue.trim() });
      setEditingId(null);
      setEditingValue("");
      onKeywordUpdate();
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error || "キーワードの更新に失敗しました";
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // キーワードを削除
  const handleDeleteKeyword = async (id: number) => {
    if (!confirm("このキーワードを削除してもよろしいですか？")) return;

    setLoading(true);
    setError(null);

    try {
      await keywordService.deleteKeyword(id);
      onKeywordUpdate();
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error || "キーワードの削除に失敗しました";
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // データ収集
  const handleCollectData = async (keyword: Keyword) => {
    setCollectingId(keyword.id);
    setError(null);

    try {
      const result = await dataService.collectKeywordData(keyword.id);
      alert(
        `データ収集完了！\n${result.keyword}: ${result.items_collected}件のデータを収集しました。`
      );
      onKeywordUpdate(); // トレンドデータが更新されるので再読み込み
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error || "データ収集に失敗しました";
      setError(errorMessage);
    } finally {
      setCollectingId(null);
    }
  };

  return (
    <div className="space-y-6">
      {/* エラー表示 */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm flex items-center space-x-2"
          >
            <AlertCircle className="w-4 h-4" />
            <span>{error}</span>
          </motion.div>
        )}
      </AnimatePresence>

      {/* 新しいキーワード追加フォーム */}
      <form onSubmit={handleAddKeyword} className="space-y-3">
        <div className="space-y-2">
          <label className="text-sm font-medium">新しいキーワード</label>
          <div className="flex space-x-2">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                type="text"
                value={newKeyword}
                onChange={(e) => setNewKeyword(e.target.value)}
                placeholder="キーワードを入力..."
                className="pl-10"
                disabled={loading}
              />
            </div>
            <Button
              type="submit"
              disabled={loading || !newKeyword.trim()}
              size="sm"
              className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700"
            >
              {loading ? (
                <Spinner size="sm" />
              ) : (
                <>
                  <Plus className="w-4 h-4 mr-1" />
                  追加
                </>
              )}
            </Button>
          </div>
        </div>
      </form>

      {/* キーワード一覧 */}
      <div className="space-y-3">
        <h4 className="text-sm font-medium text-muted-foreground">
          キーワード一覧 ({keywords.length})
        </h4>

        {keywords.length === 0 ? (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="p-8 text-center"
          >
            <div className="w-16 h-16 mx-auto mb-4 bg-muted/50 rounded-full flex items-center justify-center">
              <TrendingUp className="w-8 h-8 text-muted-foreground" />
            </div>
            <p className="text-muted-foreground text-sm">
              まだキーワードがありません
            </p>
            <p className="text-muted-foreground text-xs mt-1">
              上のフォームから追加してください
            </p>
          </motion.div>
        ) : (
          <div className="space-y-2">
            <AnimatePresence>
              {keywords.map((keyword, index) => (
                <motion.div
                  key={keyword.id}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  transition={{ delay: index * 0.05 }}
                  className={`p-3 rounded-lg border transition-all duration-200 ${
                    selectedKeyword?.id === keyword.id
                      ? "bg-primary/10 border-primary/30 shadow-sm"
                      : "bg-background/50 border-border hover:bg-accent/50 hover:border-accent"
                  }`}
                >
                  {editingId === keyword.id ? (
                    // 編集モード
                    <motion.div
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      className="space-y-3"
                    >
                      <Input
                        type="text"
                        value={editingValue}
                        onChange={(e) => setEditingValue(e.target.value)}
                        disabled={loading}
                        onKeyDown={(e) => {
                          if (e.key === "Enter") {
                            handleSaveEdit(keyword.id);
                          } else if (e.key === "Escape") {
                            handleCancelEdit();
                          }
                        }}
                        autoFocus
                        className="text-sm"
                      />
                      <div className="flex space-x-2">
                        <Button
                          onClick={() => handleSaveEdit(keyword.id)}
                          disabled={loading || !editingValue.trim()}
                          size="sm"
                          className="flex-1"
                        >
                          {loading ? (
                            <Spinner size="sm" />
                          ) : (
                            <>
                              <Check className="w-3 h-3 mr-1" />
                              保存
                            </>
                          )}
                        </Button>
                        <Button
                          onClick={handleCancelEdit}
                          disabled={loading}
                          variant="outline"
                          size="sm"
                          className="flex-1"
                        >
                          <X className="w-3 h-3 mr-1" />
                          キャンセル
                        </Button>
                      </div>
                    </motion.div>
                  ) : (
                    // 表示モード
                    <div className="space-y-3">
                      <button
                        onClick={() => onKeywordSelect(keyword)}
                        className="w-full text-left"
                      >
                        <div className="flex items-center justify-between">
                          <span className="font-medium text-sm">
                            {keyword.keyword}
                          </span>
                          {selectedKeyword?.id === keyword.id && (
                            <div className="w-2 h-2 bg-primary rounded-full" />
                          )}
                        </div>
                      </button>

                      <div className="flex space-x-1">
                        <Button
                          onClick={() => handleCollectData(keyword)}
                          disabled={loading || collectingId === keyword.id}
                          variant="outline"
                          size="sm"
                          className="flex-1 text-xs"
                        >
                          {collectingId === keyword.id ? (
                            <Spinner size="sm" />
                          ) : (
                            <>
                              <Database className="w-3 h-3 mr-1" />
                              収集
                            </>
                          )}
                        </Button>
                        <Button
                          onClick={() => handleStartEdit(keyword)}
                          disabled={loading}
                          variant="outline"
                          size="sm"
                          className="flex-1 text-xs"
                        >
                          <Edit3 className="w-3 h-3 mr-1" />
                          編集
                        </Button>
                        <Button
                          onClick={() => handleDeleteKeyword(keyword.id)}
                          disabled={loading}
                          variant="outline"
                          size="sm"
                          className="flex-1 text-xs hover:bg-destructive/10 hover:text-destructive hover:border-destructive/30"
                        >
                          <Trash2 className="w-3 h-3 mr-1" />
                          削除
                        </Button>
                      </div>
                    </div>
                  )}
                </motion.div>
              ))}
            </AnimatePresence>
          </div>
        )}
      </div>
    </div>
  );
};

export default KeywordManager;
