package main // CompileDaemon -command="./jwtnew"

import (
	"github.com/NurymGM/jwtnew/controllers"
	"github.com/NurymGM/jwtnew/initializers"
	"github.com/NurymGM/jwtnew/middleware"
	"github.com/NurymGM/jwtnew/migrations"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	migrations.Migrate()
}

func main() {
	r := gin.Default()
	r.GET("/", controllers.RootRoute)
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.LogIn)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.Run()
}
