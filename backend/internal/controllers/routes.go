package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/trendscout/backend/internal/auth"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(r *gin.Engine) {
	// Initialize services
	authService := auth.NewService()

	// Initialize controllers
	authController := NewAuthController(authService)
	keywordController := NewKeywordController()
	trendController := NewTrendController()
	dataController := NewDataController()

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(authService)

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/auth/register", authController.Register)
		public.POST("/auth/login", authController.Login)
		public.POST("/auth/refresh", authController.Refresh)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(authMiddleware.Authenticate())
	{
		// Auth routes
		protected.POST("/auth/logout", authController.Logout)

		// Keyword routes
		protected.GET("/keywords", keywordController.GetKeywords)
		protected.POST("/keywords", keywordController.CreateKeyword)
		protected.PUT("/keywords/:id", keywordController.UpdateKeyword)
		protected.DELETE("/keywords/:id", keywordController.DeleteKeyword)

		// Trend routes
		protected.GET("/trends/", trendController.GetTrendData)
		protected.POST("/trends/analysis", trendController.GetTrendAnalysis)
		protected.POST("/trends/prediction", trendController.GetTrendPrediction)
		protected.POST("/trends/predict", trendController.GetTrendPrediction) // Legacy alias
		protected.POST("/trends/sentiment", trendController.GetSentimentAnalysis)
		protected.GET("/trends/comparison", trendController.GetMultiKeywordComparison)

		// Data collection routes
		protected.POST("/data/collect/:id", dataController.CollectKeywordData)
	}
} 