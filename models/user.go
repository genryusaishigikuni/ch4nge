package models

import "time"

type User struct {
	ID            uint      `json:"userId" gorm:"primaryKey"`
	Username      string    `json:"username" gorm:"unique;not null"`
	Email         string    `json:"email" gorm:"unique;not null"`
	Password      string    `json:"-" gorm:"not null"`
	ProfilePicURL string    `json:"profilePicUrl"`
	Streak        int       `json:"streak" gorm:"default:0"`
	Points        int       `json:"points" gorm:"default:0"`
	GHGIndex      float64   `json:"ghgIndex" gorm:"default:0.0"`
	Latitude      float64   `json:"-"`
	Longitude     float64   `json:"-"`
	IsAdmin       bool      `json:"-" gorm:"default:false"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	Friends       []*User   `json:"-" gorm:"many2many:user_friends;constraint:OnDelete:CASCADE"`
}
