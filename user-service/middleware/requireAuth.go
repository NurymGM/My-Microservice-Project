package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/NurymGM/jwtnew/initializers"
	"github.com/NurymGM/jwtnew/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func RequireAuth(c *gin.Context) {
	// get cookie from request body
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing auth token"})
		return
	}
	// decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}
	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
		return
	}
	// Validate expiration
	exp, ok := claims["exp"].(float64)
	if !ok || float64(time.Now().Unix()) > exp {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token expired"})
		return
	}
	// Extract user ID safely
	sub, ok := claims["sub"].(float64)  // JWT numbers are stored as float64
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token subject"})
		return
	}
	userID := uint(sub)
	// Find the user in the database
	user := models.Userr{}
	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User doesn't exist"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}
	// Attach user to the request
	c.Set("user", user)
	// Continue the request
	c.Next()
}
