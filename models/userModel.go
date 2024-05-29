package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	Email     string    `gorm:"unique" json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required"`
	Expenses  []Expense `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" `
	Incomes   []Income  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" `
}
type Expense struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	UserID    uint      `json:"user_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required"`
	Category  string    `json:"category" gorm:"default" validate:"required"`
}
type Income struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	UserID    uint      `json:"user_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required"`
	Category  string    `json:"category" gorm:"default" validate:"required"`
}
