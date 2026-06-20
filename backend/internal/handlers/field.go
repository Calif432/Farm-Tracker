package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// CreateField - POST /api/v1/fields
func CreateField(c *gin.Context) {
	var input struct {
		Name         string  `json:"name" binding:"required"`
		SizeHectares float64 `json:"size_hectares" binding:"required"`
		SoilType     string  `json:"soil_type"`
		Location     string  `json:"location"`
		FarmID       string  `json:"farm_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	field := models.Field{
		Name:         input.Name,
		SizeHectares: input.SizeHectares,
		SoilType:     input.SoilType,
		Location:     input.Location,
	}

	if input.FarmID != "" {
		farmID, err := uuid.Parse(input.FarmID)
		if err == nil {
			field.FarmID = farmID
		}
	}

	if err := database.DB.Create(&field).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create field: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, field)
}

// GetFields - GET /api/v1/fields
func GetFields(c *gin.Context) {
	var fields []models.Field
	farmID := c.Query("farm_id")

	query := database.DB
	if farmID != "" {
		query = query.Where("farm_id = ?", farmID)
	}

	if err := query.Preload("Plantings").Find(&fields).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fields"})
		return
	}

	c.JSON(http.StatusOK, fields)
}

// GetFieldByID - GET /api/v1/fields/:id
func GetFieldByID(c *gin.Context) {
	id := c.Param("id")
	var field models.Field

	if err := database.DB.Preload("Plantings").First(&field, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Field not found"})
		return
	}

	c.JSON(http.StatusOK, field)
}

// UpdateField - PUT /api/v1/fields/:id
func UpdateField(c *gin.Context) {
	id := c.Param("id")
	var field models.Field

	if err := database.DB.First(&field, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Field not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&field).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update field"})
		return
	}

	c.JSON(http.StatusOK, field)
}

// DeleteField - DELETE /api/v1/fields/:id
func DeleteField(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.Field{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete field"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Field not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Field deleted successfully"})
}