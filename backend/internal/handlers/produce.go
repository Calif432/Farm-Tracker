package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// CreateProduce - POST /api/v1/produce
func CreateProduce(c *gin.Context) {
	var input struct {
		AnimalID    string  `json:"animal_id" binding:"required"`
		ProduceType string  `json:"produce_type" binding:"required"`
		Quantity    float64 `json:"quantity" binding:"required"`
		Unit        string  `json:"unit" binding:"required"`
		RecordedAt  string  `json:"recorded_at"`
		Notes       string  `json:"notes"`
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

	produce := models.AnimalProduce{
		AnimalID:    animalID,
		ProduceType: input.ProduceType,
		Quantity:    input.Quantity,
		Unit:        input.Unit,
		Notes:       input.Notes,
	}

	if input.RecordedAt != "" {
		recordedAt, err := time.Parse("2006-01-02", input.RecordedAt)
		if err == nil {
			produce.RecordedAt = recordedAt
		}
	}

	if err := database.DB.Create(&produce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create produce record"})
		return
	}

	c.JSON(http.StatusCreated, produce)
}

// GetProduce - GET /api/v1/produce
func GetProduce(c *gin.Context) {
	var produce []models.AnimalProduce
	animalID := c.Query("animal_id")
	produceType := c.Query("produce_type")

	query := database.DB
	if animalID != "" {
		query = query.Where("animal_id = ?", animalID)
	}
	if produceType != "" {
		query = query.Where("produce_type = ?", produceType)
	}

	if err := query.Preload("Animal").Order("recorded_at DESC").Find(&produce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch produce records"})
		return
	}

	c.JSON(http.StatusOK, produce)
}

// GetProduceSummary - GET /api/v1/produce/summary
func GetProduceSummary(c *gin.Context) {
	animalID := c.Query("animal_id")
	produceType := c.Query("produce_type")
	period := c.DefaultQuery("period", "month") // day, week, month, year

	var results []struct {
		ProduceType string  `json:"produce_type"`
		Total       float64 `json:"total"`
		Unit        string  `json:"unit"`
	}

	query := database.DB.Table("animal_produces")
	if animalID != "" {
		query = query.Where("animal_id = ?", animalID)
	}
	if produceType != "" {
		query = query.Where("produce_type = ?", produceType)
	}

	// Filter by period
	var dateFilter time.Time
	switch period {
	case "day":
		dateFilter = time.Now().AddDate(0, 0, -1)
	case "week":
		dateFilter = time.Now().AddDate(0, 0, -7)
	case "month":
		dateFilter = time.Now().AddDate(0, -1, 0)
	case "year":
		dateFilter = time.Now().AddDate(-1, 0, 0)
	default:
		dateFilter = time.Now().AddDate(0, -1, 0)
	}

	query = query.Where("recorded_at >= ?", dateFilter)

	query.Select("produce_type, SUM(quantity) as total, unit").
		Group("produce_type, unit").
		Scan(&results)

	c.JSON(http.StatusOK, results)
}