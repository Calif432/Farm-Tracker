package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// CreateAnimal - POST /api/v1/animals
func CreateAnimal(c *gin.Context) {
	var input struct {
		TagID       string `json:"tag_id" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Breed       string `json:"breed"`
		Gender      string `json:"gender" binding:"required"`
		DateOfBirth string `json:"date_of_birth" binding:"required"`
		MotherID    string `json:"mother_id"`
		FatherID    string `json:"father_id"`
		Notes       string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dob, err := time.Parse("2006-01-02", input.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	animal := models.Animal{
		TagID:       input.TagID,
		Type:        input.Type,
		Breed:       input.Breed,
		Gender:      input.Gender,
		DateOfBirth: dob,
		Status:      "alive",
		Notes:       input.Notes,
	}

	if input.MotherID != "" {
		if motherUUID, err := uuid.Parse(input.MotherID); err == nil {
			animal.MotherID = &motherUUID
		}
	}
	if input.FatherID != "" {
		if fatherUUID, err := uuid.Parse(input.FatherID); err == nil {
			animal.FatherID = &fatherUUID
		}
	}

	result := database.DB.Create(&animal)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, animal)
}

// GetAnimals - GET /api/v1/animals
func GetAnimals(c *gin.Context) {
	var animals []models.Animal
	result := database.DB.Find(&animals)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, animals)
}

// GetAnimalByID - GET /api/v1/animals/:id
func GetAnimalByID(c *gin.Context) {
	id := c.Param("id")
	var animal models.Animal

	result := database.DB.First(&animal, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
		return
	}
	c.JSON(http.StatusOK, animal)
}

// UpdateAnimal - PUT /api/v1/animals/:id
func UpdateAnimal(c *gin.Context) {
	id := c.Param("id")
	var animal models.Animal

	if result := database.DB.First(&animal, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Model(&animal).Updates(input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, animal)
}

// DeleteAnimal - DELETE /api/v1/animals/:id
func DeleteAnimal(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.Animal{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Animal deleted successfully"})
}