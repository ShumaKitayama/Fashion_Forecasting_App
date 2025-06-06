package models

import (
	"context"
	"strings"
	"time"
)

// TrendRecord represents a trend record in the database
type TrendRecord struct {
	ID         int64     `json:"id"`
	KeywordID  int       `json:"keyword_id"`
	RecordDate time.Time `json:"record_date"`
	Volume     int       `json:"volume"`
	Sentiment  float64   `json:"sentiment"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateTrendRecord adds a new trend record for a keyword
func CreateTrendRecord(ctx context.Context, keywordID int, recordDate time.Time, volume int, sentiment float64) (*TrendRecord, error) {
	var tr TrendRecord
	err := PgPool.QueryRow(ctx,
		`INSERT INTO trend_records (keyword_id, record_date, volume, sentiment) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (keyword_id, record_date) DO UPDATE 
		SET volume = $3, sentiment = $4
		RETURNING id, keyword_id, record_date, volume, sentiment, created_at`,
		keywordID, recordDate, volume, sentiment).Scan(
		&tr.ID, &tr.KeywordID, &tr.RecordDate, &tr.Volume, &tr.Sentiment, &tr.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &tr, nil
}

// GetTrendRecordsForKeyword retrieves all trend records for a specific keyword within a date range
func GetTrendRecordsForKeyword(ctx context.Context, keywordID int, fromDate, toDate time.Time) ([]*TrendRecord, error) {
	rows, err := PgPool.Query(ctx,
		`SELECT id, keyword_id, record_date, volume, sentiment, created_at 
		FROM trend_records 
		WHERE keyword_id = $1 AND record_date BETWEEN $2 AND $3 
		ORDER BY record_date ASC`,
		keywordID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*TrendRecord
	for rows.Next() {
		var tr TrendRecord
		if err := rows.Scan(&tr.ID, &tr.KeywordID, &tr.RecordDate, &tr.Volume, &tr.Sentiment, &tr.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, &tr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// GetTrendRecordsForDate retrieves all trend records for a specific date
func GetTrendRecordsForDate(ctx context.Context, date time.Time) ([]*TrendRecord, error) {
	rows, err := PgPool.Query(ctx,
		`SELECT id, keyword_id, record_date, volume, sentiment, created_at 
		FROM trend_records 
		WHERE record_date = $1
		ORDER BY keyword_id ASC`,
		date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*TrendRecord
	for rows.Next() {
		var tr TrendRecord
		if err := rows.Scan(&tr.ID, &tr.KeywordID, &tr.RecordDate, &tr.Volume, &tr.Sentiment, &tr.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, &tr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// GetLatestTrendRecords retrieves the most recent trend records for a keyword (up to limit count)
func GetLatestTrendRecords(ctx context.Context, keywordID int, limit int) ([]*TrendRecord, error) {
	if limit <= 0 {
		limit = 30 // Default to 30 days if not specified
	}

	rows, err := PgPool.Query(ctx,
		`SELECT id, keyword_id, record_date, volume, sentiment, created_at 
		FROM trend_records 
		WHERE keyword_id = $1 
		ORDER BY record_date DESC 
		LIMIT $2`,
		keywordID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*TrendRecord
	for rows.Next() {
		var tr TrendRecord
		if err := rows.Scan(&tr.ID, &tr.KeywordID, &tr.RecordDate, &tr.Volume, &tr.Sentiment, &tr.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, &tr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Reverse the slice to have chronological order
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}

	return records, nil
}

// SaveTrendData saves trend data to the database
func SaveTrendData(ctx context.Context, keywordID int, posts []SocialMediaPost, articles []BlogArticle) error {
	// 日付ごとのボリュームを集計
	volumeMap := make(map[string]int)
	sentimentMap := make(map[string]float64)
	sentimentCountMap := make(map[string]int)
	
	// 投稿データの集計
	for _, post := range posts {
		dateStr := post.PostDate.Format("2006-01-02")
		volumeMap[dateStr]++
		
		// センチメント計算（仮実装）
		sentimentScore := 0.0
		if strings.Contains(strings.ToLower(post.Caption), "love") || 
		   strings.Contains(strings.ToLower(post.Caption), "great") {
			sentimentScore = 0.8
		} else if strings.Contains(strings.ToLower(post.Caption), "hate") || 
		   strings.Contains(strings.ToLower(post.Caption), "bad") {
			sentimentScore = -0.8
		}
		
		sentimentMap[dateStr] += sentimentScore
		sentimentCountMap[dateStr]++
	}
	
	// 記事データの集計
	for _, article := range articles {
		dateStr := article.PublishDate.Format("2006-01-02")
		volumeMap[dateStr] += 3  // 記事は投稿より重み付け
	}
	
	// トレンドレコードの保存
	for dateStr, volume := range volumeMap {
		date, _ := time.Parse("2006-01-02", dateStr)
		
		// センチメントスコア計算
		sentiment := 0.0
		if count := sentimentCountMap[dateStr]; count > 0 {
			sentiment = sentimentMap[dateStr] / float64(count)
		}
		
		// レコード保存
		record := &TrendRecord{
			KeywordID:  keywordID,
			RecordDate: date,
			Volume:     volume,
			Sentiment:  sentiment,
		}
		
		if err := SaveTrendRecord(ctx, record); err != nil {
			return err
		}
	}
	
	return nil
}

// SaveTrendRecord saves a trend record to the database
func SaveTrendRecord(ctx context.Context, record *TrendRecord) error {
	_, err := PgPool.Exec(ctx,
		`INSERT INTO trend_records (keyword_id, record_date, volume, sentiment) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (keyword_id, record_date) DO UPDATE 
		SET volume = $3, sentiment = $4`,
		record.KeywordID, record.RecordDate, record.Volume, record.Sentiment)
	return err
} 