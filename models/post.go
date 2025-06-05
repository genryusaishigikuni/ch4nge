package models

import "time"

type Post struct {
	ID           uint      `json:"postId" gorm:"primaryKey"`
	UserID       uint      `json:"userId"`
	Title        string    `json:"title" gorm:"not null"`
	ImageURL     string    `json:"imageUrl"`
	LikeNumber   int       `json:"likeNumber" gorm:"default:0"`
	SharesNumber int       `json:"sharesNumber" gorm:"default:0"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	User         User      `json:"-" gorm:"foreignKey:UserID"`
}
