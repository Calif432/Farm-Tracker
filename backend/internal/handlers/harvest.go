package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// CreateHarvest - POST /api/v1/harvests
func CreateHarvest(c *gin.Context) {
	var input struct {
		PlantingID   string  `json:"planting_id" binding:"required"`
		HarvestDate  string  `json:"harvest_date" binding:"required"`
		Quantity     float64 `json:"quantity" binding:"required"`
		Unit         string  `json:"unit" binding:"required"`
		QualityGrade string  `json:"quality_grade"`
		Notes        string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plantingID, err := uuid.Parse(input.PlantingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	harvestDate, err := time.Parse("2006-01-02", input.HarvestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid harvest date format, use YYYY-MM-DD"})
		return
	}

	harvest := models.Harvest{
		PlantingID:   plantingID,
		HarvestDate:  harvestDate,
		Quantity:     input.Quantity,
		Unit:         input.Unit,
		QualityGrade: input.QualityGrade,
		Notes:        input.Notes,
	}

	if err := database.DB.Create(&harvest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create harvest record: " + err.Error()})
		return
	}

	// Update planting status to harvested
	database.DB.Model(&models.Planting{}).Where("id = ?", plantingID).Updates(map[string]interface{}{
		"status":           "harvested",
		"actual_harvest_date": harvestDate,
	})

	c.JSON(http.StatusCreated, harvest)
}

// GetHarvests - GET /api/v1/harvests
func GetHarvests(c *gin.Context) {
	var harvests []models.Harvest
	plantingID := c.Query("planting_id")

	query := database.DB
	if plantingID != "" {
		query = query.Where("planting_id = ?", plantingID)
	}

	if err := query.Preload("Planting.Field").Order("harvest_date DESC").Find(&harvests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch harvests"})
		return
	}

	c.JSON(http.StatusOK, harvests)
}

// GetHarvestByID - GET /api/v1/harvests/:id
func GetHarvestByID(c *gin.Context) {
	id := c.Param("id")
	var harvest models.Harvest

	if err := database.DB.Preload("Planting.Field").First(&harvest, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Harvest not found"})
		return
	}

	c.JSON(http.StatusOK, harvest)
}