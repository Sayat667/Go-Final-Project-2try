package routes

import (
	// "fmt"
	"gofp/controllers"
	// "net/http"

	"github.com/gin-gonic/gin"
	"log"
	"gorm.io/gorm"
)


var DB *gorm.DB

func AuthRoutes(r *gin.Engine) {
	router := r.Group("/")

	router.GET("/", controllers.HomePage)
	router.POST("/signup", controllers.SignUp)
	router.POST("/signin", controllers.SignIn)
	// router.GET("/api", controllers.JWTMiddleware(), controllers.APIHandler)

	log.Println("Server starting on http://localhost:8080/")
	log.Fatal(r.Run(":8080"))
}
