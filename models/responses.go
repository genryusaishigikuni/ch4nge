package models

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

type RegisterResponse struct {
	Token string `json:"token,omitempty"`
}

type UserResponse struct {
	ID            uint      `json:"userId"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	ProfilePicURL string    `json:"profilePicUrl"`
	Streak        int       `json:"streak"`
	Points        int       `json:"points"`
	GHGIndex      float64   `json:"ghgIndex"`
	Location      []float64 `json:"location"`
	FriendsIds    []string  `json:"friendsIds"`
}

type AchievementResponse struct {
	AchievementID uint   `json:"achievementId"`
	UserID        uint   `json:"userId"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	IsAchieved    bool   `json:"isAchieved"`
}

type MiniChallengeResponse struct {
	MiniChallengeID uint   `json:"miniChallengeId"`
	UserID          uint   `json:"userId"`
	Title           string `json:"title"`
	Subtitle        string `json:"subtitle"`
	IsAchieved      bool   `json:"isAchieved"`
	Points          int    `json:"points"`
}

type WeeklyChallengeResponse struct {
	WeeklyChallengeID uint   `json:"weeklyChallengeId"`
	UserID            uint   `json:"userId"`
	Title             string `json:"title"`
	Subtitle          string `json:"subtitle"`
	CurrentValue      int    `json:"currentValue"`
	TotalValue        int    `json:"totalValue"`
	Points            int    `json:"points"`
}

type PostLikeResponse struct {
	Post    Post   `json:"post"`
	IsLiked bool   `json:"is_liked"`
	Message string `json:"message"`
}

type PostEngagementResponse struct {
	PostID      uint   `json:"post_id"`
	LikesCount  int    `json:"likes_count"`
	SharesCount int    `json:"shares_count"`
	LikedBy     []User `json:"liked_by,omitempty"`
}

type UserLikeStatusResponse struct {
	PostID  uint `json:"post_id"`
	UserID  uint `json:"user_id"`
	IsLiked bool `json:"is_liked"`
}
