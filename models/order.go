
package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	UserID      uint   `json:"user_id" gorm:"not null;index"`
	TotalAmount int64  `json:"total_amount" gorm:"not null"`
	Status      string `json:"status" gorm:"size:30;default:PENDING_PAYMENT"`

	Items []OrderItem `json:"items,omitempty"`
}
