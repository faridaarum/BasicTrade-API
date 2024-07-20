package controllers

import (
	"BasicTradeApp/config"
	"BasicTradeApp/models"
	"BasicTradeApp/services"
	"BasicTradeApp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register registers a new user
func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	admin := models.Admin{Email: input.Email, Password: hashedPassword}
	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User already exist "})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// Login logs in a user and returns JWT tokens
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
	if err := config.DB.Where("email = ?", input.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Email"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, admin.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	tokenDetails, err := services.GenerateJWT(admin.ID, admin.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	})
}

// RefreshToken refreshes the JWT tokens
func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	td, err := services.RefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": td.AccessToken, "refresh_token": td.RefreshToken})
}
