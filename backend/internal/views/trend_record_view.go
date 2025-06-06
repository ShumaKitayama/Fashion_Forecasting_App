package views

import (
	"time"

	"github.com/trendscout/backend/internal/models"
)

// TrendRecordResponse represents the trend record data returned in API responses
type TrendRecordResponse struct {
	ID         int64     `json:"id"`
	KeywordID  int       `json:"keyword_id"`
	RecordDate time.Time `json:"date"`
	Volume     int       `json:"volume"`
	Sentiment  float64   `json:"sentiment"`
}

// NewTrendRecordResponse creates a new trend record response from a trend record model
func NewTrendRecordResponse(record *models.TrendRecord) *TrendRecordResponse {
	return &TrendRecordResponse{
		ID:         record.ID,
		KeywordID:  record.KeywordID,
		RecordDate: record.RecordDate,
		Volume:     record.Volume,
		Sentiment:  record.Sentiment,
	}
}

// TrendRecordListResponse represents a list of trend records
type TrendRecordListResponse struct {
	Records []*TrendRecordResponse `json:"records"`
	Count   int                    `json:"count"`
}

// NewTrendRecordListResponse creates a new trend record list response
func NewTrendRecordListResponse(records []*models.TrendRecord) *TrendRecordListResponse {
	responses := make([]*TrendRecordResponse, len(records))
	for i, record := range records {
		responses[i] = NewTrendRecordResponse(record)
	}

	return &TrendRecordListResponse{
		Records: responses,
		Count:   len(responses),
	}
}

// TrendPrediction represents a prediction point
type TrendPrediction struct {
	Date   time.Time `json:"date"`
	Volume int       `json:"volume"`
}

// TrendPredictionResponse represents the prediction response
type TrendPredictionResponse struct {
	Predictions []*TrendPrediction `json:"predictions"`
	KeywordID   int                `json:"keyword_id"`
	Horizon     int                `json:"horizon"`
}

// NewTrendPredictionResponse creates a new trend prediction response
func NewTrendPredictionResponse(keywordID int, horizon int, predictions []*TrendPrediction) *TrendPredictionResponse {
	return &TrendPredictionResponse{
		KeywordID:   keywordID,
		Horizon:     horizon,
		Predictions: predictions,
	}
}

// SentimentAnalysisResponse represents the sentiment analysis response
type SentimentAnalysisResponse struct {
	KeywordID int       `json:"keyword_id"`
	Date      time.Time `json:"date"`
	Positive  float64   `json:"positive"`
	Neutral   float64   `json:"neutral"`
	Negative  float64   `json:"negative"`
} 