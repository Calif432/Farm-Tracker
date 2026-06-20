package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// CreateDeath - POST /api/v1/deaths
func CreateDeath(c *gin.Context) {
	var input struct {
		AnimalID       string `json:"animal_id" binding:"required"`
		DeathDate      string `json:"death_date" binding:"required"`
		Reason         string `json:"reason" binding:"required"`
		PreventionNotes string `json:"prevention_notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	animalID, err := uuid.Parse(input.AnimalID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	// Check if animal exists
	var animal models.Animal
	if err := database.DB.First(&animal, "id = ?", animalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
		return
	}

	// Update animal status to dead
	animal.Status = "dead"
	if err := database.DB.Save(&animal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update animal status"})
		return
	}

	deathDate, err := time.Parse("2006-01-02", input.DeathDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid death date format, use YYYY-MM-DD"})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")
	userUUID := userID.(uuid.UUID)

	death := models.Death{
		AnimalID:       animalID,
		DeathDate:      deathDate,
		Reason:         input.Reason,
		PreventionNotes: input.PreventionNotes,
		RecordedBy:     &userUUID,
	}

	if err := database.DB.Create(&death).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create death record"})
		return
	}

	c.JSON(http.StatusCreated, death)
}

// GetDeaths - GET /api/v1/deaths
func GetDeaths(c *gin.Context) {
	var deaths []models.Death
	animalID := c.Query("animal_id")

	query := database.DB
	if animalID != "" {
		query = query.Where("animal_id = ?", animalID)
	}

	if err := query.Preload("Animal").Preload("RecordedByUser").Order("death_date DESC").Find(&deaths).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch death records"})
		return
	}

	c.JSON(http.StatusOK, deaths)
}

// GetDeathStatistics - GET /api/v1/deaths/statistics
func GetDeathStatistics(c *gin.Context) {
	var stats []struct {
		Reason string `json:"reason"`
		Count  int64  `json:"count"`
	}

	database.DB.Table("deaths").
		Select("reason, COUNT(*) as count").
		Group("reason").
		Order("count DESC").
		Scan(&stats)

	c.JSON(http.StatusOK, stats)
}