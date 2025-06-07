package models

import "time"

type WeeklyChallenge struct {
	ID          uint      `json:"weeklyChallengeId" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Subtitle    string    `json:"subtitle"`
	Points      int       `json:"points" gorm:"default:0"`
	TargetValue float64   `json:"targetValue" gorm:"default:1.0"` // Изменено с int на float64
	IsActive    bool      `json:"isActive" gorm:"default:true"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
