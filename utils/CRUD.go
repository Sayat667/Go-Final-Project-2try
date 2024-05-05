package utils

import (
	"fmt"
	"gofp/models"
    // "errors"

	"gorm.io/gorm"
)

func CreateUsers(db *gorm.DB, users ...*models.User) error {
    if len(users) == 0 {
        return fmt.Errorf("no users provided to create")
    }

    for _, user := range users {
        result := db.Create(user)
        if result.Error != nil {
            return fmt.Errorf("failed to create user: %v", result.Error)
        }
        fmt.Printf("User created successfully: %+v\n", *user) // Логируем созданных пользователей
    }

    return nil
}

func GetUserByID(db *gorm.DB, id uint) (*models.User, error) {
    var user models.User
    if err := db.First(&user, id).Error; err != nil {
        return nil, fmt.Errorf("failed to retrieve user: %v", err)
    }
    fmt.Printf("Retrieved user: %+v\n", user)
    return &user, nil
}

func UpdateUserByID(db *gorm.DB, id uint, newName string, newEmail string, newAge int) error {
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}

	user.Name = newName
	user.Age = newAge
	user.Email = newEmail

	result := db.Save(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %v", result.Error)
	}

	fmt.Println("User updated successfully")
	return nil
}

func DeleteUser(db *gorm.DB, user *models.User) error {
	result := db.Delete(user)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %v", result.Error)
	}
	fmt.Println("User deleted successfully")
	return nil
}

func DeleteUsersByID(db *gorm.DB, ids ...uint) error {
    if len(ids) == 0 {
        return fmt.Errorf("no user IDs provided for deletion")
    }

    for _, id := range ids {
        // Создаем пустую структуру User с указанным ID для удаления
        user := models.User{ID: id}

        // Удаляем пользователя по ID
        result := db.Delete(&user)
        if result.Error != nil {
            return fmt.Errorf("failed to delete user with ID %d: %v", id, result.Error)
        }
        fmt.Printf("User with ID %d deleted successfully\n", id)
    }

    return nil
}

