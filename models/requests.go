package models

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateFriendsRequest struct {
	FriendIds []string `json:"friendIds"`
}

type FriendsActivityRequest struct {
	UserIds []string `json:"userIds"`
}

type GreenActionRequest struct {
	ActionType string                 `json:"actionType" binding:"required"`
	Payload    map[string]interface{} `json:"payload" binding:"required"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type TransportationActionRequest struct {
	ActionType string                 `json:"actionType" binding:"required"`
	Payload    map[string]interface{} `json:"payload" binding:"required"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type PostRequest struct {
	UserID string `json:"userId" form:"userId" binding:"required"`
	Title  string `json:"title" form:"title" binding:"required"`
}

type LikePostRequest struct {
	UserID string `json:"userId" binding:"required"`
}

type SharePostRequest struct {
	UserID string `json:"userId" binding:"required"`
}
