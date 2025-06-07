// backend/cmd/datacollector/data_collector.go
package datacollector

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/trendscout/backend/internal/collector"
	"github.com/trendscout/backend/internal/models"
)

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("Warning: .env.local file not found, trying .env file")
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

	log.Println("データ収集を開始します...")
	
	// データベース接続を初期化
	if err := models.InitDatabases(); err != nil {
		log.Fatalf("データベース初期化エラー: %v", err)
	}
	defer models.CloseDatabases()
	
	ctx := context.Background()
	
	// キーワードリストを取得
	keywords, err := models.GetAllKeywords(ctx)
	if err != nil {
		log.Fatalf("キーワード取得エラー: %v", err)
	}
	
	if len(keywords) == 0 {
		log.Println("キーワードが登録されていません。データ収集をスキップします。")
		return
	}
	
	log.Printf("%d件のキーワードに対してデータ収集を実行します", len(keywords))
	
	// 各キーワードに対してデータ収集を実行
	for _, keyword := range keywords {
		log.Printf("キーワード '%s' (ID: %d) のデータ収集を開始", keyword.Keyword, keyword.ID)
		
		// SNSデータの収集
		posts, err := collector.CollectSocialMediaData(ctx, keyword.Keyword)
		if err != nil {
			log.Printf("SNSデータ収集エラー: %v", err)
			continue
		}
		
		// キーワードIDを設定
		for i := range posts {
			posts[i].KeywordID = keyword.ID
		}
		
		// ブログデータの収集
		articles, err := collector.CollectBlogData(ctx, keyword.Keyword)
		if err != nil {
			log.Printf("ブログデータ収集エラー: %v", err)
		}
		
		// キーワードIDを設定
		for i := range articles {
			articles[i].KeywordID = keyword.ID
		}
		
		// データを保存
		if err := models.SaveTrendData(ctx, keyword.ID, posts, articles); err != nil {
			log.Printf("データ保存エラー: %v", err)
		} else {
			log.Printf("キーワード '%s' のデータを保存しました（投稿: %d件, 記事: %d件）", 
				keyword.Keyword, len(posts), len(articles))
		}
		
		// APIレート制限対策のため少し待機
		time.Sleep(2 * time.Second)
	}
	
	log.Println("データ収集が完了しました")
} 