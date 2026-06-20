package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// CreatePlanting - POST /api/v1/plantings
func CreatePlanting(c *gin.Context) {
	var input struct {
		FieldID            string  `json:"field_id" binding:"required"`
		CropType           string  `json:"crop_type" binding:"required"`
		Variety            string  `json:"variety"`
		PlantingDate       string  `json:"planting_date" binding:"required"`
		ExpectedHarvestDate string  `json:"expected_harvest_date"`
		QuantityPlanted    float64 `json:"quantity_planted"`
		Unit               string  `json:"unit"`
		Notes              string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fieldID, err := uuid.Parse(input.FieldID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field ID"})
		return
	}

	plantingDate, err := time.Parse("2006-01-02", input.PlantingDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting date format, use YYYY-MM-DD"})
		return
	}

	planting := models.Planting{
		FieldID:         fieldID,
		CropType:        input.CropType,
		Variety:         input.Variety,
		PlantingDate:    plantingDate,
		QuantityPlanted: input.QuantityPlanted,
		Unit:            input.Unit,
		Status:          "growing",
		Notes:           input.Notes,
	}

	if input.ExpectedHarvestDate != "" {
		expectedDate, err := time.Parse("2006-01-02", input.ExpectedHarvestDate)
		if err == nil {
			planting.ExpectedHarvestDate = &expectedDate
		}
	}

	if err := database.DB.Create(&planting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create planting record: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, planting)
}

// GetPlantings - GET /api/v1/plantings
func GetPlantings(c *gin.Context) {
	var plantings []models.Planting
	fieldID := c.Query("field_id")
	cropType := c.Query("crop_type")
	status := c.Query("status")

	query := database.DB
	if fieldID != "" {
		query = query.Where("field_id = ?", fieldID)
	}
	if cropType != "" {
		query = query.Where("crop_type = ?", cropType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Preload("Field").Order("planting_date DESC").Find(&plantings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plantings"})
		return
	}

	c.JSON(http.StatusOK, plantings)
}

// GetPlantingByID - GET /api/v1/plantings/:id
func GetPlantingByID(c *gin.Context) {
	id := c.Param("id")
	var planting models.Planting

	if err := database.DB.Preload("Field").Preload("Harvests").First(&planting, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting not found"})
		return
	}

	c.JSON(http.StatusOK, planting)
}