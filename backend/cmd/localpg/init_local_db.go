package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// .env.localを読み込む
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatalf("Error loading .env.local file: %v", err)
	}

	log.Println("ローカルPostgreSQLデータベースの初期化を開始します...")

	// データベース接続情報
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		dbUser, dbPassword, dbHost, dbName)

	// データベース接続
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

	// 初期化スクリプトの読み込み
	schemaFile := "configs/sql/init-scripts/01-schema.sql"
	schemaSQL, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Fatalf("スキーマファイル読み込みエラー: %v", err)
	}

	// スキーマの適用
	_, err = pgPool.Exec(ctx, string(schemaSQL))
	if err != nil {
		log.Fatalf("スキーマ適用エラー: %v", err)
	}
	log.Println("スキーマを正常に適用しました")

	// ユーザーテーブルの確認
	var count int
	err = pgPool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("ユーザーテーブル確認エラー: %v", err)
	}
	log.Printf("ユーザーテーブルのレコード数: %d", count)

	log.Println("ローカルデータベースの初期化が完了しました")
} 