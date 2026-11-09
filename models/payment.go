
package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model

	OrderID       uint   `json:"order_id" gorm:"not null;index"`
	ReferenceCode string `json:"reference_code" gorm:"size:100;uniqueIndex;not null"`
	PaymentMethod string `json:"payment_method" gorm:"size:50;not null"`
	Amount        int64  `json:"amount" gorm:"not null"`
	Status        string `json:"status" gorm:"size:30;default:PENDING"`
}
