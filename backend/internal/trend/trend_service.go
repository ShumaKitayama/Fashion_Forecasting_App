package trend

import (
	"context"
	"time"

	"github.com/trendscout/backend/internal/models"
)

// Service handles trend analysis and predictions
type Service struct {
	engine *PredictionEngine
}

// NewService creates a new trend service
func NewService() *Service {
	return &Service{
		engine: NewPredictionEngine(),
	}
}

// TrendPrediction represents a single prediction point
type TrendPrediction struct {
	Date      time.Time `json:"date"`
	Volume    int       `json:"volume"`
	Sentiment float64   `json:"sentiment"`
}

// PredictTrend predicts future trend values for a keyword
func (s *Service) PredictTrend(ctx context.Context, keywordID int, days int) ([]TrendPrediction, error) {
	// Get historical data (last 30 days for better prediction)
	historicalRecords, err := models.GetLatestTrendRecords(ctx, keywordID, 30)
	if err != nil {
		return nil, err
	}

	if len(historicalRecords) == 0 {
		return []TrendPrediction{}, nil
	}

	// Convert to TrendPoint format
	var historical []TrendPoint
	for _, record := range historicalRecords {
		historical = append(historical, TrendPoint{
			Date:      record.Date,
			Volume:    float64(record.Volume),
			Sentiment: record.Sentiment,
		})
	}

	// Perform prediction
	predictions, err := s.engine.PredictTrend(historical, days)
	if err != nil {
		return nil, err
	}

	// Convert to TrendPrediction format
	var result []TrendPrediction
	for _, pred := range predictions {
		result = append(result, TrendPrediction{
			Date:      pred.Date,
			Volume:    int(pred.Volume),
			Sentiment: pred.Sentiment,
		})
	}

	return result, nil
}

// GetTrendAnalysis provides detailed trend analysis for a keyword
func (s *Service) GetTrendAnalysis(ctx context.Context, keywordID int) (map[string]interface{}, error) {
	// Get historical data for analysis
	records, err := models.GetLatestTrendRecords(ctx, keywordID, 30)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return map[string]interface{}{
			"status":      "no_data",
			"data_points": 0,
		}, nil
	}

	// Convert to TrendPoint format for analysis
	var historical []TrendPoint
	for _, record := range records {
		historical = append(historical, TrendPoint{
			Date:      record.Date,
			Volume:    float64(record.Volume),
			Sentiment: record.Sentiment,
		})
	}

	// Get predictions for insight generation
	predictions, err := s.engine.PredictTrend(historical, 7)
	if err != nil {
		predictions = []EnhancedPredictionResult{} // Continue without predictions if error
	}

	// Generate insights
	insights := s.engine.GetTrendInsights(historical, predictions)
	
	return insights, nil
} 