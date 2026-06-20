package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// GetDashboardStats - GET /api/v1/dashboard/stats
func GetDashboardStats(c *gin.Context) {
	// Get farm ID from query or use default
	farmID := c.Query("farm_id")

	stats := gin.H{}

	// 1. Animal Statistics
	var totalAnimals int64
	var animalsByType []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}

	query := database.DB.Model(&models.Animal{})
	if farmID != "" {
		query = query.Where("farm_id = ?", farmID)
	}
	query.Count(&totalAnimals)

	database.DB.Model(&models.Animal{}).
		Select("type, COUNT(*) as count").
		Where("status = ?", "alive").
		Group("type").
		Scan(&animalsByType)

	stats["total_animals"] = totalAnimals
	stats["animals_by_type"] = animalsByType

	// 2. Recent Produce (last 7 days)
	var recentProduce []models.AnimalProduce
	database.DB.Where("recorded_at >= ?", time.Now().AddDate(0, 0, -7)).
		Preload("Animal").
		Order("recorded_at DESC").
		Limit(10).
		Find(&recentProduce)
	stats["recent_produce"] = recentProduce

	// 3. Upcoming Vaccinations (next 7 days)
	var upcomingVaccinations []models.Vaccination
	database.DB.Where("next_due_date BETWEEN ? AND ?", time.Now(), time.Now().AddDate(0, 0, 7)).
		Preload("Animal").
		Find(&upcomingVaccinations)
	stats["upcoming_vaccinations"] = upcomingVaccinations

	// 4. Crop Statistics
	var totalFields int64
	var totalPlantings int64
	var plantingsByStatus []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	database.DB.Model(&models.Field{}).Count(&totalFields)
	database.DB.Model(&models.Planting{}).Count(&totalPlantings)
	database.DB.Model(&models.Planting{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&plantingsByStatus)

	stats["total_fields"] = totalFields
	stats["total_plantings"] = totalPlantings
	stats["plantings_by_status"] = plantingsByStatus

	// 5. Inventory Summary
	var totalInventoryItems int64
	var lowStockItems int64

	database.DB.Model(&models.InventoryItem{}).Count(&totalInventoryItems)
	database.DB.Model(&models.InventoryItem{}).Where("quantity <= min_stock_level").Count(&lowStockItems)

	stats["total_inventory_items"] = totalInventoryItems
	stats["low_stock_items"] = lowStockItems

	// 6. Recent Harvests
	var recentHarvests []models.Harvest
	database.DB.Where("harvest_date >= ?", time.Now().AddDate(0, -1, 0)).
		Preload("Planting.Field").
		Order("harvest_date DESC").
		Limit(10).
		Find(&recentHarvests)
	stats["recent_harvests"] = recentHarvests

	c.JSON(http.StatusOK, stats)
}

// GetAnimalStatistics - GET /api/v1/dashboard/animals
func GetAnimalStatistics(c *gin.Context) {
	stats := gin.H{}

	// Total animals
	var total int64
	database.DB.Model(&models.Animal{}).Count(&total)
	stats["total"] = total

	// By status
	var byStatus []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	database.DB.Model(&models.Animal{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&byStatus)
	stats["by_status"] = byStatus

	// By type
	var byType []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	database.DB.Model(&models.Animal{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&byType)
	stats["by_type"] = byType

	// Recent deaths (last 30 days)
	var recentDeaths []models.Death
	database.DB.Where("death_date >= ?", time.Now().AddDate(0, -1, 0)).
		Preload("Animal").
		Order("death_date DESC").
		Limit(10).
		Find(&recentDeaths)
	stats["recent_deaths"] = recentDeaths

	c.JSON(http.StatusOK, stats)
}

// GetCropStatistics - GET /api/v1/dashboard/crops
func GetCropStatistics(c *gin.Context) {
	stats := gin.H{}

	// Total harvests
	var totalHarvests int64
	database.DB.Model(&models.Harvest{}).Count(&totalHarvests)
	stats["total_harvests"] = totalHarvests

	// Total harvest quantity by crop type
	var harvestByCrop []struct {
		CropType string  `json:"crop_type"`
		Total    float64 `json:"total"`
		Unit     string  `json:"unit"`
	}
	database.DB.Table("harvests").
		Select("plantings.crop_type, SUM(harvests.quantity) as total, harvests.unit").
		Joins("JOIN plantings ON harvests.planting_id = plantings.id").
		Group("plantings.crop_type, harvests.unit").
		Scan(&harvestByCrop)
	stats["harvest_by_crop"] = harvestByCrop

	// Recent harvests
	var recentHarvests []models.Harvest
	database.DB.Where("harvest_date >= ?", time.Now().AddDate(0, -1, 0)).
		Preload("Planting.Field").
		Order("harvest_date DESC").
		Limit(10).
		Find(&recentHarvests)
	stats["recent_harvests"] = recentHarvests

	c.JSON(http.StatusOK, stats)
}