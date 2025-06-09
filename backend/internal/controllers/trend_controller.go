package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trendscout/backend/internal/auth"
	"github.com/trendscout/backend/internal/models"
	"github.com/trendscout/backend/internal/trend"
	"github.com/trendscout/backend/internal/views"
)

// TrendController handles trend-related requests
type TrendController struct {
	predictionEngine *trend.PredictionEngine
}

// NewTrendController creates a new trend controller
func NewTrendController() *TrendController {
	return &TrendController{
		predictionEngine: trend.NewPredictionEngine(),
	}
}

// GetTrendData handles trend data retrieval requests
func (c *TrendController) GetTrendData(ctx *gin.Context) {
	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get query parameters
	keywordIDStr := ctx.Query("q")
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")

	if keywordIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "keyword_id (q) parameter is required"})
		return
	}

	keywordID, err := strconv.Atoi(keywordIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword_id"})
		return
	}

	// Verify keyword ownership
	if !c.verifyKeywordOwnership(ctx, keywordID, userID) {
		return
	}

	// Parse dates with defaults
	var startDate, endDate time.Time
	if fromStr != "" {
		startDate, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from date format"})
			return
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // Default: 30 days ago
	}

	if toStr != "" {
		endDate, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to date format"})
			return
		}
	} else {
		endDate = time.Now() // Default: today
	}

	// Get trend records
	records, err := models.GetTrendRecords(ctx, keywordID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trend data"})
		return
	}

	// Convert to response format using the helper function
	response := views.NewTrendRecordListResponse(records)

	ctx.JSON(http.StatusOK, response)
}

// TrendAnalysisRequest represents the request for trend analysis
type TrendAnalysisRequest struct {
	KeywordID int    `json:"keyword_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

// TrendPredictionRequest represents the request for trend prediction
type TrendPredictionRequest struct {
	KeywordID int `json:"keyword_id" binding:"required"`
	Days      int `json:"days" binding:"required,min=1,max=60"`
}

// TrendSentimentRequest represents the request for sentiment analysis
type TrendSentimentRequest struct {
	KeywordID int `json:"keyword_id" binding:"required"`
	Period    int `json:"period" binding:"required,min=1,max=90"`
}

// GetTrendAnalysis handles trend analysis requests
func (c *TrendController) GetTrendAnalysis(ctx *gin.Context) {
	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request
	var req TrendAnalysisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify keyword ownership
	if !c.verifyKeywordOwnership(ctx, req.KeywordID, userID) {
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	// Get trend data
	trends, err := models.GetTrendRecords(ctx, req.KeywordID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trend data"})
		return
	}

	// Convert to trend points for analysis
	var trendPoints []trend.TrendPoint
	for _, t := range trends {
		trendPoints = append(trendPoints, trend.TrendPoint{
			Date:      t.Date,
			Volume:    float64(t.Volume),
			Sentiment: t.Sentiment,
		})
	}

	// Generate insights
	var insights map[string]interface{}
	if len(trendPoints) > 0 {
		// Try to generate predictions for insights
		predictions, _ := c.predictionEngine.PredictTrend(trendPoints, 7)
		insights = c.predictionEngine.GetTrendInsights(trendPoints, predictions)
	} else {
		insights = make(map[string]interface{})
	}

	// Return response
	response := views.TrendAnalysisResponse{
		KeywordID: req.KeywordID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Data:      trends,
		Insights:  insights,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTrendPrediction handles trend prediction requests
func (c *TrendController) GetTrendPrediction(ctx *gin.Context) {
	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request
	var req TrendPredictionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify keyword ownership
	if !c.verifyKeywordOwnership(ctx, req.KeywordID, userID) {
		return
	}

	// Get historical data (last 90 days for better prediction)
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -90)

	trends, err := models.GetTrendRecords(ctx, req.KeywordID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get historical data"})
		return
	}

	if len(trends) < 7 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":          "Insufficient data for prediction. At least 7 data points required.",
			"available_points": len(trends),
		})
		return
	}

	// Convert to trend points
	var trendPoints []trend.TrendPoint
	for _, t := range trends {
		trendPoints = append(trendPoints, trend.TrendPoint{
			Date:      t.Date,
			Volume:    float64(t.Volume),
			Sentiment: t.Sentiment,
		})
	}

	// Generate predictions
	predictions, err := c.predictionEngine.PredictTrend(trendPoints, req.Days)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Prediction failed: %v", err)})
		return
	}

	// Convert predictions to response format
	var predictionData []views.PredictionData
	for _, pred := range predictions {
		predictionData = append(predictionData, views.PredictionData{
			Date:           pred.Date.Format("2006-01-02"),
			Volume:         int(pred.Volume),
			Sentiment:      pred.Sentiment,
			Confidence:     pred.Confidence,
			TrendDirection: pred.TrendDirection,
		})
	}

	// Generate insights
	insights := c.predictionEngine.GetTrendInsights(trendPoints, predictions)

	// Return response
	response := views.TrendPredictionResponse{
		KeywordID:   req.KeywordID,
		Predictions: predictionData,
		Insights:    insights,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetSentimentAnalysis handles sentiment analysis requests
func (c *TrendController) GetSentimentAnalysis(ctx *gin.Context) {
	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request
	var req TrendSentimentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify keyword ownership
	if !c.verifyKeywordOwnership(ctx, req.KeywordID, userID) {
		return
	}

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -req.Period)

	// Get trend data
	trends, err := models.GetTrendRecords(ctx, req.KeywordID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trend data"})
		return
	}

	// Get images for sentiment context
	images, err := models.GetImagesByKeywordAndDateRange(ctx, req.KeywordID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get images"})
		return
	}

	// Calculate sentiment statistics
	var totalSentiment float64
	var positiveCount, negativeCount, neutralCount int

	for _, trend := range trends {
		totalSentiment += trend.Sentiment
		if trend.Sentiment > 0.6 {
			positiveCount++
		} else if trend.Sentiment < 0.4 {
			negativeCount++
		} else {
			neutralCount++
		}
	}

	var averageSentiment float64
	if len(trends) > 0 {
		averageSentiment = totalSentiment / float64(len(trends))
	}

	// Return response
	response := views.SentimentAnalysisResponse{
		KeywordID:        req.KeywordID,
		AverageSentiment: averageSentiment,
		PositiveCount:    positiveCount,
		NegativeCount:    negativeCount,
		NeutralCount:     neutralCount,
		Data:             trends,
		Images:           images,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetMultiKeywordComparison handles multi-keyword comparison requests
func (c *TrendController) GetMultiKeywordComparison(ctx *gin.Context) {
	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get query parameters
	keywordIDsStr := ctx.Query("keyword_ids")
	daysStr := ctx.DefaultQuery("days", "30")

	if keywordIDsStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "keyword_ids parameter is required"})
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 90 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter (1-90)"})
		return
	}

	// Parse keyword IDs
	var keywordIDs []int
	for _, idStr := range []string{keywordIDsStr} {
		// TODO: Properly parse comma-separated IDs
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID"})
			return
		}
		keywordIDs = append(keywordIDs, id)
	}

	// Verify all keywords belong to user
	for _, keywordID := range keywordIDs {
		if !c.verifyKeywordOwnership(ctx, keywordID, userID) {
			return
		}
	}

	// Get comparison data
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	var comparisonData []views.KeywordComparisonData
	for _, keywordID := range keywordIDs {
		// Get keyword info
		keyword, err := models.GetKeywordByID(ctx, keywordID)
		if err != nil {
			continue
		}

		// Get trend data
		trends, err := models.GetTrendRecords(ctx, keywordID, startDate, endDate)
		if err != nil {
			continue
		}

		// Calculate metrics
		var totalVolume int
		var totalSentiment float64
		for _, trend := range trends {
			totalVolume += trend.Volume
			totalSentiment += trend.Sentiment
		}

		var avgSentiment float64
		if len(trends) > 0 {
			avgSentiment = totalSentiment / float64(len(trends))
		}

		comparisonData = append(comparisonData, views.KeywordComparisonData{
			KeywordID:    keywordID,
			Keyword:      keyword.Keyword,
			TotalVolume:  totalVolume,
			AvgSentiment: avgSentiment,
			DataPoints:   len(trends),
			Trends:       trends,
		})
	}

	// Return response
	response := views.MultiKeywordComparisonResponse{
		Period:     days,
		StartDate:  startDate.Format("2006-01-02"),
		EndDate:    endDate.Format("2006-01-02"),
		Keywords:   comparisonData,
		Insights:   c.generateComparisonInsights(comparisonData),
	}

	ctx.JSON(http.StatusOK, response)
}

// verifyKeywordOwnership checks if a keyword belongs to the user
func (c *TrendController) verifyKeywordOwnership(ctx *gin.Context, keywordID, userID int) bool {
	keyword, err := models.GetKeywordByID(ctx, keywordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify keyword ownership"})
		return false
	}

	if keyword == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return false
	}

	if keyword.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return false
	}

	return true
}

// generateComparisonInsights generates insights for multi-keyword comparison
func (c *TrendController) generateComparisonInsights(data []views.KeywordComparisonData) map[string]interface{} {
	insights := make(map[string]interface{})

	if len(data) == 0 {
		return insights
	}

	// Find highest performing keyword
	var topKeyword views.KeywordComparisonData
	maxVolume := 0

	for _, keyword := range data {
		if keyword.TotalVolume > maxVolume {
			maxVolume = keyword.TotalVolume
			topKeyword = keyword
		}
	}

	insights["top_keyword"] = topKeyword.Keyword
	insights["top_volume"] = maxVolume

	// Calculate average sentiment across all keywords
	var totalSentiment float64
	var sentimentCount int

	for _, keyword := range data {
		if keyword.DataPoints > 0 {
			totalSentiment += keyword.AvgSentiment
			sentimentCount++
		}
	}

	if sentimentCount > 0 {
		insights["overall_avg_sentiment"] = totalSentiment / float64(sentimentCount)
	}

	// Calculate total data points
	totalDataPoints := 0
	for _, keyword := range data {
		totalDataPoints += keyword.DataPoints
	}
	insights["total_data_points"] = totalDataPoints

	return insights
} 