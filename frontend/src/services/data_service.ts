import apiService from "./api";

export interface DataCollectionResponse {
  message: string;
  keyword: string;
  items_collected: number;
  date: string;
}

class DataService {
  // キーワードのデータ収集
  async collectKeywordData(keywordId: number): Promise<DataCollectionResponse> {
    return await apiService.post<DataCollectionResponse>(
      `/data/collect/${keywordId}`,
      {}
    );
  }
}

export default new DataService();
