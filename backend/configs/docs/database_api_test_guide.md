# TrendScout データベース・API テスト手順書

このドキュメントでは、TrendScout アプリケーションのデータベース接続テストと API 動作確認の手順について説明します。

## 1. Docker 環境の起動

まず、Docker 環境を起動して各種サービスを利用可能な状態にします。

```bash
# プロジェクトのルートディレクトリで実行
docker-compose up -d
```

これにより、以下のサービスが起動します：

- PostgreSQL: ユーザー情報、キーワード、トレンドデータの保存
- MongoDB: 画像・投稿データの保存
- Redis: トークン管理
- バックエンドサーバー: Go アプリケーション

## 2. ローカル環境の準備（ローカルテスト時のみ）

ローカルマシンから直接テストを実行する場合は、データベースの初期化が必要な場合があります。

### 2.1 ローカル PostgreSQL の初期化

初めてテストを実行する場合や、スキーマを更新した場合は、以下のコマンドでローカルデータベースを初期化します：

```bash
# プロジェクトのルートディレクトリで実行
cd backend
go run cmd/localpg/init_local_db.go
```

このコマンドは以下の処理を行います：

- PostgreSQL への接続
- スキーマの適用（テーブル作成）
- テストユーザーの作成

## 3. データベース接続テスト

データベース接続テストを行うには、作成したテストツールを実行します。

### 3.1 テスト実行方法について

**注意**: 本番用 Docker コンテナには Go 開発環境が含まれていないため、コンテナ内で直接 `go run` コマンドを実行することはできません。これは軽量なコンテナを維持するためのベストプラクティスです。

テストはローカル環境で実行してください。Docker 環境を起動した状態で、テストツールはローカルマシンから各データベースに接続します。

```bash
# プロジェクトのルートディレクトリで実行
cd backend
go run cmd/testdb/main.go
```

> **補足**: コンテナ内でテストを実行する場合は、開発環境を含むカスタム Docker イメージを作成するか、既存のビルダーステージを使用する必要があります。必要に応じて DevOps チームにご相談ください。

このコマンドは以下の接続テストを行います：

1. PostgreSQL 接続テスト：

   - データベース接続
   - ユーザーテーブルの確認
   - テストユーザーの確認

2. MongoDB 接続テスト：

   - データベース接続
   - コレクション一覧の取得
   - テストドキュメントの作成と削除

3. Redis 接続テスト：
   - Redis サーバーへの接続
   - キーの設定と取得
   - キーの削除

すべてのテストが成功すると、「すべてのデータベース接続テストが成功しました！」というメッセージが表示されます。

## 4. API 動作確認テスト

API 動作確認テストには、2 つの方法があります：

### 4.1 テストツールを使用する方法

作成した API テストツールを実行します。

```bash
# プロジェクトのルートディレクトリで実行
cd backend
go run cmd/testapi/main.go
```

このコマンドは以下の API テストを行います：

1. 認証 API テスト：

   - ログイン
   - 認証トークンの取得

2. キーワード API テスト：

   - キーワード一覧取得
   - キーワード作成

3. トレンド API テスト：
   - トレンド予測
   - センチメント分析

### 4.2 Postman を使用する方法

1. Postman をインストールします（[https://www.postman.com/downloads/](https://www.postman.com/downloads/)）

2. 作成したコレクションをインポートします：

   - Postman を起動
   - 「Import」ボタンをクリック
   - `backend/configs/postman/trendscout_api_collection.json` ファイルを選択

3. 環境変数を設定します：

   - 「Environments」タブで新しい環境を作成
   - `base_url` 変数に `http://localhost:8080/api` を設定
   - `auth_token` 変数は空のままにしておく（ログイン時に自動で設定されます）

4. API テストを実行します：
   - 認証 API > ログイン：管理者ユーザーでログイン
   - キーワード API > キーワード一覧取得：登録済みキーワードの確認
   - キーワード API > キーワード作成：新しいキーワードの登録
   - トレンド API > トレンド予測：キーワードに基づく予測の実行

## 5. トラブルシューティング

### データベース接続エラー

1. Docker コンテナが起動しているか確認します：

```bash
docker ps
```

2. 環境変数が正しく設定されているか確認します：

```bash
cat backend/.env
cat backend/.env.local
```

3. 必要に応じてコンテナを再起動します：

```bash
docker-compose restart postgres mongo redis
```

4. ローカルテスト時のホスト名解決エラーが発生した場合：

```bash
# エラー例: hostname resolving error (lookup postgres: no such host)
```

.env.local ファイルのホスト設定が正しいか確認してください：

```
DB_HOST=localhost   # Docker 外からアクセスする場合は localhost を使用
MONGO_URI=mongodb://localhost:27017
REDIS_HOST=localhost
```

### API 接続エラー

1. バックエンドサーバーが起動しているか確認します：

```bash
docker logs trendscout_backend
```

2. サーバーが応答するか確認します：

```bash
curl http://localhost:8080/
```

3. 必要に応じてバックエンドを再起動します：

```bash
docker-compose restart backend
```

## 6. テスト完了後

すべてのテストが完了したら、Docker 環境を停止します：

```bash
docker-compose down
```

開発を続ける場合は、必要なコンテナのみを起動したままにしておくことができます。
