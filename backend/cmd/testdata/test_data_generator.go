package testdata

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	"github.com/trendscout/backend/internal/models"
)

// ダミーデータ生成用の定数
const (
	NUM_USERS    = 5
	NUM_KEYWORDS = 10
	DAYS_OF_DATA = 30
)

// テスト用ユーザーのメールドメイン
var emailDomain = "test.com"

// キーワードリスト
var fashionKeywords = []string{
	"minimalism", "vintage", "streetwear", "sustainable", "athleisure",
	"genderfluid", "cottagecore", "y2k", "normcore", "techwear",
	"cyberpunk", "darkacademia", "bohemian", "preppy", "grunge",
}

func main() {
	// カレントディレクトリを取得
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("カレントディレクトリの取得に失敗しました: %v", err)
	}
	log.Printf("カレントディレクトリ: %s", wd)

	// .env.localファイルのパスを構築
	envLocalPath := filepath.Join(wd, ".env.local")
	
	// 環境変数の読み込み
	if err := godotenv.Load(envLocalPath); err != nil {
		log.Printf("Warning: %s file not found: %v", envLocalPath, err)
		// .envファイルも試す
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	} else {
		log.Printf("環境変数を %s から読み込みました", envLocalPath)
	}

	// 接続情報の確認
	log.Printf("DB接続情報: host=%s user=%s db=%s", 
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"))

	log.Println("テストデータの生成を開始します...")

	// データベース接続を初期化
	if err := models.InitDatabases(); err != nil {
		log.Fatalf("データベース初期化エラー: %v", err)
	}
	defer models.CloseDatabases()

	// コンテキスト作成
	ctx := context.Background()

	// ユーザー生成
	userIDs, err := createUsers(ctx)
	if err != nil {
		log.Fatalf("ユーザー生成エラー: %v", err)
	}

	// キーワード生成
	keywordIDs, err := createKeywords(ctx, userIDs)
	if err != nil {
		log.Fatalf("キーワード生成エラー: %v", err)
	}

	// トレンドデータ生成
	if err := createTrendData(ctx, keywordIDs); err != nil {
		log.Fatalf("トレンドデータ生成エラー: %v", err)
	}

	log.Println("テストデータの生成が完了しました")
}

// ユーザーを生成する関数
func createUsers(ctx context.Context) ([]int, error) {
	log.Printf("%d人のテストユーザーを生成します", NUM_USERS)
	userIDs := make([]int, 0, NUM_USERS)

	// 既存のユーザーが存在するか確認
	existingUser, err := models.GetUserByEmail(ctx, "admin@example.com")
	if err == nil && existingUser != nil && existingUser.ID > 0 {
		log.Println("管理者ユーザーが既に存在します")
		userIDs = append(userIDs, existingUser.ID)
	}

	// テストユーザーの生成
	for i := 1; i <= NUM_USERS; i++ {
		email := fmt.Sprintf("user%d@%s", i, emailDomain)

		// 既に存在するかチェック
		existingUser, err := models.GetUserByEmail(ctx, email)
		if err == nil && existingUser != nil && existingUser.ID > 0 {
			log.Printf("ユーザー %s は既に存在します（ID: %d）", email, existingUser.ID)
			userIDs = append(userIDs, existingUser.ID)
			continue
		}

		// 新規ユーザー作成
		user, err := models.CreateUser(ctx, email, "password")
		if err != nil {
			log.Printf("ユーザー %s の作成に失敗しました: %v", email, err)
			continue
		}

		log.Printf("ユーザー %s を作成しました（ID: %d）", email, user.ID)
		userIDs = append(userIDs, user.ID)
	}

	return userIDs, nil
}

// キーワードを生成する関数
func createKeywords(ctx context.Context, userIDs []int) ([]int, error) {
	log.Printf("%d個のキーワードを生成します", NUM_KEYWORDS)
	keywordIDs := make([]int, 0, NUM_KEYWORDS)

	for i, keyword := range fashionKeywords {
		if i >= NUM_KEYWORDS {
			break
		}

		// ランダムなユーザーを選択
		userID := userIDs[rand.Intn(len(userIDs))]

		// 既に存在するキーワードを検索するには、ユーザーのキーワードを取得して確認
		existingKeywords, err := models.GetKeywordsForUser(ctx, userID)
		if err != nil {
			log.Printf("キーワード取得エラー: %v", err)
			continue
		}

		// キーワードが既に存在するか確認
		keywordExists := false
		for _, k := range existingKeywords {
			if k.Keyword == keyword {
				log.Printf("キーワード '%s' は既に存在します（ID: %d）", keyword, k.ID)
				keywordIDs = append(keywordIDs, k.ID)
				keywordExists = true
				break
			}
		}

		if keywordExists {
			continue
		}

		// 新規キーワード作成
		keywordObj, err := models.CreateKeyword(ctx, userID, keyword)
		if err != nil {
			log.Printf("キーワード '%s' の作成に失敗しました: %v", keyword, err)
			continue
		}

		log.Printf("キーワード '%s' を作成しました（ID: %d）", keyword, keywordObj.ID)
		keywordIDs = append(keywordIDs, keywordObj.ID)
	}

	return keywordIDs, nil
}

// トレンドデータを生成する関数
func createTrendData(ctx context.Context, keywordIDs []int) error {
	log.Printf("%d日分のトレンドデータを生成します", DAYS_OF_DATA)

	// 現在の日付
	now := time.Now()

	// 各キーワードに対してデータを生成
	for _, keywordID := range keywordIDs {
		log.Printf("キーワードID %d のトレンドデータを生成中...", keywordID)

		// 過去DAYS_OF_DATA日分のデータを生成
		for i := 0; i < DAYS_OF_DATA; i++ {
			date := now.AddDate(0, 0, -i)
			formattedDate := date.Format("2006-01-02")

			// 既にデータが存在するか確認する方法がないため、トレンドレコードを直接作成
			// ランダムな値を生成
			volume := rand.Intn(1000) + 100
			sentiment := (rand.Float64() * 2) - 1 // -1.0から1.0の範囲

			// トレンドレコード作成
			record, err := models.CreateTrendRecord(ctx, keywordID, date, volume, sentiment)
			if err != nil {
				log.Printf("トレンドレコード作成エラー: %v", err)
				continue
			}

			log.Printf("日付 %s のトレンドデータを作成しました（ボリューム: %d, センチメント: %.2f）", 
				formattedDate, record.Volume, record.Sentiment)
		}
	}

	return nil
}

func init() {
	// 乱数の初期化
	rand.Seed(time.Now().UnixNano())
} 