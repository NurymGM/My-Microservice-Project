package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/NurymGM/jwtnew/initializers"
	"github.com/NurymGM/jwtnew/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RootRoute(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Hello wws!"})
}

func SignUp(c *gin.Context) {
	input := models.Userr{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Incorrect Inputs"})
		return
	}
	user := models.Userr{}
	err := initializers.DB.Where("email = ?", input.Email).First(&user).Error
	if err == nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Email already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	validRoles := map[string]bool{"owner": true, "guest": true}
	if !validRoles[input.Role] {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid role, must be 'owner' or 'guest'"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Couldnt Hash Password"})
		return
	}
	input.Password = string(hashed)
	result := initializers.DB.Create(&input) // before creating, email, password formats can be checked
	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Couldnt Create User", "error": result.Error.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func LogIn(c *gin.Context) {
	input := models.Userr{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Incorrect Inputs"})
		return
	}
	user := models.Userr{}
	err := initializers.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User doesn't exist"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}
	// generating jwt token for 1 month
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// signing jwt token using secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY"))) // dont forget the []byte here
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create token"})
		return
	}
	// sending it back as a cookie
	c.SetSameSite(http.SameSiteLaxMode) // this is for safety
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", true, true)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Logged in!"})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.IndentedJSON(http.StatusOK, gin.H{"message (from validate)": user})
}
