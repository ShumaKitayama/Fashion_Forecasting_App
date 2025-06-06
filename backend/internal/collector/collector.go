package collector

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/trendscout/backend/internal/models"
)

// CollectSocialMediaData はSNSからデータを収集します
func CollectSocialMediaData(ctx context.Context, keyword string) ([]models.SocialMediaPost, error) {
	log.Printf("キーワード '%s' のSNSデータを収集中...", keyword)
	
	// ここに実際のAPI呼び出しコードを実装
	// テスト用のダミーデータを返す
	posts := []models.SocialMediaPost{
		{
			ID:           fmt.Sprintf("post_%d", time.Now().Unix()),
			KeywordID:    0, // 呼び出し元で設定される
			Platform:     "instagram",
			PostID:       fmt.Sprintf("ig_%d", time.Now().Unix()),
			Username:     "fashion_user",
			Caption:      fmt.Sprintf("Latest %s trends are amazing! #fashion #trend", keyword),
			ImageURL:     "https://example.com/image1.jpg",
			LikeCount:    100,
			CommentCount: 20,
			PostDate:     time.Now(),
			CreatedAt:    time.Now(),
		},
	}
	
	return posts, nil
}

// CollectBlogData はブログからデータを収集します
func CollectBlogData(ctx context.Context, keyword string) ([]models.BlogArticle, error) {
	log.Printf("キーワード '%s' のブログデータを収集中...", keyword)
	
	// ここに実際のRSS取得コードを実装
	// テスト用のダミーデータを返す
	articles := []models.BlogArticle{
		{
			ID:          fmt.Sprintf("blog_%d", time.Now().Unix()),
			KeywordID:   0, // 呼び出し元で設定される
			Title:       fmt.Sprintf("The Rise of %s in Fashion Industry", keyword),
			URL:         "https://example.com/fashion-blog/article1",
			Author:      "Fashion Expert",
			Content:     fmt.Sprintf("This season, %s is becoming increasingly popular...", keyword),
			PublishDate: time.Now(),
			CreatedAt:   time.Now(),
		},
	}
	
	return articles, nil
} 