package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name     string    `json:"name" validate:"required"`
	PhotoURL string    `json:"photo_url"`
	AdminID  uint      `json:"admin_id"`
	Variants []Variant `json:"variants"`
}
