
package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model

	OrderID   uint  `json:"order_id" gorm:"not null;index"`
	ProductID uint  `json:"product_id" gorm:"not null"`
	Quantity  int   `json:"quantity" gorm:"not null"`
	Price     int64 `json:"price" gorm:"not null"`
	Subtotal  int64 `json:"subtotal" gorm:"not null"`
	Product Product `json:"product,omitempty"`
}
