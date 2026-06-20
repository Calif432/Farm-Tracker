package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// CreateVaccination - POST /api/v1/vaccinations
func CreateVaccination(c *gin.Context) {
	var input struct {
		AnimalID          string `json:"animal_id" binding:"required"`
		VaccineName       string `json:"vaccine_name" binding:"required"`
		AdministeredDate  string `json:"administered_date" binding:"required"`
		NextDueDate       string `json:"next_due_date"`
		AdministeredBy    string `json:"administered_by"`
		BatchNumber       string `json:"batch_number"`
		Notes             string `json:"notes"`
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

	adminDate, err := time.Parse("2006-01-02", input.AdministeredDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid administered date format, use YYYY-MM-DD"})
		return
	}

	vaccination := models.Vaccination{
		AnimalID:          animalID,
		VaccineName:       input.VaccineName,
		AdministeredDate:  adminDate,
		AdministeredBy:    input.AdministeredBy,
		BatchNumber:       input.BatchNumber,
		Notes:             input.Notes,
	}

	if input.NextDueDate != "" {
		nextDate, err := time.Parse("2006-01-02", input.NextDueDate)
		if err == nil {
			vaccination.NextDueDate = &nextDate
		}
	}

	if err := database.DB.Create(&vaccination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vaccination record"})
		return
	}

	c.JSON(http.StatusCreated, vaccination)
}

// GetVaccinations - GET /api/v1/vaccinations
func GetVaccinations(c *gin.Context) {
	var vaccinations []models.Vaccination
	animalID := c.Query("animal_id")

	query := database.DB
	if animalID != "" {
		query = query.Where("animal_id = ?", animalID)
	}

	if err := query.Preload("Animal").Find(&vaccinations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vaccinations"})
		return
	}

	c.JSON(http.StatusOK, vaccinations)
}

// GetVaccinationByID - GET /api/v1/vaccinations/:id
func GetVaccinationByID(c *gin.Context) {
	id := c.Param("id")
	var vaccination models.Vaccination

	if err := database.DB.Preload("Animal").First(&vaccination, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vaccination record not found"})
		return
	}

	c.JSON(http.StatusOK, vaccination)
}

// UpdateVaccination - PUT /api/v1/vaccinations/:id
func UpdateVaccination(c *gin.Context) {
	id := c.Param("id")
	var vaccination models.Vaccination

	if err := database.DB.First(&vaccination, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vaccination record not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&vaccination).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vaccination"})
		return
	}

	c.JSON(http.StatusOK, vaccination)
}

// DeleteVaccination - DELETE /api/v1/vaccinations/:id
func DeleteVaccination(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.Vaccination{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vaccination"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vaccination record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vaccination record deleted successfully"})
}

// GetUpcomingVaccinations - GET /api/v1/vaccinations/upcoming
func GetUpcomingVaccinations(c *gin.Context) {
	var vaccinations []models.Vaccination
	today := time.Now()
	nextWeek := today.AddDate(0, 0, 7)

	if err := database.DB.
		Where("next_due_date BETWEEN ? AND ?", today, nextWeek).
		Preload("Animal").
		Find(&vaccinations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upcoming vaccinations"})
		return
	}

	c.JSON(http.StatusOK, vaccinations)
}