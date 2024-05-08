package controllers

import (
	// "fmt"
	"gofp/utils"
	"gofp/models"
	// h "gofp/helpers"
	// "gofp/utils"
	// "net/http"
	// "strconv"
	// "strings"
	// "encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var DB *gorm.DB
var jwtKey = []byte("B0qtauga-Bo1a=ma?kot1nd1/qyswy+nelzya")


func HomePage(c *gin.Context){
	// СЮДА КОД ГЛАВНОЙ СТРАНИЦЫ ................ Damir Aitbay
	c.JSON(200, gin.H{
        "message": "That is the main page",
    })


	
    cookie, err := c.Cookie("token")

    if err != nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    claims, err := utils.ParseToken(cookie)

    if err != nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    if claims.Role != "user" && claims.Role != "admin" {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    c.JSON(200, gin.H{"success": "home page", "role": claims.Role})
}

func Premium(c *gin.Context) {

    cookie, err := c.Cookie("token")

    if err != nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    claims, err := utils.ParseToken(cookie)

    if err != nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    if claims.Role != "admin" {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    c.JSON(200, gin.H{"success": "premium page", "role": claims.Role})
}

func SignUp(c *gin.Context) {
	var user models.User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var existingUser models.User

    models.DB.Where("email = ?", user.Email).First(&existingUser)

    if existingUser.ID != 0 {
        c.JSON(400, gin.H{"error": "user already exists"})
        return
    }

    var errHash error
    user.Password, errHash = utils.GenerateHashPassword(user.Password)

    if errHash != nil {
        c.JSON(500, gin.H{"error": "could not generate password hash"})
        return
    }

    models.DB.Create(&user)

    c.JSON(200, gin.H{"success": "user created"})
}

func SignIn(c *gin.Context) {
	
    var user models.User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var existingUser models.User

    models.DB.Where("email = ?", user.Email).First(&existingUser)

    if existingUser.ID == 0 {
        c.JSON(400, gin.H{"error": "user does not exist"})
        return
    }

    errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

    if !errHash {
        c.JSON(400, gin.H{"error": "invalid password"})
        return
    }

    expirationTime := time.Now().Add(5 * time.Minute)

    claims := &models.Claims{
        Role: existingUser.Role,
        StandardClaims: jwt.StandardClaims{
            Subject:   existingUser.Email,
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString(jwtKey)

    if err != nil {
        c.JSON(500, gin.H{"error": "could not generate token"})
        return
    }

    c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
    c.JSON(200, gin.H{"success": "user logged in"})
}

func Logout(c *gin.Context) {
    c.SetCookie("token", "", -1, "/", "localhost", false, true)
    c.JSON(200, gin.H{"success": "user logged out"})
}

// func JWTMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method")
// 			}
// 			return jwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		c.Next()
// 	}
// }

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
