package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Name     string    `json:"name" validate:"required"`
	PhotoURL string    `json:"photo_url"`
	AdminID  uint      `json:"admin_id"`
	Variants []Variant `json:"variants"`
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return
}
