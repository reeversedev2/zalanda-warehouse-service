package models

import "github.com/jinzhu/gorm"

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
