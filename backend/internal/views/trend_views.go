package views

import (
	"github.com/trendscout/backend/internal/models"
)

// TrendAnalysisResponse represents the response for trend analysis
type TrendAnalysisResponse struct {
	KeywordID int                    `json:"keyword_id"`
	StartDate string                 `json:"start_date"`
	EndDate   string                 `json:"end_date"`
	Data      []models.TrendRecord   `json:"data"`
	Insights  map[string]interface{} `json:"insights"`
}

// PredictionData represents a single prediction data point
type PredictionData struct {
	Date           string  `json:"date"`
	Volume         int     `json:"volume"`
	Sentiment      float64 `json:"sentiment"`
	Confidence     float64 `json:"confidence"`
	TrendDirection string  `json:"trend_direction"`
}

// TrendPredictionResponse represents the response for trend prediction
type TrendPredictionResponse struct {
	KeywordID   int                    `json:"keyword_id"`
	Predictions []PredictionData       `json:"predictions"`
	Insights    map[string]interface{} `json:"insights"`
}

// SentimentAnalysisResponse represents the response for sentiment analysis
type SentimentAnalysisResponse struct {
	KeywordID        int                  `json:"keyword_id"`
	AverageSentiment float64              `json:"average_sentiment"`
	PositiveCount    int                  `json:"positive_count"`
	NegativeCount    int                  `json:"negative_count"`
	NeutralCount     int                  `json:"neutral_count"`
	Data             []models.TrendRecord `json:"data"`
	Images           []models.Image       `json:"images"`
}

// KeywordComparisonData represents trend data for a single keyword in comparison
type KeywordComparisonData struct {
	KeywordID    int                  `json:"keyword_id"`
	Keyword      string               `json:"keyword"`
	TotalVolume  int                  `json:"total_volume"`
	AvgSentiment float64              `json:"avg_sentiment"`
	DataPoints   int                  `json:"data_points"`
	Trends       []models.TrendRecord `json:"trends"`
}

// MultiKeywordComparisonResponse represents the response for multi-keyword comparison
type MultiKeywordComparisonResponse struct {
	Period    int                     `json:"period"`
	StartDate string                  `json:"start_date"`
	EndDate   string                  `json:"end_date"`
	Keywords  []KeywordComparisonData `json:"keywords"`
	Insights  map[string]interface{}  `json:"insights"`
}

// TrendRecordResponse represents a trend record response
type TrendRecordResponse struct {
	ID        int     `json:"id"`
	KeywordID int     `json:"keyword_id"`
	Date      string  `json:"date"`
	Volume    int     `json:"volume"`
	Sentiment float64 `json:"sentiment"`
}

// TrendRecordListResponse represents a list of trend records
type TrendRecordListResponse struct {
	Records []*TrendRecordResponse `json:"records"`
	Count   int                    `json:"count"`
}

// NewTrendRecordResponse creates a new trend record response
func NewTrendRecordResponse(record *models.TrendRecord) *TrendRecordResponse {
	return &TrendRecordResponse{
		ID:        record.ID,
		KeywordID: record.KeywordID,
		Date:      record.Date.Format("2006-01-02"),
		Volume:    record.Volume,
		Sentiment: record.Sentiment,
	}
}

// NewTrendRecordListResponse creates a new trend record list response
func NewTrendRecordListResponse(records []models.TrendRecord) TrendRecordListResponse {
	var responses []*TrendRecordResponse
	for _, record := range records {
		responses = append(responses, NewTrendRecordResponse(&record))
	}

	return TrendRecordListResponse{
		Records: responses,
		Count:   len(responses),
	}
}

// PredictionResult represents the result of a trend prediction (legacy compatibility)
type PredictionResult struct {
	Date      string  `json:"date"`
	Volume    float64 `json:"volume"`
	Sentiment float64 `json:"sentiment"`
}

// TrendPredictionLegacyResponse represents the legacy prediction response
type TrendPredictionLegacyResponse struct {
	KeywordID   int                `json:"keyword_id"`
	Horizon     int                `json:"horizon"`
	Predictions []PredictionResult `json:"predictions"`
}

// NewTrendPredictionResponse creates a new trend prediction response (legacy)
func NewTrendPredictionResponse(keywordID, horizon int, predictions []PredictionResult) TrendPredictionLegacyResponse {
	return TrendPredictionLegacyResponse{
		KeywordID:   keywordID,
		Horizon:     horizon,
		Predictions: predictions,
	}
}

// SentimentResult represents sentiment analysis result (legacy compatibility)
type SentimentResult struct {
	Date      string  `json:"date"`
	Sentiment float64 `json:"sentiment"`
	Volume    int     `json:"volume"`
}

// NewSentimentResponse creates a new sentiment response
func NewSentimentResponse(result SentimentResult) SentimentResult {
	return result
} 