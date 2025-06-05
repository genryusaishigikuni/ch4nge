package models

type PostLike struct {
	ID     uint `gorm:"primaryKey"`
	PostID uint `gorm:"not null"`
	UserID uint `gorm:"not null"`
	Post   Post `gorm:"foreignKey:PostID"`
	User   User `gorm:"foreignKey:UserID"`
}
