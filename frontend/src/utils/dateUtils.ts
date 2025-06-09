// 日付フォーマット関数
export const formatDate = (date: Date | string): string => {
  const d = new Date(date);
  return d.toLocaleDateString("ja-JP", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};

// 日付を YYYY-MM-DD 形式にフォーマット
export const formatDateForInput = (date: Date | string): string => {
  const d = new Date(date);
  return d.toISOString().split("T")[0];
};

// 指定日数前の日付を取得
export const getDateBefore = (days: number): string => {
  const date = new Date();
  date.setDate(date.getDate() - days);
  return formatDateForInput(date);
};

// 今日の日付を取得
export const getToday = (): string => {
  return formatDateForInput(new Date());
};

// 日付範囲が有効かチェック
export const isValidDateRange = (from: string, to: string): boolean => {
  const fromDate = new Date(from);
  const toDate = new Date(to);
  return fromDate <= toDate;
};
