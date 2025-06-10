package action

import (
	"fmt"
	"math"
	"strings"
	"time"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// CalculateTransportationImpact Улучшенная функция расчета GHG и поинтов
func CalculateTransportationImpact(distance, fuelConsumption float64, passengers int, transportType string) (points int, ghg float64, isEcoFriendly bool) {
	if passengers <= 0 {
		passengers = 1
	}

	// Базовый расчет CO2 (кг CO2 на км)
	var co2PerKm float64

	switch transportType {
	case "bicycle", "walking", "scooter":
		co2PerKm = 0 // Экологичный транспорт
		isEcoFriendly = true
	case "public_transport", "bus", "metro", "train":
		co2PerKm = 0.05 // Общественный транспорт
		isEcoFriendly = true
	case "electric_car":
		co2PerKm = 0.1 // Электромобиль
		isEcoFriendly = true
	case "car", "private_vehicle":
		co2PerKm = (fuelConsumption / 100.0) * 2.3 // Личный автомобиль
		isEcoFriendly = false
	case "motorcycle":
		co2PerKm = (fuelConsumption / 100.0) * 2.1
		isEcoFriendly = false
	default:
		co2PerKm = (fuelConsumption / 100.0) * 2.3
		isEcoFriendly = false
	}

	totalCO2 := co2PerKm * distance
	ghg = totalCO2 / float64(passengers)

	// Расчет поинтов на основе экологичности
	if isEcoFriendly {
		// Позитивные поинты за экологичный транспорт
		basePoints := int(distance * 0.5) // 0.5 поинта за км
		if basePoints > 50 {
			basePoints = 50
		}
		if basePoints < 5 {
			basePoints = 5
		}
		points = basePoints
	} else {
		// Негативные поинты за неэкологичный транспорт
		co2PerPerson := ghg
		switch {
		case co2PerPerson > 50:
			points = -15 // Очень плохо
		case co2PerPerson > 30:
			points = -10 // Плохо
		case co2PerPerson > 15:
			points = -5 // Не очень
		case co2PerPerson > 5:
			points = 0 // Нейтрально
		default:
			points = 2 // Немного хорошо
		}
	}

	return points, ghg, isEcoFriendly
}

// UploadTransportationAction Улучшенная функция загрузки транспортного действия
func UploadTransportationAction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.TransportationActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the transportation data from the request
	distance, _ := req.Payload["distance"].(float64)
	fuelConsumption, _ := req.Payload["fuelConsumption"].(float64)
	passengers, _ := req.Payload["passengers"].(float64)
	transportType, _ := req.Payload["transportType"].(string)
	vehicle, _ := req.Payload["vehicle"].(string)

	// Normalize transport type (handle variations like "e-scooter")
	transportType = strings.ToLower(strings.ReplaceAll(transportType, "-", "_"))

	// Calculate the transportation impact (GHG, points, eco-friendliness)
	points, ghg, isEcoFriendly := CalculateTransportationImpact(distance, fuelConsumption, int(passengers), transportType)

	// **New check** for First Eco-Friendly Action
	var userActionsCount int64
	db.DB.Model(&models.Action{}).Where("user_id = ?", userID).Count(&userActionsCount)

	if userActionsCount == 1 && isEcoFriendly { // First eco-friendly action
		// Trigger "First Green Action" achievement
		var firstActionAchievement models.Achievement
		if err := db.DB.Where("title = ?", "First Green Action").First(&firstActionAchievement).Error; err == nil {
			userAchievement := models.UserAchievement{
				UserID:        userID.(uint),
				AchievementID: firstActionAchievement.ID,
			}

			// Create achievement only if it hasn't been awarded already
			db.DB.FirstOrCreate(&userAchievement, models.UserAchievement{
				UserID:        userID.(uint),
				AchievementID: firstActionAchievement.ID,
			})

			// Mark the achievement as achieved
			userAchievement.IsAchieved = true
			achievedAt := time.Now()
			userAchievement.AchievedAt = &achievedAt
			db.DB.Save(&userAchievement)
		}
	}

	// Continue with your normal action saving logic
	action := models.Action{
		UserID:     userID.(uint),
		ActionType: req.ActionType,
		Payload:    req.Payload,
		Metadata:   req.Metadata,
		Points:     points,
	}

	// Save action to the database
	if err := db.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload transportation action"})
		return
	}

	// Update user's points and GHG index
	db.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", points))
	db.DB.Model(&models.User{}).Where("id = ?", userID).Update("ghg_index", gorm.Expr("ghg_index + ?", ghg))

	// Handle location update if available
	if locationArray, ok := req.Payload["location"].([]interface{}); ok && len(locationArray) == 2 {
		if latitude, ok := locationArray[0].(float64); ok {
			if longitude, ok := locationArray[1].(float64); ok {
				db.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
					"latitude":  latitude,
					"longitude": longitude,
				})
			}
		}
	}

	// Format action title based on transport type
	actionTitle := formatTransportationActionTitle(transportType, vehicle, distance, isEcoFriendly)

	// Create activity record
	activity := models.Activity{
		UserID: userID.(uint),
		Title:  actionTitle,
		Value:  points,
	}
	db.DB.Create(&activity)

	// Check achievements and challenges asynchronously
	go checkAchievementsAndChallenges(userID.(uint), "transportation", distance, float64(points), isEcoFriendly)

	// Response
	response := gin.H{
		"message":      "Transportation action uploaded successfully",
		"points":       points,
		"ghg_impact":   ghg,
		"eco_friendly": isEcoFriendly,
	}

	c.JSON(http.StatusCreated, response)
}

func UploadGreenAction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Check if it's the first green action for the user
	var userActionsCount int64
	db.DB.Model(&models.Action{}).Where("user_id = ?", userID).Count(&userActionsCount)

	// Check if it’s the first green action (the user has just performed their first green action)
	if userActionsCount == 1 { // First action
		var firstActionAchievement models.Achievement
		// Ensure we only award the "First Green Action" achievement
		if err := db.DB.Where("title = ?", "First Green Action").First(&firstActionAchievement).Error; err == nil {
			var userAchievement models.UserAchievement
			// Check if the user already has the achievement
			db.DB.Where("user_id = ? AND achievement_id = ?", userID.(uint), firstActionAchievement.ID).First(&userAchievement)

			// Create the achievement only if it hasn't been awarded yet
			if userAchievement.ID == 0 { // Achievement not found
				userAchievement = models.UserAchievement{
					UserID:        userID.(uint),
					AchievementID: firstActionAchievement.ID,
					IsAchieved:    true,
				}

				// Mark as achieved
				achievedAt := time.Now()
				userAchievement.AchievedAt = &achievedAt
				db.DB.Create(&userAchievement) // Save to DB
			}
		}
	}

	var req models.GreenActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate points for the green action
	points := calculateGreenActionPoints(req.Payload)

	action := models.Action{
		UserID:     userID.(uint),
		ActionType: req.ActionType,
		Payload:    req.Payload,
		Metadata:   req.Metadata,
		Points:     points,
	}

	// Save the action to the database
	if err := db.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload green action"})
		return
	}

	// Update user's points after the green action
	db.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", points))

	// Update user location if provided in the payload
	if locationArray, ok := req.Payload["location"].([]interface{}); ok && len(locationArray) == 2 {
		if latitude, ok := locationArray[0].(float64); ok {
			if longitude, ok := locationArray[1].(float64); ok {
				db.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
					"latitude":  latitude,
					"longitude": longitude,
				})
			}
		}
	}

	// Format the title of the green action based on the option
	option, _ := req.Payload["option"].(string)
	actionTitle := formatGreenActionTitle(option)

	// Create activity record based on the green action
	activity := models.Activity{
		UserID: userID.(uint),
		Title:  actionTitle,
		Value:  points,
	}
	db.DB.Create(&activity)

	// Check achievements and challenges asynchronously
	go checkAchievementsAndChallenges(userID.(uint), "green", 1.0, float64(points), true)

	// Send success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Green action uploaded successfully",
		"points":  points,
	})
}

// Функция расчета поинтов для зеленых действий
func calculateGreenActionPoints(payload map[string]interface{}) int {
	option, _ := payload["option"].(string)

	switch option {
	case "planted_tree":
		return 50 // Посадка дерева - максимальные поинты
	case "solar_power":
		return 30 // Солнечная энергия - высокие поинты
	case "composting":
		return 25 // Компостирование - высокие поинты
	case "recycling":
		return 20 // Переработка - хорошие поинты
	case "water_conservation":
		return 15 // Сохранение воды - средние поинты
	case "energy_saving":
		return 15 // Энергосбережение - средние поинты
	case "waste_reduction":
		return 15 // Сокращение отходов - средние поинты
	case "used_bike":
		return 20 // Использование велосипеда - хорошие поинты
	case "public_transport":
		return 15 // Общественный транспорт - средние поинты
	case "lights_off":
		return 10 // Выключение света - базовые поинты
	default:
		return 10 // По умолчанию
	}
}

// Улучшенная функция форматирования заголовков
func formatTransportationActionTitle(transportType, vehicle string, distance float64, isEcoFriendly bool) string {
	ecoIcon := ""
	if isEcoFriendly {
		ecoIcon = "🌱 "
	} else {
		ecoIcon = "🚗 "
	}

	distanceStr := fmt.Sprintf("%.1f km", distance)

	switch transportType {
	case "bicycle":
		return fmt.Sprintf("%sBiked %s", ecoIcon, distanceStr)
	case "walking":
		return fmt.Sprintf("%sWalked %s", ecoIcon, distanceStr)
	case "public_transport":
		return fmt.Sprintf("%sUsed public transport for %s", ecoIcon, distanceStr)
	case "electric_car":
		return fmt.Sprintf("%sDrove electric car %s", ecoIcon, distanceStr)
	default:
		if vehicle != "" {
			return fmt.Sprintf("%sUsed %s for %s", ecoIcon, vehicle, distanceStr)
		}
		return fmt.Sprintf("%sTransportation: %s", ecoIcon, distanceStr)
	}
}

// Асинхронная проверка достижений и челленджей
func checkAchievementsAndChallenges(userID uint, actionType string, value, points float64, isEcoFriendly bool) {
	// Проверяем достижения
	checkUserAchievements(userID, actionType, value, points, isEcoFriendly)

	// Обновляем прогресс weekly challenge
	updateWeeklyChallengeProgress(userID, actionType, value, points, isEcoFriendly)
}

// Проверка достижений пользователя
func checkUserAchievements(userID uint, actionType string, value, points float64, isEcoFriendly bool) {
	var userAchievements []models.UserAchievement
	db.DB.Preload("Achievement").Where("user_id = ? AND is_achieved = ?", userID, false).Find(&userAchievements)

	for _, ua := range userAchievements {
		shouldAchieve := false

		// Простая логика проверки достижений
		switch ua.Achievement.Title {
		case "First Green Action":
			if actionType == "green" || (actionType == "transportation" && isEcoFriendly) {
				shouldAchieve = true
			}

		case "Eco Traveler":
			if actionType == "transportation" && isEcoFriendly {
				shouldAchieve = true
			}
		case "Point Collector":
			// Проверяем общее количество поинтов пользователя
			var user models.User
			if db.DB.First(&user, userID).Error == nil && user.Points >= 100 {
				shouldAchieve = true
			}
		case "Distance Master":
			if actionType == "transportation" && value >= 50.0 {
				shouldAchieve = true
			}
		case "Green Warrior":
			// Проверяем количество зеленых действий
			var greenActionCount int64
			db.DB.Model(&models.Action{}).Where("user_id = ? AND action_type = ?", userID, "green").Count(&greenActionCount)
			if greenActionCount >= 10 {
				shouldAchieve = true
			}
		}

		if shouldAchieve {
			ua.IsAchieved = true
			achievedAt := time.Now()
			ua.AchievedAt = &achievedAt
			db.DB.Save(&ua)

			// Добавляем поинты за достижение
			if ua.Achievement.Points > 0 {
				db.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", ua.Achievement.Points))
			}
		}
	}
}

// Update the progress of the user's weekly challenge.
func updateWeeklyChallengeProgress(userID uint, actionType string, value, points float64, isEcoFriendly bool) {
	var userChallenge models.UserWeeklyChallenge
	if db.DB.Preload("WeeklyChallenge").Where("user_id = ? AND is_completed = ?", userID, false).First(&userChallenge).Error != nil {
		return // No active challenge
	}

	// Initialize the progress increment
	var progressIncrement float64 = 0

	// Add the respective progress based on the challenge title
	switch userChallenge.WeeklyChallenge.Title {
	case "Eco Transport Week":
		if actionType == "transportation" && isEcoFriendly {
			progressIncrement = value // Add traveled distance
		}
	case "Green Actions Week":
		if actionType == "green" {
			progressIncrement = 1 // Add 1 action for each green action
		}
	case "Point Challenge":
		if points > 0 {
			progressIncrement = points // Add points to the challenge
		}
	case "Carbon Reduction":
		if isEcoFriendly {
			progressIncrement = math.Max(value*0.1, 1) // Add carbon reduction progress
		}
	}

	// If there's progress, update the challenge
	if progressIncrement > 0 {
		userChallenge.CurrentValue += progressIncrement

		// Check if the challenge is completed
		if userChallenge.CurrentValue >= userChallenge.WeeklyChallenge.TargetValue && !userChallenge.IsCompleted {
			userChallenge.IsCompleted = true
			completedAt := time.Now()
			userChallenge.CompletedAt = &completedAt

			// Add points for completing the challenge
			if userChallenge.WeeklyChallenge.Points > 0 {
				db.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", userChallenge.WeeklyChallenge.Points))
			}
		}

		db.DB.Save(&userChallenge) // Save the updated challenge progress
	}
}
