package routes

import (
	// "fmt"
	"gofp/controllers"
	// "net/http"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

func AuthRoutes(r *gin.Engine) {
	r.GET("/", controllers.AllUsersHandler)

	// Создание пользователя
	r.POST("/create-user", controllers.CreateUserHandler)

	// Обновление пользователя
	r.POST("/update-user", func(c *gin.Context) {
		controllers.UpdateUserHandler(c.Writer, c.Request)
	})

	// Загрузка шаблона main.html
	r.LoadHTMLGlob("templates/*")
}
