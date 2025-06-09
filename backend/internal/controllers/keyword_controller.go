package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trendscout/backend/internal/auth"
	"github.com/trendscout/backend/internal/models"
	"github.com/trendscout/backend/internal/scraper"
	"github.com/trendscout/backend/internal/views"
)

// KeywordController handles keyword-related requests
type KeywordController struct {
	scraperService *scraper.Service
}

// NewKeywordController creates a new keyword controller
func NewKeywordController() *KeywordController {
	return &KeywordController{
		scraperService: scraper.NewService(),
	}
}

// KeywordCreateRequest represents the request for creating a keyword
type KeywordCreateRequest struct {
	Keyword string `json:"keyword" binding:"required,min=1,max=100"`
}

// GetKeywords handles retrieving all keywords for the authenticated user
func (c *KeywordController) GetKeywords(ctx *gin.Context) {
	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get keywords from database
	keywords, err := models.GetKeywordsForUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get keywords"})
		return
	}

	// Return keywords using view
	ctx.JSON(http.StatusOK, views.NewKeywordListResponse(keywords))
}

// CreateKeyword handles creating a new keyword
func (c *KeywordController) CreateKeyword(ctx *gin.Context) {
	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req KeywordCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create keyword in database
	keyword, err := models.CreateKeyword(ctx, userID, req.Keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create keyword"})
		return
	}

	// Start historical data collection in background
	go c.collectHistoricalData(keyword.ID, keyword.Keyword)

	// Return created keyword
	ctx.JSON(http.StatusCreated, views.NewKeywordResponse(keyword))
}

// collectHistoricalData collects historical data for a new keyword
func (c *KeywordController) collectHistoricalData(keywordID int, keywordText string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Collect data for the past 90 days
	now := time.Now()
	startDate := now.AddDate(0, 0, -90)

	// Simulate historical data collection by running multiple scraping sessions
	// In a real implementation, you might want to use different date ranges or historical APIs
	for days := 0; days < 90; days += 7 { // Collect data weekly for past 90 days
		targetDate := startDate.AddDate(0, 0, days)
		
		items, err := c.scraperService.ScrapeKeyword(ctx, keywordText)
		if err != nil {
			continue // Skip failed collections
		}

		// Store items with backdated timestamps
		if err := c.storeHistoricalItems(ctx, keywordID, items, targetDate); err != nil {
			continue // Skip failed storage
		}

		// Rate limiting to avoid overwhelming the sources
		time.Sleep(2 * time.Second)
	}
}

// storeHistoricalItems stores scraped items with historical dates
func (c *KeywordController) storeHistoricalItems(ctx context.Context, keywordID int, items []scraper.ScrapedItem, targetDate time.Time) error {
	// Group items by simulated date intervals
	itemsByDate := make(map[time.Time][]scraper.ScrapedItem)
	
	// Distribute items across the week leading up to targetDate
	for i, item := range items {
		// Distribute items across 7 days before target date
		dateOffset := i % 7
		itemDate := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day()-dateOffset, 0, 0, 0, 0, time.UTC)
		itemsByDate[itemDate] = append(itemsByDate[itemDate], item)
	}

	// Store items in MongoDB and update trend records
	for date, dateItems := range itemsByDate {
		// Skip future dates
		if date.After(time.Now()) {
			continue
		}

		// Calculate volume and sentiment
		volume := len(dateItems)
		sentiment := calculateSentiment(dateItems)

		// Store trend record in PostgreSQL
		_, err := models.CreateTrendRecord(ctx, keywordID, date, volume, sentiment)
		if err != nil {
			continue // Skip duplicate dates
		}

		// Store items in MongoDB with historical dates
		for _, item := range dateItems {
			image := &models.Image{
				KeywordID: keywordID,
				ImageURL:  item.ImageURL,
				Caption:   item.Title + " - " + item.Content,
				Tags:      item.Tags,
				FetchedAt: date, // Use historical date instead of current time
			}

			models.CreateImage(ctx, image) // Ignore errors for historical data
		}
	}

	return nil
}

// calculateSentiment calculates sentiment from scraped items
func calculateSentiment(items []scraper.ScrapedItem) float64 {
	if len(items) == 0 {
		return 0.5 // neutral
	}
	
	positiveWords := []string{"amazing", "beautiful", "stunning", "gorgeous", "elegant", "chic", "trendy", "stylish", "love", "perfect", "fabulous", "hot", "cool", "awesome"}
	negativeWords := []string{"ugly", "terrible", "awful", "disappointing", "boring", "outdated", "hate", "worst", "bad", "horrible"}
	
	var totalScore float64
	var count int
	
	for _, item := range items {
		content := strings.ToLower(item.Title + " " + item.Content)
		var score float64 = 0.5 // neutral baseline
		
		for _, word := range positiveWords {
			if strings.Contains(content, word) {
				score += 0.1
			}
		}
		
		for _, word := range negativeWords {
			if strings.Contains(content, word) {
				score -= 0.1
			}
		}
		
		// Clamp score between 0 and 1
		if score > 1.0 {
			score = 1.0
		}
		if score < 0.0 {
			score = 0.0
		}
		
		totalScore += score
		count++
	}
	
	if count == 0 {
		return 0.5
	}
	
	return totalScore / float64(count)
}

// UpdateKeyword handles updating a keyword
func (c *KeywordController) UpdateKeyword(ctx *gin.Context) {
	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse keyword ID from URL
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID"})
		return
	}

	// Check if keyword exists and belongs to user
	keyword, err := models.GetKeywordByID(ctx, id)
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

	// Parse request body
	var req KeywordCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update keyword in database
	if err := models.UpdateKeyword(ctx, id, req.Keyword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update keyword"})
		return
	}

	// Get updated keyword
	updatedKeyword, err := models.GetKeywordByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated keyword"})
		return
	}

	// Return updated keyword
	ctx.JSON(http.StatusOK, views.NewKeywordResponse(updatedKeyword))
}

// DeleteKeyword handles deleting a keyword
func (c *KeywordController) DeleteKeyword(ctx *gin.Context) {
	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse keyword ID from URL
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID"})
		return
	}

	// Check if keyword exists and belongs to user
	keyword, err := models.GetKeywordByID(ctx, id)
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

	// Delete keyword from database
	if err := models.DeleteKeyword(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete keyword"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Keyword deleted successfully"})
} 