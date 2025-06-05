package models

import "time"

type Achievement struct {
	ID        uint      `json:"achievementId" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Subtitle  string    `json:"subtitle"`
	Points    int       `json:"points" gorm:"default:0"`
	Threshold int       `json:"threshold" gorm:"default:1"`
	Category  string    `json:"category"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
