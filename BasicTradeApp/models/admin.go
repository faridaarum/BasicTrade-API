package models

import (
	"time"
)

type Admin struct {
	ID        uint      `gorm:"primaryKey"`
	UUID      string    `gorm:"type:char(36);not null"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdateAt  time.Time `gorm:"autoUpdateTime"`
}
