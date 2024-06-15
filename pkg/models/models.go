package models

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	Name     string  `json:"name" gorm:"text;not null;default:null"`
	Company  string  `json:"company" gorm:"text;not null;default:null"`
	Price    float64 `json:"price" gorm:"float;not null;default:0"`
	Category string  `json:"category" gorm:"text;not null;default:null"`
	Expire   string  `json:"expire" gorm:"text;not null;default:null"`
	Status   string  `json:"status" gorm:"text;not null;default:null"`
	Image    string  `json:"image" gorm:"text;not null;default:'https://via.placeholder.com/150'"`
}
