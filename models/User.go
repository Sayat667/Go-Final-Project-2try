package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	ID       string `gorm:"unique"`
	Email    string `gorm:"primaryKey"`
	Password string 
	Balance  float64
	Perm     Permissions
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
