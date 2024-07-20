package controllers

import (
	"BasicTradeApp/config"
	"BasicTradeApp/models"
	"BasicTradeApp/utils"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	// Parse name from form-data
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Check if the product already exists
	var existingProduct models.Product
	if err := config.DB.Where("name = ?", name).First(&existingProduct).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product already exist"})
		return
	}

	// Parse file from form-data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// Save the uploaded file to a temporary directory
	filePath := "./" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer os.Remove(filePath) // Clean up the file after upload

	// Initialize Cloudinary
	cld, err := cloudinary.NewFromParams("drkuqdska", "129773322849332", "PY8ecDt2ChkHyfyla-Aq3Gvqiqg")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary initialization failed"})
		return
	}

	// Upload the file to Cloudinary
	uploadResult, err := cld.Upload.Upload(c, filePath, uploader.UploadParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to Cloudinary"})
		return
	}

	// Extract admin ID from token
	adminID, err := utils.ExtractAdminID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create product with the uploaded image URL
	product := models.Product{
		Name:     name,
		ImageURL: uploadResult.SecureURL,
		AdminID:  adminID,
	}
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func GetProducts(c *gin.Context) {
	var products []models.Product
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	offset := (page - 1) * pageSize

	query := config.DB.Preload("Variants").Limit(pageSize).Offset(offset)
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
		return
	}
	if name != "" && len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not exist"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	var product models.Product
	if err := config.DB.Preload("Variants").First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func UpdateProduct(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	adminID, err := utils.ExtractAdminID(c)
	if err != nil || product.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if the product is already edited with the same values
	if input.Name == product.Name && input.ImageURL == product.ImageURL {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No changes"})
		return
	}

	// Update product fields if they are not empty and different from current values
	if input.Name != "" && input.Name != product.Name {
		product.Name = input.Name
	}
	if input.ImageURL != "" && input.ImageURL != product.ImageURL {
		product.ImageURL = input.ImageURL
	}

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func DeleteProduct(c *gin.Context) {
	var product models.Product
	if err := config.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	adminID, err := utils.ExtractAdminID(c)
	if err != nil || product.AdminID != adminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
