
package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model

	Name        string `json:"name" gorm:"size:150;not null"`
	Description string `json:"description" gorm:"type:text"`
	Price       int64  `json:"price" gorm:"not null"`
	Stock       int    `json:"stock" gorm:"not null;default:0"`
	Image       string `json:"image"`

	CategoryID uint     `json:"category_id" gorm:"not null"`
	Category   Category `json:"category,omitempty"`
}
