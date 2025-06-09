import apiService from "./api";

export interface Keyword {
  id: number;
  user_id: number;
  keyword: string;
  created_at: string;
}

export interface KeywordListResponse {
  keywords: Keyword[];
  count: number;
}

export interface CreateKeywordRequest {
  keyword: string;
}

export interface CreateKeywordResponse {
  id: number;
  keyword: string;
}

export interface UpdateKeywordRequest {
  keyword: string;
}

class KeywordService {
  // キーワード一覧取得
  async getKeywords(): Promise<KeywordListResponse> {
    return await apiService.get<KeywordListResponse>("/keywords/");
  }

  // キーワード作成
  async createKeyword(
    data: CreateKeywordRequest
  ): Promise<CreateKeywordResponse> {
    return await apiService.post<CreateKeywordResponse>("/keywords/", data);
  }

  // キーワード更新
  async updateKeyword(id: number, data: UpdateKeywordRequest): Promise<void> {
    return await apiService.put<void>(`/keywords/${id}`, data);
  }

  // キーワード削除
  async deleteKeyword(id: number): Promise<void> {
    return await apiService.delete<void>(`/keywords/${id}`);
  }
}

export default new KeywordService();
