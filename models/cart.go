
package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model

	UserID uint `json:"user_id" gorm:"uniqueIndex;not null"`

	Items []CartItem `json:"items,omitempty"`
}
