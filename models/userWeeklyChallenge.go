package models

import "time"

type UserWeeklyChallenge struct {
	ID                uint            `json:"-" gorm:"primaryKey"`
	UserID            uint            `json:"userId"`
	WeeklyChallengeID uint            `json:"weeklyChallengeId"`
	CurrentValue      float64         `json:"currentValue"` // Убрали default значение
	IsCompleted       bool            `json:"isCompleted" gorm:"default:false"`
	CompletedAt       *time.Time      `json:"completedAt"`
	AssignedAt        time.Time       `json:"-"`
	WeeklyChallenge   WeeklyChallenge `json:"-" gorm:"foreignKey:WeeklyChallengeID"`
	User              User            `json:"-" gorm:"foreignKey:UserID"`
}
