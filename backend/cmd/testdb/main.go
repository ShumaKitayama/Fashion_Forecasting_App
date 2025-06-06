package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/trendscout/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// 環境変数の読み込み
	// まず.env.localを読み込もうとし、失敗したら.envを読み込む
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("Warning: .env.local file not found, trying .env file")
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	} else {
		log.Println("Loaded configuration from .env.local")
	}

	// データベース接続テスト
	fmt.Println("データベース接続テストを開始します...")

	// PostgreSQL接続テスト
	fmt.Println("\n=== PostgreSQL 接続テスト ===")
	if err := testPostgres(); err != nil {
		log.Fatalf("PostgreSQL接続テスト失敗: %v", err)
	}
	fmt.Println("PostgreSQL接続テスト成功！")

	// MongoDB接続テスト
	fmt.Println("\n=== MongoDB 接続テスト ===")
	if err := testMongoDB(); err != nil {
		log.Fatalf("MongoDB接続テスト失敗: %v", err)
	}
	fmt.Println("MongoDB接続テスト成功！")

	// Redis接続テスト
	fmt.Println("\n=== Redis 接続テスト ===")
	if err := testRedis(); err != nil {
		log.Fatalf("Redis接続テスト失敗: %v", err)
	}
	fmt.Println("Redis接続テスト成功！")

	fmt.Println("\nすべてのデータベース接続テストが成功しました！")
}

// PostgreSQL接続テスト
func testPostgres() error {
	// データベース接続の初期化
	if err := models.InitDatabases(); err != nil {
		return fmt.Errorf("データベース初期化エラー: %w", err)
	}
	defer models.CloseDatabases()

	// ユーザーテーブルからデータ取得テスト
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := models.PgPool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("ユーザーテーブルクエリエラー: %w", err)
	}

	fmt.Printf("ユーザーテーブルのレコード数: %d\n", count)

	// テストユーザーの存在確認
	var email string
	err = models.PgPool.QueryRow(ctx, "SELECT email FROM users WHERE email = 'admin@example.com'").Scan(&email)
	if err != nil {
		return fmt.Errorf("テストユーザー検索エラー: %w", err)
	}

	fmt.Printf("テストユーザー検索結果: %s\n", email)
	return nil
}

// MongoDB接続テスト
func testMongoDB() error {
	// データベース接続の初期化
	if err := models.InitDatabases(); err != nil {
		return fmt.Errorf("データベース初期化エラー: %w", err)
	}
	defer models.CloseDatabases()

	// MongoDBの接続確認
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// データベース情報の取得
	fmt.Printf("MongoDB データベース名: %s\n", models.MongoDB.Name())

	// 単純なping操作
	if err := models.MongoClient.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping エラー: %w", err)
	}
	fmt.Println("MongoDB ping 成功")

	// コレクション一覧取得の代替方法
	// まずテストコレクションを作成
	testColl := models.MongoDB.Collection("test_collection")
	_, err := testColl.InsertOne(ctx, bson.M{
		"test_id":   "test123",
		"timestamp": time.Now(),
	})
	if err != nil {
		return fmt.Errorf("テストドキュメント作成エラー: %w", err)
	}
	fmt.Println("テストドキュメントの作成に成功しました")

	// コレクション一覧を確認
	collections, err := models.MongoDB.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("コレクション一覧取得エラー: %w", err)
	}

	fmt.Println("MongoDB コレクション一覧:")
	if len(collections) == 0 {
		fmt.Println("- コレクションが存在しません（初期状態）")
	} else {
		for _, name := range collections {
			fmt.Printf("- %s\n", name)
		}
	}

	// テストデータのクリーンアップ
	_, err = testColl.DeleteMany(ctx, bson.M{
		"test_id": "test123",
	})
	if err != nil {
		return fmt.Errorf("テストドキュメント削除エラー: %w", err)
	}

	fmt.Println("テストドキュメントのクリーンアップに成功しました")
	return nil
}

// Redis接続テスト
func testRedis() error {
	// Redis接続の初期化
	if err := models.InitRedis(); err != nil {
		return fmt.Errorf("Redis初期化エラー: %w", err)
	}
	defer models.CloseRedis()

	// テストキーの設定と取得
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testKey := "test:connection"
	testValue := fmt.Sprintf("test_value_%d", time.Now().Unix())

	// キーの設定
	err := models.SetWithTTL(ctx, testKey, testValue, 1*time.Minute)
	if err != nil {
		return fmt.Errorf("Redisキー設定エラー: %w", err)
	}
	fmt.Println("Redisキーの設定に成功しました")

	// キーの取得
	val, err := models.Get(ctx, testKey)
	if err != nil {
		return fmt.Errorf("Redisキー取得エラー: %w", err)
	}
	fmt.Printf("取得したキー値: %s\n", val)

	// キーの削除
	err = models.Del(ctx, testKey)
	if err != nil {
		return fmt.Errorf("Redisキー削除エラー: %w", err)
	}
	fmt.Println("Redisキーの削除に成功しました")

	return nil
} 