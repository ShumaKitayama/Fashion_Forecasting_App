import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";
import authService from "../services/auth_service";

interface AuthContextType {
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (
    email: string,
    password: string,
    passwordConfirm: string
  ) => Promise<void>;
  logout: () => Promise<void>;
  loading: boolean;
  error: string | null;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // アプリ起動時に認証状態をチェック
    const checkAuthStatus = () => {
      const authenticated = authService.isAuthenticated();
      setIsAuthenticated(authenticated);
      setLoading(false);
    };

    checkAuthStatus();
  }, []);

  const login = async (email: string, password: string): Promise<void> => {
    try {
      setLoading(true);
      setError(null);

      await authService.login({ email, password });
      setIsAuthenticated(true);
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error || "ログインに失敗しました";
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const register = async (
    email: string,
    password: string,
    passwordConfirm: string
  ): Promise<void> => {
    try {
      setLoading(true);
      setError(null);

      await authService.register({
        email,
        password,
        password_confirm: passwordConfirm,
      });

      // 登録後、自動的にログイン
      await authService.login({ email, password });
      setIsAuthenticated(true);
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || "登録に失敗しました";
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const logout = async (): Promise<void> => {
    try {
      setLoading(true);
      await authService.logout();
    } catch (err) {
      // ログアウトエラーは無視（ローカルストレージは既にクリア済み）
      console.error("Logout error:", err);
    } finally {
      setIsAuthenticated(false);
      setLoading(false);
    }
  };

  const value: AuthContextType = {
    isAuthenticated,
    login,
    register,
    logout,
    loading,
    error,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
