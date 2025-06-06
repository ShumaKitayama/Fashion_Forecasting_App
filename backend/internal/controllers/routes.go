package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/trendscout/backend/internal/auth"
	"github.com/trendscout/backend/internal/trend"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(r *gin.Engine) {
	// Create services
	authService := auth.NewService()
	trendService := trend.NewService()

	// Create middleware
	authMiddleware := auth.NewMiddleware(authService)

	// Create controllers
	authController := NewAuthController(authService)
	keywordController := NewKeywordController()
	trendController := NewTrendController(trendService)

	// API routes
	api := r.Group("/api")
	{
		// Auth routes
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authController.Register)
			authRoutes.POST("/login", authController.Login)
			authRoutes.POST("/refresh", authController.Refresh)
			authRoutes.POST("/logout", authMiddleware.Authenticate(), authController.Logout)
		}

		// Authenticated routes
		authenticated := api.Group("")
		authenticated.Use(authMiddleware.Authenticate())
		{
			// Keyword routes
			keywordRoutes := authenticated.Group("/keywords")
			{
				keywordRoutes.GET("/", keywordController.GetKeywords)
				keywordRoutes.POST("/", keywordController.CreateKeyword)
				keywordRoutes.PUT("/:id", keywordController.UpdateKeyword)
				keywordRoutes.DELETE("/:id", keywordController.DeleteKeyword)
			}

			// Trend routes
			trendRoutes := authenticated.Group("/trends")
			{
				trendRoutes.GET("/", trendController.GetTrends)
				trendRoutes.POST("/predict", trendController.PredictTrends)
				trendRoutes.POST("/sentiment", trendController.GetSentiment)
			}
		}
	}
} 