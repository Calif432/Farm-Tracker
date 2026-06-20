package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// CreateFeedRecord - POST /api/v1/feed
func CreateFeedRecord(c *gin.Context) {
	var input struct {
		Name          string  `json:"name" binding:"required"`
		Category      string  `json:"category" binding:"required"`
		Quantity      float64 `json:"quantity" binding:"required"`
		Unit          string  `json:"unit" binding:"required"`
		CostPerUnit   float64 `json:"cost_per_unit"`
		Supplier      string  `json:"supplier"`
		MinStockLevel float64 `json:"min_stock_level"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := models.InventoryItem{
		Name:          input.Name,
		Category:      input.Category,
		Quantity:      input.Quantity,
		Unit:          input.Unit,
		CostPerUnit:   input.CostPerUnit,
		Supplier:      input.Supplier,
		MinStockLevel: input.MinStockLevel,
		LastUpdated:   time.Now(),
		// FarmID can be nil (optional)
	}

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feed record: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}
// GetFeedInventory - GET /api/v1/feed
func GetFeedInventory(c *gin.Context) {
	var items []models.InventoryItem
	category := c.Query("category")

	query := database.DB
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Order("name").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feed inventory"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// UpdateFeedStock - PUT /api/v1/feed/:id
func UpdateFeedStock(c *gin.Context) {
	id := c.Param("id")
	var item models.InventoryItem

	if err := database.DB.First(&item, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feed item not found"})
		return
	}

	var input struct {
		Quantity float64 `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.Quantity = input.Quantity
	item.LastUpdated = time.Now()

	if err := database.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetLowStockItems - GET /api/v1/feed/low-stock
func GetLowStockItems(c *gin.Context) {
	var items []models.InventoryItem

	if err := database.DB.Where("quantity <= min_stock_level").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch low stock items"})
		return
	}

	c.JSON(http.StatusOK, items)
}