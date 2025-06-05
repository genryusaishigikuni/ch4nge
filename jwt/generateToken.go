package jwt

import (
	"github.com/genryusaishigikuni/ch4nge/config"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func GenerateToken(userID uint, email string, isAdmin bool) (string, error) {
	claims := &models.Claims{
		UserID:  userID,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ch4nge-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
