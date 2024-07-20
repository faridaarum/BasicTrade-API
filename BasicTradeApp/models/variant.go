package models

import (
	"time"
)

type Variant struct {
	ID          uint      `gorm:"primaryKey"`
	UUID        string    `gorm:"type:char(36);not null"`
	VariantName string    `gorm:"type:varchar(100);not null"`
	Quantity    int       `gorm:"not null"`
	ProductID   uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
