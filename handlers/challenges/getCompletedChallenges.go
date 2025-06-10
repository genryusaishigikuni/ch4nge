package challenges

import (
	"log"
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
		log.Printf("Invalid user ID provided: %v", userIdStr) // Log the invalid user ID error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	log.Printf("Fetching completed challenges for user ID: %d", userId) // Log the incoming request for completed challenges

	var response CompletedChallengesResponse

	// Fetch completed mini challenges
	if err := database.DB.Where("user_id = ? AND completed = ?", userId, true).
		Preload("MiniChallenge").
		Order("completed_at DESC").
		Find(&response.MiniChallenges).Error; err != nil {
		log.Printf("Error fetching mini challenges for user %d: %v", userId, err) // Log database query error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mini challenges"})
		return
	}

	// Fetch completed weekly challenges
	if err := database.DB.Where("user_id = ? AND completed = ?", userId, true).
		Preload("WeeklyChallenge").
		Order("completed_at DESC").
		Find(&response.WeeklyChallenges).Error; err != nil {
		log.Printf("Error fetching weekly challenges for user %d: %v", userId, err) // Log database query error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weekly challenges"})
		return
	}

	// Calculate the total mini challenge points
	var miniPoints int
	if err := database.DB.Table("user_mini_challenges").
		Select("COALESCE(SUM(mini_challenges.points), 0)").
		Joins("JOIN mini_challenges ON user_mini_challenges.mini_challenge_id = mini_challenges.id").
		Where("user_mini_challenges.user_id = ? AND user_mini_challenges.completed = ?", userId, true).
		Scan(&miniPoints).Error; err != nil {
		log.Printf("Error calculating mini challenge points for user %d: %v", userId, err) // Log points calculation error
	}

	// Calculate the total weekly challenge points
	var weeklyPoints int
	if err := database.DB.Table("user_weekly_challenges").
		Select("COALESCE(SUM(weekly_challenges.points), 0)").
		Joins("JOIN weekly_challenges ON user_weekly_challenges.weekly_challenge_id = weekly_challenges.id").
		Where("user_weekly_challenges.user_id = ? AND user_weekly_challenges.completed = ?", userId, true).
		Scan(&weeklyPoints).Error; err != nil {
		log.Printf("Error calculating weekly challenge points for user %d: %v", userId, err) // Log points calculation error
	}

	// Add the points for both mini and weekly challenges
	response.TotalPoints = miniPoints + weeklyPoints
	log.Printf("User ID: %d has a total of %d points from completed challenges", userId, response.TotalPoints) // Log total points

	// Calculate the total completed mini challenges
	if err := database.DB.Model(&models.UserMiniChallenge{}).
		Where("user_id = ? AND completed = ?", userId, true).
		Count(&response.CompletionStats.TotalMiniCompleted).Error; err != nil {
		log.Printf("Error counting completed mini challenges for user %d: %v", userId, err) // Log error counting mini challenges
	}

	// Calculate the total completed weekly challenges
	if err := database.DB.Model(&models.UserWeeklyChallenge{}).
		Where("user_id = ? AND completed = ?", userId, true).
		Count(&response.CompletionStats.TotalWeeklyCompleted).Error; err != nil {
		log.Printf("Error counting completed weekly challenges for user %d: %v", userId, err) // Log error counting weekly challenges
	}

	// Calculate total challenges assigned to the user
	var totalAssigned int64
	if err := database.DB.Model(&models.UserMiniChallenge{}).Where("user_id = ?", userId).Count(&totalAssigned).Error; err != nil {
		log.Printf("Error counting total mini challenges assigned to user %d: %v", userId, err) // Log error counting assigned mini challenges
	}

	var totalWeeklyAssigned int64
	if err := database.DB.Model(&models.UserWeeklyChallenge{}).Where("user_id = ?", userId).Count(&totalWeeklyAssigned).Error; err != nil {
		log.Printf("Error counting total weekly challenges assigned to user %d: %v", userId, err) // Log error counting assigned weekly challenges
	}

	// Calculate completion rate
	totalChallenges := totalAssigned + totalWeeklyAssigned
	totalCompleted := response.CompletionStats.TotalMiniCompleted + response.CompletionStats.TotalWeeklyCompleted
	if totalChallenges > 0 {
		response.CompletionStats.CompletionRate = float64(totalCompleted) / float64(totalChallenges) * 100
	}

	log.Printf("User ID: %d has completed %d out of %d challenges (%.2f%% completion rate)", userId, totalCompleted, totalChallenges, response.CompletionStats.CompletionRate) // Log completion rate

	// Send the response with all data
	c.JSON(http.StatusOK, response)
}
