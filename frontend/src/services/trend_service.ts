import apiService from "./api";

export interface TrendRecord {
  id: number;
  keyword_id: number;
  date: string;
  volume: number;
  sentiment: number;
}

export interface TrendDataResponse {
  records: TrendRecord[];
  count: number;
}

export interface TrendDataRequest {
  q: number; // keyword_id
  from?: string; // YYYY-MM-DD
  to?: string; // YYYY-MM-DD
}

export interface PredictionRequest {
  keyword_id: number; // keyword_id
  days: number; // 予測日数（バックエンドのフィールド名に合わせる）
}

export interface PredictionData {
  date: string;
  volume: number;
}

export interface PredictionResponse {
  predictions: PredictionData[];
  keyword_id: number;
  horizon: number;
}

export interface SentimentRequest {
  keyword_id: number; // keyword_id
  period: number; // 期間（日数）
}

export interface SentimentApiResponse {
  keyword_id: number;
  average_sentiment: number;
  positive_count: number;
  negative_count: number;
  neutral_count: number;
  data: TrendRecord[];
  images: any[];
}

export interface SentimentResponse {
  positive: number;
  neutral: number;
  negative: number;
}

class TrendService {
  // トレンドデータ取得
  async getTrendData(params: TrendDataRequest): Promise<TrendDataResponse> {
    const queryParams = new URLSearchParams();
    queryParams.append("q", params.q.toString());

    if (params.from) {
      queryParams.append("from", params.from);
    }

    if (params.to) {
      queryParams.append("to", params.to);
    }

    return await apiService.get<TrendDataResponse>(
      `/trends/?${queryParams.toString()}`
    );
  }

  // トレンド予測
  async predictTrend(data: PredictionRequest): Promise<PredictionResponse> {
    return await apiService.post<PredictionResponse>("/trends/predict", data);
  }

  // センチメント詳細取得
  async getSentimentDetail(
    data: SentimentRequest
  ): Promise<SentimentApiResponse> {
    return await apiService.post<SentimentApiResponse>(
      "/trends/sentiment",
      data
    );
  }
}

export default new TrendService();
