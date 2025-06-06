package models

type DashboardStats struct {
	User                  User                 `json:"user"`
	TotalPoints           int                  `json:"total_points"`
	TotalActions          int64                `json:"total_actions"`
	ActionsThisWeek       int64                `json:"actions_this_week"`
	CurrentStreak         int                  `json:"current_streak"`
	AchievementsCount     int64                `json:"achievements_count"`
	NextAchievement       *Achievement         `json:"next_achievement,omitempty"`
	ActiveMiniChallenges  []UserMiniChallenge  `json:"active_mini_challenges"`
	ActiveWeeklyChallenge *UserWeeklyChallenge `json:"active_weekly_challenge,omitempty"`
	RecentActions         []Action             `json:"recent_actions"`
	FriendsCount          int64                `json:"friends_count"`
	PostsCount            int64                `json:"posts_count"`
	WeeklyProgress        map[string]int64     `json:"weekly_progress"`
}
