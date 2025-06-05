package models

import "time"

type Activity struct {
	ID        uint      `json:"activityId" gorm:"primaryKey"`
	UserID    uint      `json:"userId"`
	Title     string    `json:"title" gorm:"not null"`
	Value     int       `json:"value"`
	CreatedAt time.Time `json:"-"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
}
