package trend

import (
	"fmt"
	"math"
	"time"
)

// PredictionEngine handles trend prediction calculations
type PredictionEngine struct{}

// NewPredictionEngine creates a new prediction engine
func NewPredictionEngine() *PredictionEngine {
	return &PredictionEngine{}
}

// TrendPoint represents a single data point in the trend
type TrendPoint struct {
	Date      time.Time
	Volume    float64
	Sentiment float64
}

// EnhancedPredictionResult represents a predicted value
type EnhancedPredictionResult struct {
	Date            time.Time
	Volume          float64
	Sentiment       float64
	Confidence      float64
	TrendDirection  string
	SeasonalFactor  float64
}

// PredictTrend performs advanced trend prediction using multiple algorithms
func (e *PredictionEngine) PredictTrend(historical []TrendPoint, horizon int) ([]EnhancedPredictionResult, error) {
	if len(historical) < 7 {
		return nil, fmt.Errorf("insufficient data for prediction (minimum 7 points required)")
	}

	// Sort data by date
	sortByDate(historical)

	// Perform different types of predictions
	volumePredictions := e.predictVolume(historical, horizon)
	sentimentPredictions := e.predictSentiment(historical, horizon)
	trendAnalysis := e.analyzeTrend(historical)
	seasonalFactors := e.calculateSeasonalFactors(historical)

	// Combine predictions
	var results []EnhancedPredictionResult
	baseDate := historical[len(historical)-1].Date

	for i := 0; i < horizon; i++ {
		predDate := baseDate.AddDate(0, 0, i+1)
		seasonalIndex := i % len(seasonalFactors)
		
		result := EnhancedPredictionResult{
			Date:            predDate,
			Volume:          volumePredictions[i],
			Sentiment:       sentimentPredictions[i],
			Confidence:      e.calculateConfidence(historical, i),
			TrendDirection:  trendAnalysis,
			SeasonalFactor:  seasonalFactors[seasonalIndex],
		}

		// Apply seasonal adjustment
		result.Volume *= result.SeasonalFactor

		results = append(results, result)
	}

	return results, nil
}

// predictVolume predicts volume using Exponential Smoothing with Trend
func (e *PredictionEngine) predictVolume(data []TrendPoint, horizon int) []float64 {
	if len(data) < 2 {
		return make([]float64, horizon)
	}

	// Holt's Exponential Smoothing parameters
	alpha := 0.3 // Level smoothing parameter
	beta := 0.1  // Trend smoothing parameter

	// Initialize
	level := data[0].Volume
	trend := data[1].Volume - data[0].Volume

	// Apply Holt's method to historical data
	for i := 1; i < len(data); i++ {
		prevLevel := level
		level = alpha*data[i].Volume + (1-alpha)*(level+trend)
		trend = beta*(level-prevLevel) + (1-beta)*trend
	}

	// Generate predictions
	predictions := make([]float64, horizon)
	for i := 0; i < horizon; i++ {
		predictions[i] = level + float64(i+1)*trend
		
		// Ensure non-negative values
		if predictions[i] < 0 {
			predictions[i] = 0
		}
	}

	return predictions
}

// predictSentiment predicts sentiment using moving average with momentum
func (e *PredictionEngine) predictSentiment(data []TrendPoint, horizon int) []float64 {
	if len(data) < 3 {
		predictions := make([]float64, horizon)
		for i := range predictions {
			predictions[i] = 0.5 // neutral
		}
		return predictions
	}

	// Calculate recent trend and momentum
	recentWindow := 7
	if len(data) < recentWindow {
		recentWindow = len(data)
	}

	recentData := data[len(data)-recentWindow:]
	
	// Calculate weighted moving average (more weight to recent data)
	weightSum := 0.0
	valueSum := 0.0
	
	for i, point := range recentData {
		weight := float64(i + 1) // Linear weight increase
		weightSum += weight
		valueSum += point.Sentiment * weight
	}
	
	avgSentiment := valueSum / weightSum

	// Calculate momentum (rate of change)
	momentum := 0.0
	if len(recentData) >= 3 {
		first := recentData[0].Sentiment
		last := recentData[len(recentData)-1].Sentiment
		momentum = (last - first) / float64(len(recentData)-1)
	}

	// Generate predictions with dampening
	predictions := make([]float64, horizon)
	dampening := 0.95 // Reduce momentum impact over time

	for i := 0; i < horizon; i++ {
		predicted := avgSentiment + momentum*float64(i+1)*math.Pow(dampening, float64(i))
		
		// Clamp between 0 and 1
		if predicted > 1.0 {
			predicted = 1.0
		} else if predicted < 0.0 {
			predicted = 0.0
		}
		
		predictions[i] = predicted
	}

	return predictions
}

// analyzeTrend determines the overall trend direction
func (e *PredictionEngine) analyzeTrend(data []TrendPoint) string {
	if len(data) < 3 {
		return "stable"
	}

	// Calculate slope using linear regression
	n := float64(len(data))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, point := range data {
		x := float64(i)
		y := point.Volume
		
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	// Classify trend based on slope
	if slope > 2.0 {
		return "strong_upward"
	} else if slope > 0.5 {
		return "upward"
	} else if slope > -0.5 {
		return "stable"
	} else if slope > -2.0 {
		return "downward"
	} else {
		return "strong_downward"
	}
}

// calculateSeasonalFactors calculates weekly seasonal patterns
func (e *PredictionEngine) calculateSeasonalFactors(data []TrendPoint) []float64 {
	// Default seasonal factors (Sunday to Saturday)
	// Fashion trends often peak on weekends and Wednesdays
	defaultFactors := []float64{1.1, 0.9, 0.95, 1.05, 1.0, 1.1, 1.15}

	if len(data) < 14 { // Need at least 2 weeks for seasonal analysis
		return defaultFactors
	}

	// Calculate day-of-week averages
	dayAverages := make([]float64, 7)
	dayCounts := make([]int, 7)

	for _, point := range data {
		dayOfWeek := int(point.Date.Weekday())
		dayAverages[dayOfWeek] += point.Volume
		dayCounts[dayOfWeek]++
	}

	// Calculate averages and normalize
	overallAverage := 0.0
	validDays := 0

	for i := 0; i < 7; i++ {
		if dayCounts[i] > 0 {
			dayAverages[i] /= float64(dayCounts[i])
			overallAverage += dayAverages[i]
			validDays++
		}
	}

	if validDays == 0 {
		return defaultFactors
	}

	overallAverage /= float64(validDays)

	// Normalize to create seasonal factors
	seasonalFactors := make([]float64, 7)
	for i := 0; i < 7; i++ {
		if dayCounts[i] > 0 && overallAverage > 0 {
			seasonalFactors[i] = dayAverages[i] / overallAverage
		} else {
			seasonalFactors[i] = defaultFactors[i]
		}

		// Clamp factors to reasonable range
		if seasonalFactors[i] < 0.7 {
			seasonalFactors[i] = 0.7
		} else if seasonalFactors[i] > 1.5 {
			seasonalFactors[i] = 1.5
		}
	}

	return seasonalFactors
}

// calculateConfidence determines prediction confidence based on data quality
func (e *PredictionEngine) calculateConfidence(data []TrendPoint, predictionDay int) float64 {
	baseConfidence := 0.8

	// Reduce confidence based on prediction distance
	distanceDecay := math.Exp(-float64(predictionDay) * 0.1)
	
	// Reduce confidence based on data variability
	if len(data) > 2 {
		variance := e.calculateVariance(data)
		variabilityFactor := math.Max(0.3, 1.0-variance*0.01)
		baseConfidence *= variabilityFactor
	}

	// Reduce confidence based on data recency
	if len(data) > 0 {
		lastDataAge := time.Since(data[len(data)-1].Date).Hours() / 24
		recencyFactor := math.Max(0.5, 1.0-lastDataAge*0.05)
		baseConfidence *= recencyFactor
	}

	return math.Max(0.1, math.Min(0.95, baseConfidence*distanceDecay))
}

// calculateVariance calculates the variance in volume data
func (e *PredictionEngine) calculateVariance(data []TrendPoint) float64 {
	if len(data) < 2 {
		return 0
	}

	mean := 0.0
	for _, point := range data {
		mean += point.Volume
	}
	mean /= float64(len(data))

	variance := 0.0
	for _, point := range data {
		diff := point.Volume - mean
		variance += diff * diff
	}
	variance /= float64(len(data) - 1)

	return math.Sqrt(variance)
}

// sortByDate sorts trend points by date
func sortByDate(data []TrendPoint) {
	for i := 0; i < len(data)-1; i++ {
		for j := 0; j < len(data)-i-1; j++ {
			if data[j].Date.After(data[j+1].Date) {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}

// GetTrendInsights provides detailed trend analysis and insights
func (e *PredictionEngine) GetTrendInsights(data []TrendPoint, predictions []EnhancedPredictionResult) map[string]interface{} {
	insights := make(map[string]interface{})

	if len(data) == 0 {
		return insights
	}

	// Overall metrics
	totalVolume := 0.0
	avgSentiment := 0.0
	for _, point := range data {
		totalVolume += point.Volume
		avgSentiment += point.Sentiment
	}
	avgSentiment /= float64(len(data))

	// Growth metrics
	if len(data) >= 7 {
		recentWeek := data[len(data)-7:]
		weekVolume := 0.0
		for _, point := range recentWeek {
			weekVolume += point.Volume
		}

		if len(data) >= 14 {
			prevWeek := data[len(data)-14 : len(data)-7]
			prevWeekVolume := 0.0
			for _, point := range prevWeek {
				prevWeekVolume += point.Volume
			}

			if prevWeekVolume > 0 {
				weeklyGrowth := ((weekVolume - prevWeekVolume) / prevWeekVolume) * 100
				insights["weekly_growth"] = weeklyGrowth
			}
		}
	}

	// Prediction summary
	if len(predictions) > 0 {
		avgPredictedVolume := 0.0
		avgConfidence := 0.0
		for _, pred := range predictions {
			avgPredictedVolume += pred.Volume
			avgConfidence += pred.Confidence
		}
		avgPredictedVolume /= float64(len(predictions))
		avgConfidence /= float64(len(predictions))

		insights["predicted_avg_volume"] = avgPredictedVolume
		insights["prediction_confidence"] = avgConfidence
		insights["trend_direction"] = predictions[0].TrendDirection
	}

	insights["total_volume"] = totalVolume
	insights["avg_sentiment"] = avgSentiment
	insights["data_points"] = len(data)

	return insights
} 