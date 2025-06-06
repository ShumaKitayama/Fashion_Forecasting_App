package scraper

import (
	"context"
	"fmt"
	"time"

	"github.com/trendscout/backend/internal/models"
)

// RunScrapeTest は指定されたキーワードに対してスクレイピングをテストする関数です
// この関数は本番環境では使用せず、テスト目的でのみ使用してください
func RunScrapeTest(keyword string) error {
	// データベース接続を初期化
	if err := models.InitDatabases(); err != nil {
		return fmt.Errorf("failed to initialize databases: %w", err)
	}
	// 終了時にデータベース接続を閉じる
	defer models.CloseDatabases()

	// スクレイパーサービスを初期化
	service := NewService()

	// コンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 既存のキーワードをチェック
	keyword = models.SanitizeString(keyword)
	existingKeyword, err := models.GetKeywordByName(ctx, keyword)
	if err != nil {
		return fmt.Errorf("failed to check existing keyword: %w", err)
	}

	var userID int
	if existingKeyword != nil {
		fmt.Printf("Using existing keyword: %s (ID: %d)\n", existingKeyword.Keyword, existingKeyword.ID)
		userID = existingKeyword.UserID
	} else {
		// テスト用ユーザーを作成（または既存のテストユーザーを検索）
		testUser, err := models.CreateTestUser(ctx)
		if err != nil {
			return fmt.Errorf("failed to create test user: %w", err)
		}
		userID = testUser.ID

		// キーワードを作成
		keywordObj, err := models.CreateKeyword(ctx, userID, keyword)
		if err != nil {
			return fmt.Errorf("failed to create keyword: %w", err)
		}
		fmt.Printf("Created new keyword: %s (ID: %d)\n", keywordObj.Keyword, keywordObj.ID)
	}

	// スクレイピングを実行
	fmt.Printf("Starting scraping for keyword: %s...\n", keyword)
	startTime := time.Now()

	items, err := service.ScrapeKeyword(ctx, keyword)
	if err != nil {
		return fmt.Errorf("scraping failed: %w", err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Scraping completed in %s\n", elapsed)
	fmt.Printf("Found %d items\n", len(items))

	// 結果のサマリーを表示
	if len(items) > 0 {
		fmt.Println("\nSample results:")
		for i, item := range items {
			if i >= 5 {
				break // 最初の5件だけ表示
			}
			fmt.Printf("---\n")
			fmt.Printf("Title: %s\n", item.Title)
			fmt.Printf("Source: %s\n", item.Source)
			fmt.Printf("URL: %s\n", item.URL)
			fmt.Printf("Image: %s\n", item.ImageURL)
			fmt.Printf("Date: %s\n", item.PublishedAt.Format("2006-01-02"))
			fmt.Printf("Tags: %v\n", item.Tags)
		}
	} else {
		fmt.Println("No items found for the keyword.")
	}

	return nil
} 