
package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	Name        string `json:"name" gorm:"size:100;not null"`
	Description string `json:"description" gorm:"type:text"`

	Products []Product `json:"products,omitempty"`
}
