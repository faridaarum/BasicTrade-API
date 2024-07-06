package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Variant struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Name      string    `json:"name" validate:"required"`
	ProductID uuid.UUID `json:"product_id"`
}

func (variant *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	variant.ID = uuid.New()
	return
}
