package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/handlers"
	"github.com/Calif432/farmtrack/backend/internal/middleware"
)

func main() {
	log.Println("🚀 Starting FarmTrack Backend...")

	// Connect to database
	database.Connect()

	// Create default owner account
	handlers.CreateFirstOwner()

	// Create Gin router
	r := gin.Default()

	// CORS middleware (for React frontend)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Health check endpoint (public)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "FarmTrack API is running",
		})
	})

	// ============ Auth Routes (Public) ============
	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
	}

	// ============ Protected Routes ============
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// User profile (only once!)
		api.GET("/auth/me", handlers.GetCurrentUser)

		// Animal routes
		api.POST("/animals", handlers.CreateAnimal)
		api.GET("/animals", handlers.GetAnimals)
		api.GET("/animals/:id", handlers.GetAnimalByID)
		api.PUT("/animals/:id", handlers.UpdateAnimal)
		api.DELETE("/animals/:id", handlers.DeleteAnimal)

		// Vaccination routes
		api.POST("/vaccinations", handlers.CreateVaccination)
		api.GET("/vaccinations", handlers.GetVaccinations)
		api.GET("/vaccinations/upcoming", handlers.GetUpcomingVaccinations)
		api.GET("/vaccinations/:id", handlers.GetVaccinationByID)
		api.PUT("/vaccinations/:id", handlers.UpdateVaccination)
		api.DELETE("/vaccinations/:id", handlers.DeleteVaccination)

		// Produce routes
		api.POST("/produce", handlers.CreateProduce)
		api.GET("/produce", handlers.GetProduce)
		api.GET("/produce/summary", handlers.GetProduceSummary)

		// Death routes
		api.POST("/deaths", handlers.CreateDeath)
		api.GET("/deaths", handlers.GetDeaths)
		api.GET("/deaths/statistics", handlers.GetDeathStatistics)

			// ============ Lineage Routes ============
	api.GET("/animals/:id/ancestors", handlers.GetAncestors)
	api.GET("/animals/:id/descendants", handlers.GetDescendants)
	api.GET("/animals/:id/family-tree", handlers.GetFamilyTree)

	// ============ Feed Management Routes ============
	api.POST("/feed", handlers.CreateFeedRecord)
	api.GET("/feed", handlers.GetFeedInventory)
	api.GET("/feed/low-stock", handlers.GetLowStockItems)
	api.PUT("/feed/:id", handlers.UpdateFeedStock)
	}
			// ============ Farm Routes ============
	api.POST("/farms", handlers.CreateFarm)
	api.GET("/farms", handlers.GetFarms)
	api.GET("/farms/:id", handlers.GetFarmByID)
	api.PUT("/farms/:id", handlers.UpdateFarm)
	api.DELETE("/farms/:id", handlers.DeleteFarm)

	// Start server
	port := "8080"
	log.Printf("✅ Server running on http://localhost:%s", port)
	log.Fatal(r.Run(":" + port))
}