package controllers

import (
	"fmt"
	"gofp/models"
	// "gofp/utils"
	"net/http"
	// "strconv"
	"strings"
	// "encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

var DB *gorm.DB
var jwtKey = []byte("B0qtauga-Bo1a=ma?kot1nd1/qyswy+nelzya")

func SignUp(c *gin.Context) {
	var u models.User
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	u.ID = uuid.New().String()
	u.Password = string(hashedPassword)
	// Assuming 'db' is your GORM DB connection, save user to database
	DB.Create(&u)

	c.JSON(http.StatusCreated, u)
}

func SignIn(c *gin.Context) {
	var u models.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedUser models.User
	if err := DB.Where("username = ?", u.Username).First(&storedUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(u.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": u.Username,
		"iss":      "go-jwt-server",
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}

	c.Header("Authorization", "Bearer "+tokenString)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello, world!")
}
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func APIHandler(c *gin.Context) {
	c.String(http.StatusOK, "Access granted to the API.")
}
// func CreateUserHandler(c *gin.Context) {
// 	if c.Request.Method != http.MethodPost {
// 		c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
// 		return
// 	}

// 	// Parse form data
// 	err := c.Request.ParseForm()
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Failed to parse form data")
// 		return
// 	}

// 	name := c.Request.Form.Get("name")
// 	email := c.Request.Form.Get("email")
// 	ageStr := c.Request.Form.Get("age")

// 	// Convert age string to int
// 	age, err := strconv.Atoi(ageStr)
// 	if err != nil {
// 		c.String(http.StatusBadRequest, "Invalid age")
// 		return
// 	}

// 	// Create user instance
// 	user := models.User{Name: name, Email: email, Age: age}

// 	// Save the user to the database
// 	err = utils.CreateUsers(db, &user)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
// 		return
// 	}

// 	// Render the main.html page with a success message
// 	c.HTML(http.StatusOK, "main.html", gin.H{
// 		"message": "User created successfully",
// 	})
// }

// func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
// 	// Проверяем метод запроса
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Читаем данные из формы
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
// 		return
// 	}

// 	// Получаем данные из формы
// 	idStr := r.Form.Get("userId")
// 	newName := r.Form.Get("newName")
// 	newEmail := r.Form.Get("newEmail")
// 	newAgeStr := r.Form.Get("newAge")

// 	// Преобразуем строковые данные в нужные типы
// 	id, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}
// 	newAge, err := strconv.Atoi(newAgeStr)
// 	if err != nil {
// 		http.Error(w, "Invalid age", http.StatusBadRequest)
// 		return
// 	}

// 	// Обновляем данные пользователя
// 	err = utils.UpdateUserByID(db, uint(id), newName, newEmail, newAge)
// 	if err != nil {
// 		// Если произошла ошибка при обновлении пользователя, выводим сообщение об ошибке
// 		http.Error(w, fmt.Sprintf("Failed to update user: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Если обновление прошло успешно, выводим успешное сообщение
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("User updated successfully"))
// }

// func GetUserByIDHandler(c *gin.Context) {
//     // Получаем ID пользователя из параметра запроса
//     userID := c.Query("userID")

//     // Преобразуем ID пользователя в тип uint
//     id, err := strconv.ParseUint(userID, 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
//         return
//     }

//     // Получаем пользователя по ID
//     user, err := utils.GetUserByID(db, uint(id))
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     // Отправляем пользователя в виде JSON ответа
//     c.JSON(http.StatusOK, user)
// }

// func AllUsersHandler(c *gin.Context) {
// 	var users []models.User
// 	if err := db.Find(&users).Error; err != nil {
// 		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
// 			"message": "Failed to fetch users from the database",
// 		})
// 		return
// 	}

// 	c.HTML(http.StatusOK, "main.html", gin.H{
// 		"users": users,
// 	})
// }
