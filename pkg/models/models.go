package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Name      string `json:"name" gorm:"text;not null;default:null"`
	CompanyID int
	Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Price     float64 `json:"price" gorm:"float;not null;default:0"`
	Category  string  `json:"category" gorm:"text;not null;default:null"`
	Expire    string  `json:"expire" gorm:"text;not null"`
	Status    string  `json:"status" gorm:"text;not null;default:null"`
	Image     string  `json:"image" gorm:"text;not null;default:'https://via.placeholder.com/150'"`
}

type Company struct {
	gorm.Model
	Name  string `json:"name" gorm:"text;not null;default:null"`
	Image string `json:"image" gorm:"text;default:'https://via.placeholder.com/150'"`
}

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"uniqueIndex;not null;default:null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"default:'employee'"`
	Name     string `json:"name" gorm:"not null"`
	EmployeeID string `json:"employee_id" gorm:"uniqueIndex;not null;default:null"`
	IsActive   bool   `json:"is_active" gorm:"default:true"`
	LastLoginAt *time.Time `json:"last_login_at" gorm:"default:null"`
	CreatedAt   *time.Time `json:"created_at" gorm:"default:null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"default:null"`
}

type Permissions struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	Action string `json:"action" gorm:"not null"` // create, read, update, delete
	CreatedAt *time.Time `json:"created_at" gorm:"default:null"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"default:null"`
}
