package controllers

import (
	"BasicTradeApp/config"
	"BasicTradeApp/models"
	"BasicTradeApp/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateVariant(c *gin.Context) {
	var input struct {
		VariantName string `json:"variant_name" binding:"required"`
		ProductID   uint   `json:"product_id" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required"`
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

	var existingVariant models.Variant
	if err := config.DB.Where("product_id = ? AND variant_name = ?", input.ProductID, input.VariantName).First(&existingVariant).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product already exists"})
		return
	}

	adminID, err := utils.ExtractAdminID(c)
	if err != nil || product.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	variant := models.Variant{VariantName: input.VariantName, ProductID: input.ProductID, Quantity: input.Quantity}
	if err := config.DB.Create(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variant": variant})
}

func GetVariants(c *gin.Context) {
	productID := c.Param("product_id")

	if _, err := strconv.Atoi(productID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	var variants []models.Variant
	if err := config.DB.Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching variants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variants": variants})
}

func GetAllVariants(c *gin.Context) {
	var variants []models.Variant
	variantName := c.Query("variant_name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	offset := (page - 1) * pageSize

	query := config.DB.Limit(pageSize).Offset(offset)
	if variantName != "" {
		query = query.Where("variant_name LIKE ?", "%"+variantName+"%")
	}

	if err := query.Find(&variants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching variants"})
		return
	}
	if variantName != "" && len(variants) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Variant not exist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variants": variants})
}

func GetVariantByID(c *gin.Context) {
	var variant models.Variant
	variantID := c.Param("variant_id")

	if _, err := strconv.Atoi(variantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := config.DB.First(&variant, "id = ?", variantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variant": variant})
}

func UpdateVariant(c *gin.Context) {
	var input struct {
		VariantName string `json:"variant_name"`
		Quantity    int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var variant models.Variant
	variantID := c.Param("variant_id")

	if _, err := strconv.Atoi(variantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := config.DB.First(&variant, "id = ?", variantID).Error; err != nil {
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

	if input.VariantName != "" {
		variant.VariantName = input.VariantName
	}
	if input.Quantity != 0 {
		variant.Quantity = input.Quantity
	}

	if err := config.DB.Save(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variant": variant})
}

func DeleteVariant(c *gin.Context) {
	var variant models.Variant
	variantID := c.Param("variant_id")

	if _, err := strconv.Atoi(variantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := config.DB.First(&variant, "id = ?", variantID).Error; err != nil {
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
