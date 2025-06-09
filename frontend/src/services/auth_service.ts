import apiService from "./api";

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  password_confirm: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
}

export interface RegisterResponse {
  id: number;
  email: string;
}

export interface RefreshRequest {
  refresh_token: string;
}

export interface RefreshResponse {
  access_token: string;
}

class AuthService {
  // ユーザー登録
  async register(data: RegisterRequest): Promise<RegisterResponse> {
    return await apiService.post<RegisterResponse>("/auth/register", data);
  }

  // ログイン
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await apiService.post<LoginResponse>("/auth/login", data);

    // トークンをローカルストレージに保存
    localStorage.setItem("access_token", response.access_token);
    localStorage.setItem("refresh_token", response.refresh_token);

    console.log("Login successful, tokens saved:", {
      access_token: response.access_token
        ? `${response.access_token.substring(0, 20)}...`
        : "None",
      refresh_token: response.refresh_token
        ? `${response.refresh_token.substring(0, 20)}...`
        : "None",
    });

    return response;
  }

  // ログアウト
  async logout(): Promise<void> {
    try {
      const refreshToken = localStorage.getItem("refresh_token");
      if (refreshToken) {
        await apiService.post("/auth/logout", {
          refresh_token: refreshToken,
        });
      }
    } finally {
      // エラーが発生してもローカルストレージはクリア
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
    }
  }

  // トークンリフレッシュ
  async refreshToken(): Promise<RefreshResponse> {
    const refreshToken = localStorage.getItem("refresh_token");
    if (!refreshToken) {
      throw new Error("リフレッシュトークンがありません");
    }

    const response = await apiService.post<RefreshResponse>("/auth/refresh", {
      refresh_token: refreshToken,
    });

    // 新しいアクセストークンを保存
    localStorage.setItem("access_token", response.access_token);

    return response;
  }

  // 認証状態チェック
  isAuthenticated(): boolean {
    const token = localStorage.getItem("access_token");
    return !!token;
  }

  // 現在のアクセストークンを取得
  getAccessToken(): string | null {
    return localStorage.getItem("access_token");
  }

  // 現在のリフレッシュトークンを取得
  getRefreshToken(): string | null {
    return localStorage.getItem("refresh_token");
  }
}

export default new AuthService();
