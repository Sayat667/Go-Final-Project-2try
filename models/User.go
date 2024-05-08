package models

import (
	"gorm.io/gorm"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
    gorm.Model
    Username     string `json:"name"`
    Email    string `gorm:"unique" json:"email"`
    Password string `json:"password"`
    Role     string `json:"role"`
}


type Claims struct {
    Role string `json:"role"`
    jwt.StandardClaims
}

type Permissions string

const (
	Admin  Permissions = "Admin"
	Client Permissions = "Client"
)

type Dish struct {
	gorm.Model
	Name        string
	Calories    float64
	Ingredients []string `gorm:"type:json"`
	Amount      int
}

type Drink struct {
	gorm.Model
	Name   string
	Amount int
}

type Set struct {
	gorm.Model
	FirstDishID  uint
	SecondDishID uint
	DrinkID      uint
	SaladID      uint
	FruitsID     uint
	Amount       int
}

type Basket struct {
	gorm.Model
	UserID   uint
	Total    int
	SetID    uint
	Drinks   []Drink `gorm:"foreignKey:BasketID"`
	Dishes   []Dish  `gorm:"foreignKey:BasketID"`
	TotalCal float64
}
