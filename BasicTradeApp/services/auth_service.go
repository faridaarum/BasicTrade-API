package services

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var jwtSecret = []byte("your_secret_key")
var refreshSecret = []byte("your_refresh_secret_key")

type Claims struct {
	AdminID uint   `json:"admin_id"`
	Email   string `json:"email"`
	jwt.StandardClaims
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func GenerateJWT(adminID uint, email string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(15 * time.Minute).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(7 * 24 * time.Hour).Unix()
	td.RefreshUuid = uuid.New().String()

	atClaims := &Claims{
		AdminID: adminID,
		Email:   email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: td.AtExpires,
		},
	}

	rtClaims := &Claims{
		AdminID: adminID,
		Email:   email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: td.RtExpires,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, _ = accessToken.SignedString(jwtSecret)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, _ = refreshToken.SignedString(refreshSecret)

	return td, nil
}

func RefreshToken(refreshTokenString string) (*TokenDetails, error) {
	claims := &Claims{}

	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !refreshToken.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return GenerateJWT(claims.AdminID, claims.Email)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
