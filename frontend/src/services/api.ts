import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios";

const API_BASE_URL = "/api";

class ApiService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      timeout: 10000,
      headers: {
        "Content-Type": "application/json",
      },
    });

    // リクエストインターセプター - JWTトークンを自動付与
    this.api.interceptors.request.use(
      (config: any) => {
        const token = localStorage.getItem("access_token");
        console.log(
          "Request to:",
          config.url,
          "Token:",
          token ? `${token.substring(0, 20)}...` : "None"
        );
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error: any) => {
        return Promise.reject(error);
      }
    );

    // レスポンスインターセプター - 401エラー時の自動リフレッシュ
    this.api.interceptors.response.use(
      (response: AxiosResponse) => {
        return response;
      },
      async (error: any) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            const refreshToken = localStorage.getItem("refresh_token");
            if (refreshToken) {
              console.log("Token expired, refreshing...");

              // 直接axiosを使用してリフレッシュリクエストを送信
              const refreshResponse = await axios.post(
                `${API_BASE_URL}/auth/refresh`,
                {
                  refresh_token: refreshToken,
                },
                {
                  headers: {
                    "Content-Type": "application/json",
                  },
                }
              );

              const { access_token } = refreshResponse.data;

              // 新しいトークンを確実に保存
              localStorage.setItem("access_token", access_token);

              console.log(
                "Token refreshed successfully, new token:",
                access_token ? `${access_token.substring(0, 20)}...` : "None"
              );

              // 元のリクエストを新しいトークンで再実行
              if (originalRequest.headers) {
                originalRequest.headers.Authorization = `Bearer ${access_token}`;
              } else {
                originalRequest.headers = {
                  Authorization: `Bearer ${access_token}`,
                };
              }

              return this.api(originalRequest);
            } else {
              throw new Error("No refresh token available");
            }
          } catch (refreshError) {
            console.log("Token refresh failed:", refreshError);
            console.log("Redirecting to login");

            // リフレッシュも失敗した場合はログアウト
            localStorage.removeItem("access_token");
            localStorage.removeItem("refresh_token");
            window.location.href = "/login";
            return Promise.reject(refreshError);
          }
        }

        return Promise.reject(error);
      }
    );
  }

  // GET リクエスト
  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.api.get<T>(url, config);
    return response.data;
  }

  // POST リクエスト
  async post<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    const response = await this.api.post<T>(url, data, config);
    return response.data;
  }

  // PUT リクエスト
  async put<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    const response = await this.api.put<T>(url, data, config);
    return response.data;
  }

  // DELETE リクエスト
  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.api.delete<T>(url, config);
    return response.data;
  }
}

export default new ApiService();
