package models

import "gorm.io/gorm"

type Variant struct {
	gorm.Model
	Name      string `json:"name" validate:"required"`
	ProductID uint   `json:"product_id"`
}
