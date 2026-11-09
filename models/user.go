
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name     string `json:"name" gorm:"size:100;not null"`
	Email    string `json:"email" gorm:"size:150;uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"size:20;default:customer"`
}
