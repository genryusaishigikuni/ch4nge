package models

import "time"

type UserAchievement struct {
	ID            uint        `json:"-" gorm:"primaryKey"`
	UserID        uint        `json:"userId"`
	AchievementID uint        `json:"achievementId"`
	IsAchieved    bool        `json:"isAchieved" gorm:"default:false"`
	AchievedAt    *time.Time  `json:"achievedAt"`
	Achievement   Achievement `json:"-" gorm:"foreignKey:AchievementID"`
	User          User        `json:"-" gorm:"foreignKey:UserID"`
}
