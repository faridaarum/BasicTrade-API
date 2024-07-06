package controllers

import (
	"BasicTradeApp/config"
	"BasicTradeApp/models"
	"BasicTradeApp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateVariant(c *gin.Context) {
	var input struct {
		Name      string    `json:"name" binding:"required"`
		ProductID uuid.UUID `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, "id = ?", input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	adminID, err := utils.ExtractAdminID(c)
	if err != nil || product.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	variant := models.Variant{Name: input.Name, ProductID: input.ProductID}
	if err := config.DB.Create(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variant": variant})
}

func GetVariants(c *gin.Context) {
	var variants []models.Variant
	if err := config.DB.Find(&variants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching variants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variants": variants})
}

func GetVariantByID(c *gin.Context) {
	var variant models.Variant
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := config.DB.First(&variant, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variant": variant})
}

func UpdateVariant(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var variant models.Variant
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := config.DB.First(&variant, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, "id = ?", variant.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	adminID, err := utils.ExtractAdminID(c)
	if err != nil || product.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if input.Name != "" {
		variant.Name = input.Name
	}

	if err := config.DB.Save(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variant": variant})
}

func DeleteVariant(c *gin.Context) {
	var variant models.Variant
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := config.DB.First(&variant, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, "id = ?", variant.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	adminID, err := utils.ExtractAdminID(c)
	if err != nil || product.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := config.DB.Delete(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Variant deleted"})
}
