package User

import (
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// UserActionStats представляет статистику действий пользователя
type userActionStats struct {
	TotalActions        int64               `json:"total_actions"`
	GreenActions        int64               `json:"green_actions"`
	TransportActions    int64               `json:"transport_actions"`
	TotalPoints         int                 `json:"total_points"`
	EcoTransportKm      float64             `json:"eco_transport_km"`
	TotalCO2Saved       float64             `json:"total_co2_saved"`
	WeeklyStats         WeeklyActionStats   `json:"weekly_stats"`
	MonthlyStats        MonthlyActionStats  `json:"monthly_stats"`
	AchievementCount    int64               `json:"achievement_count"`
	CompletedChallenges int64               `json:"completed_challenges"`
	ActionsByType       map[string]int64    `json:"actions_by_type"`
	CurrentStreak       int                 `json:"current_streak"`
	LongestStreak       int                 `json:"longest_streak"`
	LastActionDate      *time.Time          `json:"last_action_date"`
	EnvironmentalImpact EnvironmentalImpact `json:"environmental_impact"`
}

type WeeklyActionStats struct {
	Actions     int64   `json:"actions"`
	Points      int     `json:"points"`
	EcoDistance float64 `json:"eco_distance"`
	CO2Saved    float64 `json:"co2_saved"`
}

type MonthlyActionStats struct {
	Actions     int64   `json:"actions"`
	Points      int     `json:"points"`
	EcoDistance float64 `json:"eco_distance"`
	CO2Saved    float64 `json:"co2_saved"`
}

type EnvironmentalImpact struct {
	TreesEquivalent float64 `json:"trees_equivalent"`
	EnergyKwh       float64 `json:"energy_kwh"`
	WaterLiters     float64 `json:"water_liters"`
	WasteKg         float64 `json:"waste_kg"`
}

// GetUserActionStats возвращает подробную статистику действий пользователя
func GetUserActionStats(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var stats userActionStats

	// Общее количество действий
	database.DB.Model(&models.Action{}).Where("user_id = ?", userId).Count(&stats.TotalActions)

	// Количество зеленых действий
	database.DB.Model(&models.Action{}).Where("user_id = ? AND action_type = ?", userId, "green").Count(&stats.GreenActions)

	// Количество транспортных действий
	database.DB.Model(&models.Action{}).Where("user_id = ? AND action_type = ?", userId, "transportation").Count(&stats.TransportActions)

	// Общие поинты пользователя
	var user models.User
	if database.DB.First(&user, userId).Error == nil {
		stats.TotalPoints = user.Points
	}

	// Действия по типам
	calculateActionsByType(uint(userId), &stats)

	// Расчет экологичного транспорта и сэкономленного CO2
	calculateEcoTransportStats(uint(userId), &stats)

	// Недельная статистика
	calculateWeeklyStats(uint(userId), &stats.WeeklyStats)

	// Месячная статистика
	calculateMonthlyStats(uint(userId), &stats.MonthlyStats)

	// Количество достижений
	database.DB.Model(&models.UserAchievement{}).Where("user_id = ? AND is_achieved = ?", userId, true).Count(&stats.AchievementCount)

	// Количество выполненных челленджей
	var completedMini, completedWeekly int64
	database.DB.Model(&models.UserMiniChallenge{}).Where("user_id = ? AND is_completed = ?", userId, true).Count(&completedMini)
	database.DB.Model(&models.UserWeeklyChallenge{}).Where("user_id = ? AND is_completed = ?", userId, true).Count(&completedWeekly)
	stats.CompletedChallenges = completedMini + completedWeekly

	// Расчет серий действий
	calculateStreaks(uint(userId), &stats)

	// Последнее действие
	getLastActionDate(uint(userId), &stats)

	// Экологическое воздействие
	calculateEnvironmentalImpact(uint(userId), &stats)

	c.JSON(http.StatusOK, stats)
}

// calculateActionsByType подсчитывает действия по типам
func calculateActionsByType(userId uint, stats *userActionStats) {
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
}

// calculateEcoTransportStats вычисляет статистику экологичного транспорта
func calculateEcoTransportStats(userId uint, stats *userActionStats) {
	var transportActions []models.Action
	database.DB.Where("user_id = ? AND action_type = ?", userId, "transportation").Find(&transportActions)

	var totalEcoKm float64
	var totalCO2Saved float64

	for _, action := range transportActions {
		if distance, ok := action.Payload["distance"].(float64); ok {
			transportType, _ := action.Payload["transportType"].(string)
			fuelConsumption, _ := action.Payload["fuelConsumption"].(float64)
			passengers := 1
			if p, ok := action.Payload["passengers"].(float64); ok {
				passengers = int(p)
			}

			// Используем функцию из action пакета для расчета
			_, ghg, isEcoFriendly := calculateTransportImpact(distance, fuelConsumption, passengers, transportType)

			if isEcoFriendly {
				totalEcoKm += distance
			}

			// Расчет сэкономленного CO2 (сравнение с обычным автомобилем)
			standardCarCO2 := distance * 0.21 // средний автомобиль: 0.21 кг CO2/км
			if isEcoFriendly {
				totalCO2Saved += standardCarCO2 - ghg
			}
		}
	}

	stats.EcoTransportKm = totalEcoKm
	stats.TotalCO2Saved = totalCO2Saved
}

// calculateTransportImpact - упрощенная версия функции из action пакета
func calculateTransportImpact(distance, fuelConsumption float64, passengers int, transportType string) (int, float64, bool) {
	if passengers <= 0 {
		passengers = 1
	}

	var co2PerKm float64
	var isEcoFriendly bool

	switch transportType {
	case "bicycle", "walking", "scooter":
		co2PerKm = 0
		isEcoFriendly = true
	case "public_transport", "bus", "metro", "train":
		co2PerKm = 0.05
		isEcoFriendly = true
	case "electric_car":
		co2PerKm = 0.1
		isEcoFriendly = true
	default:
		co2PerKm = (fuelConsumption / 100.0) * 2.3
		isEcoFriendly = false
	}

	totalCO2 := co2PerKm * distance
	ghg := totalCO2 / float64(passengers)

	return 0, ghg, isEcoFriendly
}

// calculateWeeklyStats вычисляет недельную статистику
func calculateWeeklyStats(userId uint, stats *WeeklyActionStats) {
	weekStart := getStartOfCurrentWeek()

	// Количество действий за неделю
	database.DB.Model(&models.Action{}).
		Where("user_id = ? AND created_at >= ?", userId, weekStart).
		Count(&stats.Actions)

	// Поинты за неделю
	var weeklyPoints int
	database.DB.Model(&models.Action{}).
		Select("COALESCE(SUM(points), 0)").
		Where("user_id = ? AND created_at >= ?", userId, weekStart).
		Scan(&weeklyPoints)
	stats.Points = weeklyPoints

	// Экологичное расстояние за неделю
	var weeklyActions []models.Action
	database.DB.Where("user_id = ? AND action_type = ? AND created_at >= ?", userId, "transportation", weekStart).Find(&weeklyActions)

	var weeklyEcoDistance, weeklyCO2Saved float64
	for _, action := range weeklyActions {
		if distance, ok := action.Payload["distance"].(float64); ok {
			transportType, _ := action.Payload["transportType"].(string)
			fuelConsumption, _ := action.Payload["fuelConsumption"].(float64)
			passengers := 1
			if p, ok := action.Payload["passengers"].(float64); ok {
				passengers = int(p)
			}

			_, ghg, isEcoFriendly := calculateTransportImpact(distance, fuelConsumption, passengers, transportType)

			if isEcoFriendly {
				weeklyEcoDistance += distance
			}

			standardCarCO2 := distance * 0.21
			if isEcoFriendly {
				weeklyCO2Saved += standardCarCO2 - ghg
			}
		}
	}

	stats.EcoDistance = weeklyEcoDistance
	stats.CO2Saved = weeklyCO2Saved
}

// calculateMonthlyStats вычисляет месячную статистику
func calculateMonthlyStats(userId uint, stats *MonthlyActionStats) {
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())

	// Количество действий за месяц
	database.DB.Model(&models.Action{}).
		Where("user_id = ? AND created_at >= ?", userId, monthStart).
		Count(&stats.Actions)

	// Поинты за месяц
	var monthlyPoints int
	database.DB.Model(&models.Action{}).
		Select("COALESCE(SUM(points), 0)").
		Where("user_id = ? AND created_at >= ?", userId, monthStart).
		Scan(&monthlyPoints)
	stats.Points = monthlyPoints

	// Экологичное расстояние за месяц
	var monthlyActions []models.Action
	database.DB.Where("user_id = ? AND action_type = ? AND created_at >= ?", userId, "transportation", monthStart).Find(&monthlyActions)

	var monthlyEcoDistance, monthlyCO2Saved float64
	for _, action := range monthlyActions {
		if distance, ok := action.Payload["distance"].(float64); ok {
			transportType, _ := action.Payload["transportType"].(string)
			fuelConsumption, _ := action.Payload["fuelConsumption"].(float64)
			passengers := 1
			if p, ok := action.Payload["passengers"].(float64); ok {
				passengers = int(p)
			}

			_, ghg, isEcoFriendly := calculateTransportImpact(distance, fuelConsumption, passengers, transportType)

			if isEcoFriendly {
				monthlyEcoDistance += distance
			}

			standardCarCO2 := distance * 0.21
			if isEcoFriendly {
				monthlyCO2Saved += standardCarCO2 - ghg
			}
		}
	}

	stats.EcoDistance = monthlyEcoDistance
	stats.CO2Saved = monthlyCO2Saved
}

// calculateStreaks вычисляет серии действий
func calculateStreaks(userId uint, stats *userActionStats) {
	var actions []models.Action
	database.DB.Select("DATE(created_at) as action_date").
		Where("user_id = ?", userId).
		Group("DATE(created_at)").
		Order("action_date DESC").
		Find(&actions)

	if len(actions) == 0 {
		stats.CurrentStreak = 0
		stats.LongestStreak = 0
		return
	}

	// Текущая серия
	currentStreak := 0
	longestStreak := 0
	tempStreak := 0
	currentDate := time.Now().Truncate(24 * time.Hour)

	for _, action := range actions {
		actionDate := action.CreatedAt.Truncate(24 * time.Hour)

		if actionDate.Equal(currentDate) || actionDate.Equal(currentDate.AddDate(0, 0, -1)) {
			if currentStreak == 0 || actionDate.Equal(currentDate.AddDate(0, 0, -currentStreak)) {
				currentStreak++
			}
			tempStreak++
			currentDate = currentDate.AddDate(0, 0, -1)
		} else {
			if tempStreak > longestStreak {
				longestStreak = tempStreak
			}
			tempStreak = 1
			currentDate = actionDate.AddDate(0, 0, -1)
		}
	}

	if tempStreak > longestStreak {
		longestStreak = tempStreak
	}

	stats.CurrentStreak = currentStreak
	stats.LongestStreak = longestStreak
}

// getLastActionDate получает дату последнего действия
func getLastActionDate(userId uint, stats *userActionStats) {
	var lastAction models.Action
	if err := database.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		First(&lastAction).Error; err == nil {
		stats.LastActionDate = &lastAction.CreatedAt
	}
}

// calculateEnvironmentalImpact вычисляет экологическое воздействие
func calculateEnvironmentalImpact(userId uint, stats *userActionStats) {
	var greenActions []models.Action
	database.DB.Where("user_id = ? AND action_type = ?", userId, "green").Find(&greenActions)

	var impact EnvironmentalImpact

	for _, action := range greenActions {
		option, _ := action.Payload["option"].(string)

		switch option {
		case "planted_tree":
			impact.TreesEquivalent += 1.0
		case "solar_power":
			impact.EnergyKwh += 10.0
		case "composting":
			impact.WasteKg += 5.0
		case "recycling":
			impact.WasteKg += 2.0
		case "water_conservation":
			impact.WaterLiters += 50.0
		case "energy_saving":
			impact.EnergyKwh += 5.0
		case "waste_reduction":
			impact.WasteKg += 3.0
		}
	}

	// Добавляем эквивалент деревьев на основе сэкономленного CO2
	if stats.TotalCO2Saved > 0 {
		// 1 дерево поглощает примерно 22 кг CO2 в год
		impact.TreesEquivalent += stats.TotalCO2Saved / 22.0
	}

	stats.EnvironmentalImpact = impact
}

// Вспомогательные функции
func getStartOfCurrentWeek() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	return now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
}

// GetUserProfile возвращает профиль пользователя с базовой информацией
func GetUserProfile(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Быстрая базовая статистика
	var totalActions int64
	database.DB.Model(&models.Action{}).Where("user_id = ?", userId).Count(&totalActions)

	var achievementsCount int64
	database.DB.Model(&models.UserAchievement{}).Where("user_id = ? AND is_achieved = ?", userId, true).Count(&achievementsCount)

	response := gin.H{
		"id":            user.ID,
		"email":         user.Email,
		"name":          user.Username,
		"points":        user.Points,
		"ghg_index":     user.GHGIndex,
		"total_actions": totalActions,
		"achievements":  achievementsCount,
		"latitude":      user.Latitude,
		"longitude":     user.Longitude,
		"created_at":    user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
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
