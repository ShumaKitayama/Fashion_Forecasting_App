package check_keywords

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

	log.Println("PostgreSQLのキーワードテーブルを確認します...")

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

	// キーワードテーブルの構造を確認
	tableInfo, err := pgPool.Query(ctx, `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'keywords'
	`)
	if err != nil {
		log.Fatalf("テーブル情報クエリエラー: %v", err)
	}
	defer tableInfo.Close()

	// テーブル構造を表示
	fmt.Println("\nキーワードテーブルの構造:")
	fmt.Println("-------------------------")
	fmt.Println("カラム名 | データ型")
	fmt.Println("-------------------------")

	for tableInfo.Next() {
		var columnName, dataType string
		if err := tableInfo.Scan(&columnName, &dataType); err != nil {
			log.Fatalf("行のスキャンエラー: %v", err)
		}
		fmt.Printf("%s | %s\n", columnName, dataType)
	}

	// キーワードテーブルの内容を確認
	rows, err := pgPool.Query(ctx, "SELECT id, user_id, keyword, created_at FROM keywords")
	if err != nil {
		log.Fatalf("キーワードテーブルクエリエラー: %v", err)
	}
	defer rows.Close()

	// 結果を表示
	fmt.Println("\nキーワードテーブルの内容:")
	fmt.Println("-------------------------")
	fmt.Println("ID | UserID | Keyword | Created At")
	fmt.Println("-------------------------")

	count := 0
	for rows.Next() {
		var id int
		var userID int
		var keyword string
		var createdAt time.Time

		if err := rows.Scan(&id, &userID, &keyword, &createdAt); err != nil {
			log.Fatalf("行のスキャンエラー: %v", err)
		}

		fmt.Printf("%d | %d | %s | %s\n", id, userID, keyword, createdAt.Format("2006-01-02 15:04:05"))
		count++
	}

	if count == 0 {
		fmt.Println("キーワードテーブルにレコードがありません")
	} else {
		fmt.Printf("\n合計 %d 件のキーワードレコードが見つかりました\n", count)
	}

	// テストデータの挿入
	fmt.Println("\nテストキーワードを挿入します...")
	// 管理者ユーザーのIDを取得
	var adminID int
	err = pgPool.QueryRow(ctx, "SELECT id FROM users WHERE email = 'admin@example.com'").Scan(&adminID)
	if err != nil {
		log.Fatalf("管理者ユーザーID取得エラー: %v", err)
	}

	// キーワードをテスト挿入
	_, err = pgPool.Exec(ctx,
		"INSERT INTO keywords (user_id, keyword) VALUES ($1, $2) ON CONFLICT (user_id, keyword) DO NOTHING",
		adminID, "test_keyword")
	if err != nil {
		log.Fatalf("キーワード挿入エラー: %v", err)
	}
	fmt.Println("テストキーワードを挿入しました")

	// 再度内容を確認
	rows, err = pgPool.Query(ctx, "SELECT id, user_id, keyword, created_at FROM keywords")
	if err != nil {
		log.Fatalf("キーワードテーブル再クエリエラー: %v", err)
	}
	defer rows.Close()

	fmt.Println("\n挿入後のキーワードテーブルの内容:")
	fmt.Println("-------------------------")
	fmt.Println("ID | UserID | Keyword | Created At")
	fmt.Println("-------------------------")

	count = 0
	for rows.Next() {
		var id int
		var userID int
		var keyword string
		var createdAt time.Time

		if err := rows.Scan(&id, &userID, &keyword, &createdAt); err != nil {
			log.Fatalf("行のスキャンエラー: %v", err)
		}

		fmt.Printf("%d | %d | %s | %s\n", id, userID, keyword, createdAt.Format("2006-01-02 15:04:05"))
		count++
	}

	fmt.Printf("\n合計 %d 件のキーワードレコードが見つかりました\n", count)
}
