package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// UintArray is a custom type for storing array of uint in database
type UintArray []uint

// Scan implements the Scanner interface for database reads
func (a *UintArray) Scan(value interface{}) error {
	if value == nil {
		*a = UintArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return errors.New("cannot scan into UintArray")
	}
}

// Value implements the driver Valuer interface for database writes
func (a UintArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "[]", nil
	}
	return json.Marshal(a)
}

type Post struct {
	ID           uint      `json:"postId" gorm:"primaryKey"`
	UserID       uint      `json:"userId"`
	Title        string    `json:"title" gorm:"not null"`
	ImageURL     string    `json:"imageUrl"`
	LikedBy      UintArray `json:"likedBy" gorm:"type:json;default:'[]'"`
	SharedBy     UintArray `json:"sharedBy" gorm:"type:json;default:'[]'"`
	LikeNumber   int       `json:"likeNumber" gorm:"-"`   // Computed field, not stored
	SharesNumber int       `json:"sharesNumber" gorm:"-"` // Computed field, not stored
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Description  string    `json:"description"`
	User         User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// AfterFind hook to calculate like and share counts
func (p *Post) AfterFind() error {
	p.LikeNumber = len(p.LikedBy)
	p.SharesNumber = len(p.SharedBy)
	return nil
}

// HasLiked checks if a user has liked this post
func (p *Post) HasLiked(userID uint) bool {
	for _, id := range p.LikedBy {
		if id == userID {
			return true
		}
	}
	return false
}

// AddLike adds a user's like to the post
func (p *Post) AddLike(userID uint) bool {
	if p.HasLiked(userID) {
		return false // Already liked
	}
	p.LikedBy = append(p.LikedBy, userID)
	return true
}

// RemoveLike removes a user's like from the post
func (p *Post) RemoveLike(userID uint) bool {
	for i, id := range p.LikedBy {
		if id == userID {
			// Remove the element at index i
			p.LikedBy = append(p.LikedBy[:i], p.LikedBy[i+1:]...)
			return true
		}
	}
	return false // Not found
}

// ToggleLike toggles a user's like status and returns whether the post is now liked
func (p *Post) ToggleLike(userID uint) bool {
	if p.HasLiked(userID) {
		p.RemoveLike(userID)
		return false
	} else {
		p.AddLike(userID)
		return true
	}
}
