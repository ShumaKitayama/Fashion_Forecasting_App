package check_user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Run() {
	// .env.localを読み込む
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatalf("Error loading .env.local file: %v", err)
	}

	log.Println("PostgreSQLのユーザーテーブルを確認します...")

	// データベース接続情報を取得
	dbHost := "localhost"
	dbUser := "trendscout"
	dbPassword := "trendscout"
	dbName := "trendscout"

	// 接続文字列を作成
	connString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		dbUser, dbPassword, dbHost, dbName)

	// データベースに接続
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("接続文字列の解析エラー: %v", err)
	}

	pgPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	defer pgPool.Close()

	// 接続テスト
	if err := pgPool.Ping(ctx); err != nil {
		log.Fatalf("データベース接続確認エラー: %v", err)
	}
	log.Println("PostgreSQLに接続しました")

	// usersテーブルを確認
	rows, err := pgPool.Query(ctx, "SELECT id, email, password_hash, created_at FROM users")
	if err != nil {
		log.Fatalf("ユーザーテーブルクエリエラー: %v", err)
	}
	defer rows.Close()

	// 結果を表示
	fmt.Println("\nユーザーテーブルの内容:")
	fmt.Println("------------------------")
	fmt.Println("ID | Email | Password Hash | Created At")
	fmt.Println("------------------------")

	count := 0
	for rows.Next() {
		var id int
		var email string
		var passwordHash string
		var createdAt time.Time

		if err := rows.Scan(&id, &email, &passwordHash, &createdAt); err != nil {
			log.Fatalf("行のスキャンエラー: %v", err)
		}

		fmt.Printf("%d | %s | %s | %s\n", id, email, passwordHash, createdAt.Format("2006-01-02 15:04:05"))
		count++
	}

	if count == 0 {
		fmt.Println("ユーザーテーブルにレコードがありません")
	} else {
		fmt.Printf("\n合計 %d 件のユーザーレコードが見つかりました\n", count)
	}

	// 特定のユーザーを検索
	var userCount int
	err = pgPool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE email = 'admin@example.com'").Scan(&userCount)
	if err != nil {
		log.Fatalf("ユーザー検索エラー: %v", err)
	}

	if userCount > 0 {
		fmt.Println("\n'admin@example.com' ユーザーが存在します")
	} else {
		fmt.Println("\n'admin@example.com' ユーザーが存在しません")
	}

	// バックエンドのルートエンドポイントへのリクエスト
	fmt.Println("\nバックエンドサーバー接続チェック...")
	fmt.Println("バックエンドサーバーに接続するには以下を実行してください:")
	fmt.Println("curl http://localhost:8080/")
}
