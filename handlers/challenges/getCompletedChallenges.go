package challenges

import (
	"net/http"
	"strconv"

	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

type CompletedChallengesResponse struct {
	MiniChallenges   []models.UserMiniChallenge   `json:"mini_challenges"`
	WeeklyChallenges []models.UserWeeklyChallenge `json:"weekly_challenges"`
	TotalPoints      int                          `json:"total_points"`
	CompletionStats  CompletionStats              `json:"completion_stats"`
}

type CompletionStats struct {
	TotalMiniCompleted   int64   `json:"total_mini_completed"`
	TotalWeeklyCompleted int64   `json:"total_weekly_completed"`
	CompletionRate       float64 `json:"completion_rate"`
}

func GetCompletedChallenges(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var response CompletedChallengesResponse

	// Get completed mini challenges
	database.DB.Where("user_id = ? AND completed = ?", userId, true).
		Preload("MiniChallenge").
		Order("completed_at DESC").
		Find(&response.MiniChallenges)

	// Get completed weekly challenges
	database.DB.Where("user_id = ? AND completed = ?", userId, true).
		Preload("WeeklyChallenge").
		Order("completed_at DESC").
		Find(&response.WeeklyChallenges)

	// Calculate total points from completed challenges
	var miniPoints int
	database.DB.Table("user_mini_challenges").
		Select("COALESCE(SUM(mini_challenges.points), 0)").
		Joins("JOIN mini_challenges ON user_mini_challenges.mini_challenge_id = mini_challenges.id").
		Where("user_mini_challenges.user_id = ? AND user_mini_challenges.completed = ?", userId, true).
		Scan(&miniPoints)

	var weeklyPoints int
	database.DB.Table("user_weekly_challenges").
		Select("COALESCE(SUM(weekly_challenges.points), 0)").
		Joins("JOIN weekly_challenges ON user_weekly_challenges.weekly_challenge_id = weekly_challenges.id").
		Where("user_weekly_challenges.user_id = ? AND user_weekly_challenges.completed = ?", userId, true).
		Scan(&weeklyPoints)

	response.TotalPoints = miniPoints + weeklyPoints

	// Calculate completion stats
	database.DB.Model(&models.UserMiniChallenge{}).
		Where("user_id = ? AND completed = ?", userId, true).
		Count(&response.CompletionStats.TotalMiniCompleted)

	database.DB.Model(&models.UserWeeklyChallenge{}).
		Where("user_id = ? AND completed = ?", userId, true).
		Count(&response.CompletionStats.TotalWeeklyCompleted)

	// Calculate completion rate
	var totalAssigned int64
	database.DB.Model(&models.UserMiniChallenge{}).Where("user_id = ?", userId).Count(&totalAssigned)

	var totalWeeklyAssigned int64
	database.DB.Model(&models.UserWeeklyChallenge{}).Where("user_id = ?", userId).Count(&totalWeeklyAssigned)

	totalChallenges := totalAssigned + totalWeeklyAssigned
	totalCompleted := response.CompletionStats.TotalMiniCompleted + response.CompletionStats.TotalWeeklyCompleted

	if totalChallenges > 0 {
		response.CompletionStats.CompletionRate = float64(totalCompleted) / float64(totalChallenges) * 100
	}

	c.JSON(http.StatusOK, response)
}
