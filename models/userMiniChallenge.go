package models

import "time"

type UserMiniChallenge struct {
	ID              uint          `json:"-" gorm:"primaryKey"`
	UserID          uint          `json:"userId"`
	MiniChallengeID uint          `json:"miniChallengeId"`
	IsAchieved      bool          `json:"isAchieved" gorm:"default:false"`
	AchievedAt      *time.Time    `json:"achievedAt"`
	AssignedAt      time.Time     `json:"-"`
	MiniChallenge   MiniChallenge `json:"-" gorm:"foreignKey:MiniChallengeID"`
	User            User          `json:"-" gorm:"foreignKey:UserID"`
}
