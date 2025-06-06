package controllers

import (
	"errors"
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
	trendService *trend.Service
}

// NewTrendController creates a new trend controller
func NewTrendController(trendService *trend.Service) *TrendController {
	return &TrendController{
		trendService: trendService,
	}
}

// TrendPredictRequest represents the request for trend prediction
type TrendPredictRequest struct {
	KeywordID int `json:"keyword_id" binding:"required"`
	Horizon   int `json:"horizon" binding:"required,min=1,max=30"`
}

// TrendSentimentRequest represents the request for sentiment analysis
type TrendSentimentRequest struct {
	KeywordID int       `json:"keyword_id" binding:"required"`
	Text      string    `json:"text" binding:"required,min=10"`
	Date      time.Time `json:"date" binding:"required" time_format:"2006-01-02"`
}

// GetTrends handles retrieving trend data for a keyword
func (c *TrendController) GetTrends(ctx *gin.Context) {
	// Parse query parameters
	keywordIDStr := ctx.Query("q")
	if keywordIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Keyword ID (q) is required"})
		return
	}

	keywordID, err := strconv.Atoi(keywordIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID"})
		return
	}

	fromStr := ctx.DefaultQuery("from", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	toStr := ctx.DefaultQuery("to", time.Now().Format("2006-01-02"))

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from date format (use YYYY-MM-DD)"})
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to date format (use YYYY-MM-DD)"})
		return
	}

	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if keyword belongs to user
	keyword, err := models.GetKeywordByID(ctx, keywordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get keyword"})
		return
	}

	if keyword == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return
	}

	if keyword.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Get trend records
	records, err := models.GetTrendRecordsForKeyword(ctx, keywordID, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trend data"})
		return
	}

	// Return trend data
	ctx.JSON(http.StatusOK, views.NewTrendRecordListResponse(records))
}

// PredictTrends handles trend prediction
func (c *TrendController) PredictTrends(ctx *gin.Context) {
	var req TrendPredictRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if keyword belongs to user
	keyword, err := models.GetKeywordByID(ctx, req.KeywordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get keyword"})
		return
	}

	if keyword == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return
	}

	if keyword.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Get trend predictions
	predictions, err := c.trendService.PredictTrend(ctx, req.KeywordID, req.Horizon)
	if err != nil {
		var statusCode int
		var message string

		switch {
		case errors.Is(err, trend.ErrInsufficientData):
			statusCode = http.StatusBadRequest
			message = "Insufficient data for prediction"
		case errors.Is(err, trend.ErrPredictionFailed):
			statusCode = http.StatusInternalServerError
			message = "Prediction failed"
		default:
			statusCode = http.StatusInternalServerError
			message = "Failed to predict trends"
		}

		ctx.JSON(statusCode, gin.H{"error": message})
		return
	}

	// Return predictions
	ctx.JSON(http.StatusOK, gin.H{"predictions": predictions})
}

// GetSentiment handles sentiment analysis
func (c *TrendController) GetSentiment(ctx *gin.Context) {
	var req TrendSentimentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if keyword belongs to user
	keyword, err := models.GetKeywordByID(ctx, req.KeywordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get keyword"})
		return
	}

	if keyword == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return
	}

	if keyword.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Get sentiment analysis
	sentiment, err := c.trendService.AnalyzeSentiment(ctx, req.KeywordID, req.Date)
	if err != nil {
		var statusCode int
		var message string

		switch {
		case errors.Is(err, trend.ErrDataNotFound):
			statusCode = http.StatusNotFound
			message = "No data available for the specified date"
		case errors.Is(err, trend.ErrAnalysisFailed):
			statusCode = http.StatusInternalServerError
			message = "Sentiment analysis failed"
		default:
			statusCode = http.StatusInternalServerError
			message = "Failed to analyze sentiment"
		}

		ctx.JSON(statusCode, gin.H{"error": message})
		return
	}

	// Return sentiment analysis
	ctx.JSON(http.StatusOK, sentiment)
} 