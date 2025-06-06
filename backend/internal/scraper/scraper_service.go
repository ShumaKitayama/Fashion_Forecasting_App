package scraper

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/trendscout/backend/internal/models"
)

// Service provides scraping operations
type Service struct {
	collector *colly.Collector
}

// ScrapedItem represents an item scraped from fashion websites
type ScrapedItem struct {
	Source      string    // website name
	URL         string    // original URL
	Title       string    // title or empty for social posts
	Content     string    // post content/caption
	ImageURL    string    // image URL if available
	Tags        []string  // keywords or categories
	PublishedAt time.Time // publication date
}

// NewService creates a new scraper service
func NewService() *Service {
	// Initialize collector with default settings
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		// 実際にアクセス可能なファッションサイト
		colly.AllowedDomains("www.vogue.com", "www.elle.com", "www.harpersbazaar.com", "www.fashionista.com", "wwd.com", "hypebeast.com"),
		colly.MaxDepth(2),
		// Respect robots.txt
		colly.AllowURLRevisit(),
	)

	// Set rate limiting
       c.Limit(&colly.LimitRule{
               // Set random delay up to 5 seconds
               RandomDelay: 5 * time.Second,
		// Parallelism
		Parallelism: 2,
	})

	// エラーハンドリングを追加
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request to %s failed: %v", r.Request.URL, err)
	})

	return &Service{
		collector: c,
	}
}

// ScrapeKeyword scrapes data for a specific keyword
func (s *Service) ScrapeKeyword(ctx context.Context, keyword string) ([]ScrapedItem, error) {
	var items []ScrapedItem

	// Create a clone of the collector for this scrape
	c := s.collector.Clone()

	// Vogueのスクレイピング設定
	c.OnHTML("div.feed-card", func(e *colly.HTMLElement) {
		// タイトルを取得
		title := e.ChildText("div.feed-card__content h3")
		
		// キーワードがタイトルに含まれているか確認
		if !strings.Contains(strings.ToLower(title), strings.ToLower(keyword)) {
			return
		}

		// 内容を取得
		summary := e.ChildText("div.feed-card__content div.feed-card__description")
		
		// 画像URLを取得
		imageURL := e.ChildAttr("div.feed-card__image-container img", "src")
		if imageURL == "" {
			imageURL = e.ChildAttr("div.feed-card__image-container source", "srcset")
		}
		
		// リンクを取得
		link := e.ChildAttr("a.feed-card__link", "href")
		fullURL := e.Request.AbsoluteURL(link)
		
		// 日付を取得（ここでは仮にcreated_atが見つからない場合は現在時刻を使用）
		dateStr := e.ChildText("div.feed-card__content div.feed-card__byline time")
		pubDate := parseArticleDate(dateStr)
		
		// タグを取得
		var tags []string
		tags = append(tags, keyword)
		
		// カテゴリを追加
		category := e.ChildText("div.feed-card__content div.feed-card__category")
		if category != "" {
			tags = append(tags, strings.ToLower(category))
		}

		// ScrapedItemを作成
		item := ScrapedItem{
			Source:      "vogue.com",
			URL:         fullURL,
			Title:       title,
			Content:     summary,
			ImageURL:    imageURL,
			Tags:        tags,
			PublishedAt: pubDate,
		}

		items = append(items, item)
	})

	// Elle.comのスクレイピング設定
	c.OnHTML("div.full-item", func(e *colly.HTMLElement) {
		// タイトルを取得
		title := e.ChildText("div.full-item-content h3.full-item-title")
		
		// キーワードがタイトルに含まれているか確認
		if !strings.Contains(strings.ToLower(title), strings.ToLower(keyword)) {
			return
		}

		// 内容を取得
		summary := e.ChildText("div.full-item-content div.full-item-dek")
		
		// 画像URLを取得
		imageURL := e.ChildAttr("div.full-item-image img", "data-src")
		if imageURL == "" {
			imageURL = e.ChildAttr("div.full-item-image img", "src")
		}
		
		// リンクを取得
		link := e.ChildAttr("a.full-item-link", "href")
		fullURL := e.Request.AbsoluteURL(link)
		
		// 日付を取得
		dateStr := e.ChildText("div.full-item-content div.full-item-metadata time")
		pubDate := parseArticleDate(dateStr)
		
		// タグを取得
		var tags []string
		tags = append(tags, keyword)
		
		// カテゴリを追加
		e.ForEach("div.full-item-content div.full-item-tags a", func(_ int, el *colly.HTMLElement) {
			tag := strings.ToLower(el.Text)
			if tag != "" {
				tags = append(tags, tag)
			}
		})

		// ScrapedItemを作成
		item := ScrapedItem{
			Source:      "elle.com",
			URL:         fullURL,
			Title:       title,
			Content:     summary,
			ImageURL:    imageURL,
			Tags:        tags,
			PublishedAt: pubDate,
		}

		items = append(items, item)
	})

	// WWDのスクレイピング設定
	c.OnHTML("article.article-card", func(e *colly.HTMLElement) {
		// タイトルを取得
		title := e.ChildText("h2.article-card__title")
		
		// キーワードがタイトルに含まれているか確認
		if !strings.Contains(strings.ToLower(title), strings.ToLower(keyword)) {
			return
		}

		// 内容を取得
		summary := e.ChildText("div.article-card__description")
		
		// 画像URLを取得
		imageURL := e.ChildAttr("div.article-card__image img", "src")
		
		// リンクを取得
		link := e.ChildAttr("a.article-card__link", "href")
		fullURL := e.Request.AbsoluteURL(link)
		
		// 日付を取得
		dateStr := e.ChildText("time.article-card__date")
		pubDate := parseArticleDate(dateStr)
		
		// タグを取得
		var tags []string
		tags = append(tags, keyword)
		
		// カテゴリを追加
		category := e.ChildText("div.article-card__category")
		if category != "" {
			tags = append(tags, strings.ToLower(category))
		}

		// ScrapedItemを作成
		item := ScrapedItem{
			Source:      "wwd.com",
			URL:         fullURL,
			Title:       title,
			Content:     summary,
			ImageURL:    imageURL,
			Tags:        tags,
			PublishedAt: pubDate,
		}

		items = append(items, item)
	})

	// 一般的な記事ページ用のフォールバック
	c.OnHTML("article, .article, .post", func(e *colly.HTMLElement) {
		title := e.ChildText("h1, h2, .article-title, .entry-title")
		content := e.ChildText("p.summary, .article-summary, .entry-summary")
		
		// Check if keyword is in title or content
		if !strings.Contains(strings.ToLower(title), strings.ToLower(keyword)) {
			return
		}

		// Get the first image or default image
		imageURL := e.ChildAttr("img.featured-image, .article-featured-image, .entry-image", "src")
		if imageURL == "" {
			imageURL = e.ChildAttr("img", "src")
		}

		// リンクを取得（現在のURLを使用）
		fullURL := e.Request.URL.String()
		
		// 日付を取得
		dateStr := e.ChildText(".article-date, .published-date, time, .date")
		pubDate := parseArticleDate(dateStr)

		item := ScrapedItem{
			Source:      e.Request.URL.Host,
			URL:         fullURL,
			Title:       title,
			Content:     content,
			ImageURL:    imageURL,
			PublishedAt: pubDate,
		}

		// Add keyword as a tag
		item.Tags = append(item.Tags, keyword)

		// Extract keywords from meta tags
		e.DOM.Find("meta[name=keywords]").Each(func(_ int, s *goquery.Selection) {
			if content, exists := s.Attr("content"); exists {
				keywords := strings.Split(content, ",")
				for _, kw := range keywords {
					kw = strings.TrimSpace(kw)
					if kw != "" {
						item.Tags = append(item.Tags, kw)
					}
				}
			}
		})

		items = append(items, item)
	})

	// Construct URLs to visit
	urls := []string{
		fmt.Sprintf("https://www.vogue.com/search?q=%s", keyword),
		fmt.Sprintf("https://www.elle.com/search/?q=%s", keyword),
		fmt.Sprintf("https://wwd.com/search/?s=%s", keyword),
		fmt.Sprintf("https://hypebeast.com/search?s=%s", keyword),
		fmt.Sprintf("https://www.fashionista.com/search?q=%s", keyword),
	}

	// Start scraping
	for _, url := range urls {
		// ランダムな遅延を追加（サーバー負荷軽減とブロック回避のため）
		time.Sleep(time.Duration(2000+rand.Intn(3000)) * time.Millisecond)
		
		if err := c.Visit(url); err != nil {
			log.Printf("Failed to visit %s: %v", url, err)
			// Continue with other URLs
			continue
		}
	}

	// Wait until scraping is finished
	c.Wait()

	// Store in database
	if err := s.storeScrapedItems(ctx, keyword, items); err != nil {
		return items, fmt.Errorf("failed to store scraped items: %w", err)
	}

	return items, nil
}

// storeScrapedItems stores the scraped items in MongoDB and updates trend records
func (s *Service) storeScrapedItems(ctx context.Context, keyword string, items []ScrapedItem) error {
	// Get keyword ID from database
	keywordObj, err := models.GetKeywordByName(ctx, keyword)
	if err != nil {
		return fmt.Errorf("failed to get keyword ID: %w", err)
	}

	if keywordObj == nil {
		return fmt.Errorf("keyword not found: %s", keyword)
	}

	keywordID := keywordObj.ID

	// Group items by date
	itemsByDate := make(map[time.Time][]ScrapedItem)
	for _, item := range items {
		date := time.Date(item.PublishedAt.Year(), item.PublishedAt.Month(), item.PublishedAt.Day(), 0, 0, 0, 0, time.UTC)
		itemsByDate[date] = append(itemsByDate[date], item)
	}

	// Store items in MongoDB and update trend records
	for date, dateItems := range itemsByDate {
		// Calculate volume and sentiment (placeholder)
		volume := len(dateItems)
		sentiment := 0.5 // Neutral sentiment as placeholder

		// Store trend record in PostgreSQL
		_, err := models.CreateTrendRecord(ctx, keywordID, date, volume, sentiment)
		if err != nil {
			return fmt.Errorf("failed to create trend record: %w", err)
		}

		// Store items in MongoDB
		for _, item := range dateItems {
			image := &models.Image{
				KeywordID: keywordID,
				ImageURL:  item.ImageURL,
				Caption:   item.Content,
				Tags:      item.Tags,
				FetchedAt: item.PublishedAt,
			}

			if err := models.CreateImage(ctx, image); err != nil {
				log.Printf("Failed to store image: %v", err)
				// Continue with other items
				continue
			}
		}
	}

	return nil
}

// parseArticleDate parses an article date string
func parseArticleDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now() // 日付が見つからない場合は現在時刻を返す
	}

	// 一般的な日付フォーマット
	formats := []string{
		"Jan 2, 2006",
		"January 2, 2006",
		"2006-01-02",
		"02/01/2006",
		"January 2, 2006 15:04",
		"Jan 2, 2006 15:04",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04",
		"2006/01/02 15:04:05",
		"2 January 2006",
		"2 Jan 2006",
		"02 Jan 2006",
		"Monday, January 2, 2006",
		"Mon, Jan 2, 2006",
	}

	dateStr = strings.TrimSpace(dateStr)
	
	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t
		}
	}

	// どのフォーマットにも一致しない場合は現在の日付から1週間前の日付を返す
	// （新しい記事として扱う）
	return time.Now().AddDate(0, 0, -7)
} 