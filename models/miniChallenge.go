package models

import "time"

type MiniChallenge struct {
	ID        uint      `json:"miniChallengeId" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Subtitle  string    `json:"subtitle"`
	Points    int       `json:"points" gorm:"default:0"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
