package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Calif432/farmtrack/backend/internal/database"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

// GetAncestors - GET /api/v1/animals/:id/ancestors
// Returns all ancestors (parents, grandparents, etc.)
func GetAncestors(c *gin.Context) {
	id := c.Param("id")
	animalID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	var ancestors []models.Animal
	var currentID = &animalID

	// Recursively find parents up to 5 generations
	for i := 0; i < 5 && currentID != nil; i++ {
		var animal models.Animal
		if err := database.DB.First(&animal, "id = ?", currentID).Error; err != nil {
			break
		}

		// Add parents if they exist
		if animal.MotherID != nil {
			var mother models.Animal
			if err := database.DB.First(&mother, "id = ?", animal.MotherID).Error; err == nil {
				ancestors = append(ancestors, mother)
			}
		}
		if animal.FatherID != nil {
			var father models.Animal
			if err := database.DB.First(&father, "id = ?", animal.FatherID).Error; err == nil {
				ancestors = append(ancestors, father)
			}
		}

		// Move up to mother for next iteration
		if animal.MotherID != nil {
			currentID = animal.MotherID
		} else if animal.FatherID != nil {
			currentID = animal.FatherID
		} else {
			currentID = nil
		}
	}

	c.JSON(http.StatusOK, ancestors)
}

// GetDescendants - GET /api/v1/animals/:id/descendants
// Returns all children, grandchildren, etc.
func GetDescendants(c *gin.Context) {
	id := c.Param("id")
	animalID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	var descendants []models.Animal
	var children []models.Animal

	// Find direct children
	if err := database.DB.Where("mother_id = ? OR father_id = ?", animalID, animalID).Find(&children).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch descendants"})
		return
	}

	descendants = append(descendants, children...)

	// Find grandchildren (children of children)
	for _, child := range children {
		var grandchildren []models.Animal
		if err := database.DB.Where("mother_id = ? OR father_id = ?", child.ID, child.ID).Find(&grandchildren).Error; err == nil {
			descendants = append(descendants, grandchildren...)
		}
	}

	c.JSON(http.StatusOK, descendants)
}

// GetFamilyTree - GET /api/v1/animals/:id/family-tree
// Returns a structured family tree with relationships
func GetFamilyTree(c *gin.Context) {
	id := c.Param("id")
	animalID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	var animal models.Animal
	if err := database.DB.First(&animal, "id = ?", animalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
		return
	}

	// Build the family tree structure
	type FamilyMember struct {
		Animal   models.Animal    `json:"animal"`
		Children []FamilyMember   `json:"children,omitempty"`
		Parents  []models.Animal  `json:"parents,omitempty"`
	}

	// Get parents
	var parents []models.Animal
	if animal.MotherID != nil {
		var mother models.Animal
		if err := database.DB.First(&mother, "id = ?", animal.MotherID).Error; err == nil {
			parents = append(parents, mother)
		}
	}
	if animal.FatherID != nil {
		var father models.Animal
		if err := database.DB.First(&father, "id = ?", animal.FatherID).Error; err == nil {
			parents = append(parents, father)
		}
	}

	// Get children
	var children []models.Animal
	database.DB.Where("mother_id = ? OR father_id = ?", animalID, animalID).Find(&children)

	// Build response
	response := FamilyMember{
		Animal:  animal,
		Parents: parents,
	}

	// Add grandchildren (optional - can be expanded)
	for _, child := range children {
		var grandchildren []models.Animal
		database.DB.Where("mother_id = ? OR father_id = ?", child.ID, child.ID).Find(&grandchildren)
		response.Children = append(response.Children, FamilyMember{
			Animal: child,
		})
		// Add grandchildren if needed
		if len(grandchildren) > 0 {
			for _, grandchild := range grandchildren {
				response.Children = append(response.Children, FamilyMember{
					Animal: grandchild,
				})
			}
		}
	}

	c.JSON(http.StatusOK, response)
}