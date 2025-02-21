package controllers

import (
	"net/http"

	"github.com/NurymGM/jwtnew/initializers"
	"github.com/NurymGM/jwtnew/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	if err := initializers.DB.Where("email = ?", input.Email).First(&user).Error; err == nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "Email already exists"})
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
	
}

func Validate(c *gin.Context) {

}
