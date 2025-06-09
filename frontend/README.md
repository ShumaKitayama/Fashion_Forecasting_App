# TrendScout Frontend

ファッション業界向けのトレンド予測・分析アプリケーションのフロントエンド部分です。React + TypeScript + Vite で構築されています。

## 🚀 主要機能

### 認証機能

- **ユーザー登録**: メールアドレスとパスワードでアカウント作成
- **ログイン/ログアウト**: JWT トークンベースの認証
- **自動トークンリフレッシュ**: セッション維持機能

### キーワード管理

- **キーワード登録**: トレンド分析したいキーワードの追加
- **キーワード編集**: 登録済みキーワードの修正
- **キーワード削除**: 不要なキーワードの削除
- **リアルタイム選択**: クリックで即座にトレンド表示

### トレンド分析

- **データ可視化**: Chart.js を使用した美しいグラフ表示
- **ボリューム推移**: 検索ボリュームの時系列変化
- **センチメント分析**: ポジティブ/ニュートラル/ネガティブの比率
- **日付範囲選択**: 任意の期間でのデータ分析

### 予測機能

- **未来予測**: 1-30 日先のトレンド予測
- **AI 予測**: Gemini API を活用した高精度予測
- **視覚的表示**: 点線グラフで予測データを表現

### センチメント詳細分析

- **日別分析**: 特定日のセンチメント詳細
- **ドーナツチャート**: 直感的なセンチメント比率表示
- **分析サマリー**: AI による分析結果の解釈

## 🏗️ 技術スタック

### コア技術

- **React 18**: UI ライブラリ
- **TypeScript**: 型安全な開発
- **Vite**: 高速ビルドツール

### UI・グラフ

- **Chart.js**: グラフ描画ライブラリ
- **react-chartjs-2**: Chart.js の React ラッパー
- **CSS3**: カスタムスタイリング

### HTTP 通信

- **Axios**: HTTP クライアント
- **JWT**: 認証トークン管理

### ルーティング

- **React Router**: SPA ルーティング

## 📁 プロジェクト構成

```
frontend/
├── public/                 # 静的ファイル
├── src/
│   ├── assets/            # 画像・アイコンなど
│   │   ├── KeywordManager.tsx    # キーワード管理
│   │   ├── TrendChart.tsx        # トレンドグラフ
│   │   ├── PredictionChart.tsx   # 予測グラフ
│   │   ├── SentimentAnalysis.tsx # センチメント分析
│   │   └── ProtectedRoute.tsx    # 認証保護ルート
│   ├── contexts/          # React Context
│   │   └── AuthContext.tsx       # 認証状態管理
│   ├── pages/             # ページコンポーネント
│   │   ├── LoginPage.tsx         # ログインページ
│   │   ├── RegisterPage.tsx      # 登録ページ
│   │   └── DashboardPage.tsx     # ダッシュボード
│   ├── services/          # API通信層
│   │   ├── api.ts               # ベースAPIクライアント
│   │   ├── auth_service.ts      # 認証API
│   │   ├── keyword_service.ts   # キーワードAPI
│   │   └── trend_service.ts     # トレンドAPI
│   ├── utils/             # ユーティリティ関数
│   │   └── dateUtils.ts         # 日付操作
│   ├── App.tsx            # メインアプリコンポーネント
│   ├── App.css            # グローバルスタイル
│   └── main.tsx           # エントリーポイント
├── package.json           # 依存関係
├── tsconfig.json          # TypeScript設定
├── vite.config.ts         # Vite設定
└── README.md             # このファイル
```

## 🎨 UI・UX の特徴

### デザインシステム

- **モダンなデザイン**: 美しいグラデーションと影
- **レスポンシブ**: モバイル・タブレット対応
- **直感的 UI**: わかりやすいアイコンとラベル

### ユーザビリティ

- **ローディング状態**: 処理中の視覚的フィードバック
- **エラーハンドリング**: わかりやすいエラーメッセージ
- **バリデーション**: リアルタイム入力検証

## 🔧 開発コマンド

### 開発サーバー起動

```bash
npm run dev
```

### 本番ビルド

```bash
npm run build
```

### ビルド確認

```bash
npm run preview
```

### 型チェック・Lint

```bash
npm run lint
```

## 🌐 API 通信

### ベース URL

```
http://localhost:8080/api
```

### 認証

- JWT トークンを`Authorization: Bearer {token}`ヘッダーで送信
- 自動トークンリフレッシュ機能
- 401 エラー時の自動ログアウト

### エンドポイント例

```typescript
// 認証
POST / auth / register;
POST / auth / login;
POST / auth / logout;
POST / auth / refresh;

// キーワード管理
GET / keywords / POST / keywords / PUT / keywords / { id };
DELETE / keywords / { id };

// トレンド分析
GET / trends / POST / trends / predict;
POST / trends / sentiment;
```

## 🔐 セキュリティ機能

### 認証・認可

- **JWT 認証**: アクセストークン（15 分）+ リフレッシュトークン（7 日）
- **自動ログアウト**: トークン期限切れ時
- **ルート保護**: 未認証時の自動リダイレクト

### データ保護

- **入力検証**: フロントエンド・バックエンド両方で実施
- **XSS 対策**: React 標準の自動エスケープ
- **CORS 設定**: 適切なオリジン制限

## 📱 レスポンシブ対応

### ブレークポイント

- **デスクトップ**: 1200px 以上
- **タブレット**: 768px-1199px
- **モバイル**: 767px 以下

### 対応機能

- **フレキシブルレイアウト**: 画面サイズに応じた配置変更
- **タッチ操作**: モバイルフレンドリーなボタンサイズ
- **コンテンツ優先**: 重要な機能から表示

## 🚧 今後の拡張予定

### 新機能

- **ダークモード**: テーマ切り替え機能
- **通知機能**: 予測アラート・トレンド変化通知
- **エクスポート**: データの CSV・PDF 出力
- **キーワード分析**: 関連キーワード提案

### パフォーマンス改善

- **メモ化**: React.memo、useMemo 活用
- **レイジーローディング**: ページ分割読み込み
- **キャッシュ最適化**: データ取得の効率化

## 🤝 開発ガイドライン

### コーディング規約

- **TypeScript**: 厳密型チェック有効
- **コンポーネント**: 関数コンポーネント + Hooks
- **状態管理**: useState、useContext 中心
- **スタイル**: CSS Modules または Styled Components

### ディレクトリ命名

- **PascalCase**: コンポーネントファイル
- **camelCase**: ユーティリティ関数
- **kebab-case**: CSS クラス名

## 📞 サポート

### トラブルシューティング

1. **依存関係エラー**: `npm install` 再実行
2. **型エラー**: TypeScript 設定確認
3. **API 接続エラー**: バックエンド起動状況確認

### 開発者向け情報

- **API ドキュメント**: `AI_Memory/api_specification.md`
- **システム設計**: `AI_Memory/system_architecture.md`
- **データベース**: `AI_Memory/database_structure.md`
