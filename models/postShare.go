package models

import "time"

type PostShare struct {
	ID       uint      `gorm:"primaryKey"`
	PostID   uint      `gorm:"not null"`
	UserID   uint      `gorm:"not null"`
	SharedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Post     Post      `gorm:"foreignKey:PostID"`
	User     User      `gorm:"foreignKey:UserID"`
}
