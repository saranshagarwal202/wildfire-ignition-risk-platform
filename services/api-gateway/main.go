package main

import (
	"log"
	"net/http"

	"api-gateway/config"
	"api-gateway/handlers"
	"api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "api-gateway",
			"version": "1.0.0",
		})
	})

	riskHandler := handlers.NewRiskAssessmentHandler(cfg)

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/wildfire-risk-jobs", riskHandler.CreateJob)
		v1.GET("/wildfire-risk-jobs/:job_id", riskHandler.GetJobStatus)
		v1.GET("/wildfire-risk-jobs/:job_id/results", riskHandler.GetJobResults)
	}

	// Starting server
	log.Printf("Starting API Gateway on port %s", cfg.Port)
	log.Printf("Orchestrator URL: %s", cfg.OrchestratorURL)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
