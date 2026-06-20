package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// CreateFarm - POST /api/v1/farms
func CreateFarm(c *gin.Context) {
	var input struct {
		Name         string  `json:"name" binding:"required"`
		Location     string  `json:"location"`
		SizeHectares float64 `json:"size_hectares"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	farm := models.Farm{
		Name:         input.Name,
		Location:     input.Location,
		SizeHectares: input.SizeHectares,
		OwnerID:      userID.(uuid.UUID),
	}

	if err := database.DB.Create(&farm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create farm"})
		return
	}

	c.JSON(http.StatusCreated, farm)
}

// GetFarms - GET /api/v1/farms
func GetFarms(c *gin.Context) {
	var farms []models.Farm

	if err := database.DB.Find(&farms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch farms"})
		return
	}

	c.JSON(http.StatusOK, farms)
}

// GetFarmByID - GET /api/v1/farms/:id
func GetFarmByID(c *gin.Context) {
	id := c.Param("id")
	var farm models.Farm

	if err := database.DB.First(&farm, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Farm not found"})
		return
	}

	c.JSON(http.StatusOK, farm)
}

// UpdateFarm - PUT /api/v1/farms/:id
func UpdateFarm(c *gin.Context) {
	id := c.Param("id")
	var farm models.Farm

	if err := database.DB.First(&farm, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Farm not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&farm).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update farm"})
		return
	}

	c.JSON(http.StatusOK, farm)
}

// DeleteFarm - DELETE /api/v1/farms/:id
func DeleteFarm(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.Farm{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete farm"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Farm not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Farm deleted successfully"})
}