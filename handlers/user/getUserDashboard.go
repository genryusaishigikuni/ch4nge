package User

import (
	"net/http"
	"strconv"
	"time"

	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func GetUserDashboard(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var dashboard models.DashboardStats

	// Get user info
	if err := database.DB.First(&dashboard.User, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get total actions
	database.DB.Model(&models.Action{}).Where("user_id = ?", userId).Count(&dashboard.TotalActions)

	// Get actions this week
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	database.DB.Model(&models.Action{}).
		Where("user_id = ? AND created_at >= ?", userId, weekStart).
		Count(&dashboard.ActionsThisWeek)

	// Calculate total points
	var achievementPoints int
	database.DB.Table("user_achievements").
		Select("COALESCE(SUM(achievements.points), 0)").
		Joins("JOIN achievements ON user_achievements.achievement_id = achievements.id").
		Where("user_achievements.user_id = ?", userId).
		Scan(&achievementPoints)

	var challengePoints int
	database.DB.Table("user_mini_challenges").
		Select("COALESCE(SUM(mini_challenges.points), 0)").
		Joins("JOIN mini_challenges ON user_mini_challenges.mini_challenge_id = mini_challenges.id").
		Where("user_mini_challenges.user_id = ? AND user_mini_challenges.completed = ?", userId, true).
		Scan(&challengePoints)

	dashboard.TotalPoints = achievementPoints + challengePoints

	// Get current streak
	dashboard.CurrentStreak = calculateActionStreak(userId)

	// Get achievements count
	database.DB.Model(&models.UserAchievement{}).Where("user_id = ?", userId).Count(&dashboard.AchievementsCount)

	// Get next achievement (first unearned achievement)
	var nextAchievement models.Achievement
	if err := database.DB.Where("id NOT IN (SELECT achievement_id FROM user_achievements WHERE user_id = ?)", userId).
		Where("is_active = ?", true).
		Order("threshold ASC").
		First(&nextAchievement).Error; err == nil {
		dashboard.NextAchievement = &nextAchievement
	}

	// Get active mini challenges
	database.DB.Where("user_id = ? AND completed = ?", userId, false).
		Preload("MiniChallenge").
		Find(&dashboard.ActiveMiniChallenges)

	// Get active weekly challenge
	var weeklyChallenge models.UserWeeklyChallenge
	if err := database.DB.Where("user_id = ? AND completed = ?", userId, false).
		Preload("WeeklyChallenge").
		First(&weeklyChallenge).Error; err == nil {
		dashboard.ActiveWeeklyChallenge = &weeklyChallenge
	}

	// Get recent actions (last 5)
	database.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(5).
		Find(&dashboard.RecentActions)

	// Get friends count
	database.DB.Model(&models.UserFriend{}).Where("user_id = ?", userId).Count(&dashboard.FriendsCount)

	// Get posts count
	database.DB.Model(&models.Post{}).Where("user_id = ?", userId).Count(&dashboard.PostsCount)

	// Get weekly progress (actions per day for the current week)
	dashboard.WeeklyProgress = make(map[string]int64)
	for i := 0; i < 7; i++ {
		day := weekStart.AddDate(0, 0, i)
		dayEnd := day.Add(24 * time.Hour)

		var dayActions int64
		database.DB.Model(&models.Action{}).
			Where("user_id = ? AND created_at >= ? AND created_at < ?", userId, day, dayEnd).
			Count(&dayActions)

		dashboard.WeeklyProgress[day.Format("2006-01-02")] = dayActions
	}

	c.JSON(http.StatusOK, dashboard)
}
