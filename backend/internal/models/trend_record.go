package models

import (
	"context"
	"time"
)

// TrendRecord represents a trend data point
type TrendRecord struct {
	ID        int       `json:"id" db:"id"`
	KeywordID int       `json:"keyword_id" db:"keyword_id"`
	Date      time.Time `json:"date" db:"date"`
	Volume    int       `json:"volume" db:"volume"`
	Sentiment float64   `json:"sentiment" db:"sentiment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTrendRecord creates a new trend record
func CreateTrendRecord(ctx context.Context, keywordID int, date time.Time, volume int, sentiment float64) (*TrendRecord, error) {
	query := `
		INSERT INTO trend_records (keyword_id, record_date, volume, sentiment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (keyword_id, record_date) DO UPDATE SET
			volume = EXCLUDED.volume,
			sentiment = EXCLUDED.sentiment,
			updated_at = EXCLUDED.updated_at
		RETURNING id, keyword_id, record_date, volume, sentiment, created_at, updated_at
	`

	now := time.Now()
	var record TrendRecord

	err := PgPool.QueryRow(ctx, query, keywordID, date, volume, sentiment, now, now).
		Scan(&record.ID, &record.KeywordID, &record.Date, &record.Volume, &record.Sentiment, &record.CreatedAt, &record.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

// GetTrendRecords retrieves trend records for a keyword within a date range
func GetTrendRecords(ctx context.Context, keywordID int, startDate, endDate time.Time) ([]TrendRecord, error) {
	query := `
		SELECT id, keyword_id, record_date, volume, sentiment, created_at, updated_at
		FROM trend_records
		WHERE keyword_id = $1 AND record_date >= $2 AND record_date <= $3
		ORDER BY record_date ASC
	`

	rows, err := PgPool.Query(ctx, query, keywordID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []TrendRecord
	for rows.Next() {
		var record TrendRecord
		err := rows.Scan(&record.ID, &record.KeywordID, &record.Date, &record.Volume, &record.Sentiment, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

// GetTrendRecordsForKeyword retrieves all trend records for a keyword (legacy function)
func GetTrendRecordsForKeyword(ctx context.Context, keywordID int, startDate, endDate time.Time) ([]TrendRecord, error) {
	return GetTrendRecords(ctx, keywordID, startDate, endDate)
}

// GetLatestTrendRecord gets the most recent trend record for a keyword
func GetLatestTrendRecord(ctx context.Context, keywordID int) (*TrendRecord, error) {
	query := `
		SELECT id, keyword_id, record_date, volume, sentiment, created_at, updated_at
		FROM trend_records
		WHERE keyword_id = $1
		ORDER BY record_date DESC
		LIMIT 1
	`

	var record TrendRecord
	err := PgPool.QueryRow(ctx, query, keywordID).
		Scan(&record.ID, &record.KeywordID, &record.Date, &record.Volume, &record.Sentiment, &record.CreatedAt, &record.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

// GetTrendRecordsByDateRange gets all trend records within a date range
func GetTrendRecordsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]TrendRecord, error) {
	query := `
		SELECT id, keyword_id, record_date, volume, sentiment, created_at, updated_at
		FROM trend_records
		WHERE record_date >= $1 AND record_date <= $2
		ORDER BY record_date ASC, keyword_id ASC
	`

	rows, err := PgPool.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []TrendRecord
	for rows.Next() {
		var record TrendRecord
		err := rows.Scan(&record.ID, &record.KeywordID, &record.Date, &record.Volume, &record.Sentiment, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

// DeleteTrendRecordsForKeyword deletes all trend records for a keyword
func DeleteTrendRecordsForKeyword(ctx context.Context, keywordID int) error {
	query := `DELETE FROM trend_records WHERE keyword_id = $1`
	_, err := PgPool.Exec(ctx, query, keywordID)
	return err
}

// GetTrendStatistics calculates statistics for a keyword's trend data
func GetTrendStatistics(ctx context.Context, keywordID int, days int) (map[string]interface{}, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	query := `
		SELECT 
			COUNT(*) as data_points,
			AVG(volume) as avg_volume,
			MAX(volume) as max_volume,
			MIN(volume) as min_volume,
			AVG(sentiment) as avg_sentiment,
			MAX(sentiment) as max_sentiment,
			MIN(sentiment) as min_sentiment
		FROM trend_records
		WHERE keyword_id = $1 AND record_date >= $2 AND record_date <= $3
	`

	var stats struct {
		DataPoints   int     `db:"data_points"`
		AvgVolume    float64 `db:"avg_volume"`
		MaxVolume    int     `db:"max_volume"`
		MinVolume    int     `db:"min_volume"`
		AvgSentiment float64 `db:"avg_sentiment"`
		MaxSentiment float64 `db:"max_sentiment"`
		MinSentiment float64 `db:"min_sentiment"`
	}

	err := PgPool.QueryRow(ctx, query, keywordID, startDate, endDate).
		Scan(&stats.DataPoints, &stats.AvgVolume, &stats.MaxVolume, &stats.MinVolume,
			&stats.AvgSentiment, &stats.MaxSentiment, &stats.MinSentiment)

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"data_points":    stats.DataPoints,
		"avg_volume":     stats.AvgVolume,
		"max_volume":     stats.MaxVolume,
		"min_volume":     stats.MinVolume,
		"avg_sentiment":  stats.AvgSentiment,
		"max_sentiment":  stats.MaxSentiment,
		"min_sentiment":  stats.MinSentiment,
	}

	return result, nil
}

// SaveTrendData saves trend data from scraped posts and articles
func SaveTrendData(ctx context.Context, keywordID int, posts []SocialMediaPost, articles []BlogArticle) error {
	// Calculate volume based on number of posts and articles
	volume := len(posts) + len(articles)
	
	// Calculate basic sentiment (placeholder implementation)
	sentiment := 0.5 // neutral sentiment as default
	
	// Use current date
	date := time.Now().Truncate(24 * time.Hour) // Remove time component
	
	// Create trend record
	_, err := CreateTrendRecord(ctx, keywordID, date, volume, sentiment)
	return err
}

// GetLatestTrendRecords gets the most recent trend records for a keyword
func GetLatestTrendRecords(ctx context.Context, keywordID int, limit int) ([]TrendRecord, error) {
	query := `
		SELECT id, keyword_id, record_date, volume, sentiment, created_at, updated_at
		FROM trend_records
		WHERE keyword_id = $1
		ORDER BY record_date ASC
		LIMIT $2
	`

	rows, err := PgPool.Query(ctx, query, keywordID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []TrendRecord
	for rows.Next() {
		var record TrendRecord
		err := rows.Scan(&record.ID, &record.KeywordID, &record.Date, &record.Volume, &record.Sentiment, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, rows.Err()
} 