package User

import (
	"net/http"
	"strconv"
	"time"

	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

type ActionStats struct {
	TotalActions     int64            `json:"total_actions"`
	ActionsByType    map[string]int64 `json:"actions_by_type"`
	ActionsThisWeek  int64            `json:"actions_this_week"`
	ActionsThisMonth int64            `json:"actions_this_month"`
	TotalPoints      int              `json:"total_points"`
	RecentStreak     int              `json:"recent_streak"`
	LastActionDate   *time.Time       `json:"last_action_date"`
}

func GetUserActionStats(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var stats ActionStats

	// Get total actions count
	database.DB.Model(&models.Action{}).Where("user_id = ?", userId).Count(&stats.TotalActions)

	// Get actions by type
	var actionTypes []struct {
		ActionType string `json:"action_type"`
		Count      int64  `json:"count"`
	}

	database.DB.Model(&models.Action{}).
		Select("action_type, COUNT(*) as count").
		Where("user_id = ?", userId).
		Group("action_type").
		Scan(&actionTypes)

	stats.ActionsByType = make(map[string]int64)
	for _, at := range actionTypes {
		stats.ActionsByType[at.ActionType] = at.Count
	}

	// Get actions this week
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	database.DB.Model(&models.Action{}).
		Where("user_id = ? AND created_at >= ?", userId, weekStart).
		Count(&stats.ActionsThisWeek)

	// Get actions this month
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	database.DB.Model(&models.Action{}).
		Where("user_id = ? AND created_at >= ?", userId, monthStart).
		Count(&stats.ActionsThisMonth)

	// Calculate total points from user's achievements and challenges
	var totalPoints int
	database.DB.Table("user_achievements").
		Select("COALESCE(SUM(achievements.points), 0)").
		Joins("JOIN achievements ON user_achievements.achievement_id = achievements.id").
		Where("user_achievements.user_id = ?", userId).
		Scan(&totalPoints)

	var challengePoints int
	database.DB.Table("user_mini_challenges").
		Select("COALESCE(SUM(mini_challenges.points), 0)").
		Joins("JOIN mini_challenges ON user_mini_challenges.mini_challenge_id = mini_challenges.id").
		Where("user_mini_challenges.user_id = ? AND user_mini_challenges.completed = ?", userId, true).
		Scan(&challengePoints)

	stats.TotalPoints = totalPoints + challengePoints

	// Get last action date
	var lastAction models.Action
	if err := database.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		First(&lastAction).Error; err == nil {
		stats.LastActionDate = &lastAction.CreatedAt
	}

	// Calculate recent streak (consecutive days with actions)
	stats.RecentStreak = calculateActionStreak(userId)

	c.JSON(http.StatusOK, stats)
}

func calculateActionStreak(userId int) int {
	var actions []models.Action
	database.DB.Select("DATE(created_at) as action_date").
		Where("user_id = ?", userId).
		Group("DATE(created_at)").
		Order("action_date DESC").
		Find(&actions)

	if len(actions) == 0 {
		return 0
	}

	streak := 0
	currentDate := time.Now().Truncate(24 * time.Hour)

	for _, action := range actions {
		actionDate := action.CreatedAt.Truncate(24 * time.Hour)

		if actionDate.Equal(currentDate) || actionDate.Equal(currentDate.AddDate(0, 0, -1)) {
			streak++
			currentDate = currentDate.AddDate(0, 0, -1)
		} else {
			break
		}
	}

	return streak
}
