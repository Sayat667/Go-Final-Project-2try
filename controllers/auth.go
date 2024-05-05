package controllers

import (
	"fmt"
	"gofp/models"
	"gofp/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(DB *gorm.DB) {
	db = DB
}

func CreateUserHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Parse form data
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse form data")
		return
	}

	name := c.Request.Form.Get("name")
	email := c.Request.Form.Get("email")
	ageStr := c.Request.Form.Get("age")

	// Convert age string to int
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid age")
		return
	}

	// Create user instance
	user := models.User{Name: name, Email: email, Age: age}

	// Save the user to the database
	err = utils.CreateUsers(db, &user)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	// Render the main.html page with a success message
	c.HTML(http.StatusOK, "main.html", gin.H{
		"message": "User created successfully",
	})
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Читаем данные из формы
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}

	// Получаем данные из формы
	idStr := r.Form.Get("userId")
	newName := r.Form.Get("newName")
	newEmail := r.Form.Get("newEmail")
	newAgeStr := r.Form.Get("newAge")

	// Преобразуем строковые данные в нужные типы
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	newAge, err := strconv.Atoi(newAgeStr)
	if err != nil {
		http.Error(w, "Invalid age", http.StatusBadRequest)
		return
	}

	// Обновляем данные пользователя
	err = utils.UpdateUserByID(db, uint(id), newName, newEmail, newAge)
	if err != nil {
		// Если произошла ошибка при обновлении пользователя, выводим сообщение об ошибке
		http.Error(w, fmt.Sprintf("Failed to update user: %v", err), http.StatusInternalServerError)
		return
	}

	// Если обновление прошло успешно, выводим успешное сообщение
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}

func GetUserByIDHandler(c *gin.Context) {
    // Получаем ID пользователя из параметра запроса
    userID := c.Query("userID")

    // Преобразуем ID пользователя в тип uint
    id, err := strconv.ParseUint(userID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Получаем пользователя по ID
    user, err := utils.GetUserByID(db, uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Отправляем пользователя в виде JSON ответа
    c.JSON(http.StatusOK, user)
}

func AllUsersHandler(c *gin.Context) {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": "Failed to fetch users from the database",
		})
		return
	}

	c.HTML(http.StatusOK, "main.html", gin.H{
		"users": users,
	})
}
