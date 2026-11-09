
package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model

	CartID    uint `json:"cart_id" gorm:"uniqueIndex:idx_cart_product;not null"`
	ProductID uint `json:"product_id" gorm:"uniqueIndex:idx_cart_product;not null"`
	Quantity  int  `json:"quantity" gorm:"not null;default:1"`

	Product Product `json:"product,omitempty"`
}
