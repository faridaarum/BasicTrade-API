package models

import (
	"time"
)

type Product struct {
	ID        uint      `gorm:"primaryKey"`
	UUID      string    `gorm:"type:char(36);not null"`
	Name      string    `gorm:"type:varchar(100);not null"`
	ImageURL  string    `gorm:"type:varchar(255)"`
	AdminID   uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Variants  []Variant `gorm:"foreignKey:ProductID"`
}
