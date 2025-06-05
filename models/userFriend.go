package models

type UserFriend struct {
	UserID   uint `gorm:"primaryKey"`
	FriendID uint `gorm:"primaryKey"`
}
