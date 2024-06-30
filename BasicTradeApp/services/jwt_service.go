package services

import (
	"BasicTradeApp/config"
	"BasicTradeApp/models"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractToken(bearerToken string) (string, error) {
	tokenString := strings.Split(bearerToken, "Bearer ")
	if len(tokenString) != 2 {
		return "", errors.New("invalid token format")
	}
	return tokenString[1], nil
}

func ExtractAdminID(c *gin.Context) (uint, error) {
	tokenString, err := ExtractToken(c.GetHeader("Authorization"))
	if err != nil {
		return 0, err
	}

	claims, err := ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	var admin models.Admin
	if err := config.DB.Where("username = ?", claims.Username).First(&admin).Error; err != nil {
		return 0, errors.New("user not found")
	}

	return admin.ID, nil
}
