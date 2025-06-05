package models

import "time"

type Action struct {
	ID         uint                   `json:"actionId" gorm:"primaryKey"`
	UserID     uint                   `json:"userId"`
	ActionType string                 `json:"actionType" gorm:"not null"`
	Payload    map[string]interface{} `json:"payload" gorm:"type:jsonb"`
	Metadata   map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	Points     int                    `json:"points" gorm:"default:0"`
	CreatedAt  time.Time              `json:"-"`
	User       User                   `json:"-" gorm:"foreignKey:UserID"`
}
