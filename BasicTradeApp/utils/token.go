package utils

import (
	"BasicTradeApp/config"
	"BasicTradeApp/models"
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("your_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

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

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	var admin models.Admin
	if err := config.DB.Where("email = ?", claims.Email).First(&admin).Error; err != nil {
		return 0, errors.New("user not found")
	}

	return admin.ID, nil
}
