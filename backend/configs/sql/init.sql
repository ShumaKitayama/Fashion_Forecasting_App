-- ユーザーテーブル
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- 認証トークンテーブル
CREATE TABLE IF NOT EXISTS auth_tokens (
  token_id UUID PRIMARY KEY,
  user_id INT REFERENCES users(id),
  expires_at TIMESTAMP NOT NULL
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_auth_user ON auth_tokens(user_id);

-- キーワードテーブル
CREATE TABLE IF NOT EXISTS keywords (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id),
  keyword VARCHAR(100) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_id, keyword)
);

-- トレンドレコードテーブル
CREATE TABLE IF NOT EXISTS trend_records (
  id BIGSERIAL PRIMARY KEY,
  keyword_id INT REFERENCES keywords(id),
  record_date DATE NOT NULL,
  volume INT NOT NULL,
  sentiment FLOAT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(keyword_id, record_date)
);

-- 管理者ユーザー追加（開発用）
INSERT INTO users (email, password_hash) 
VALUES ('admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy') -- パスワード: password
ON CONFLICT (email) DO NOTHING; 