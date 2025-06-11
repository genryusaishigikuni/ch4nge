package action

import (
	"fmt"
	"math"

	"log"
	"strings"
	"time"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func CalculateTransportationImpact(distance, fuelConsumption float64, passengers int, transportType string) (points int, ghg float64, isEcoFriendly bool) {
	if passengers <= 0 {
		passengers = 1
	}

	var co2PerKm float64
	log.Printf("Calculating CO2 impact for transport type: %s", transportType)

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
	case "car", "private_vehicle":
		co2PerKm = (fuelConsumption / 100.0) * 2.3
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

	log.Printf("Total CO2 impact: %.2f, GHG per person: %.2f", totalCO2, ghg)

	if isEcoFriendly {
		basePoints := int(distance * 0.5)
		if basePoints > 50 {
			basePoints = 50
		}
		if basePoints < 5 {
			basePoints = 5
		}
		points = basePoints
	} else {
		co2PerPerson := ghg
		switch {
		case co2PerPerson > 50:
			points = -15
		case co2PerPerson > 30:
			points = -10
		case co2PerPerson > 15:
			points = -5
		case co2PerPerson > 5:
			points = 0
		default:
			points = 2
		}
	}

	return points, ghg, isEcoFriendly
}

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

	// FIX 1: Use actionType as transportType if transportType is not provided in payload
	transportType, ok := req.Payload["transportType"].(string)
	if !ok || transportType == "" {
		transportType = req.ActionType // Use actionType (e.g., "bicycle") as transportType
		log.Printf("Using actionType '%s' as transportType", transportType)
	}

	vehicle, _ := req.Payload["vehicle"].(string)

	transportType = strings.ToLower(strings.ReplaceAll(transportType, "-", "_"))

	// FIX 2: Check if client provided isEcoFriendly in metadata and respect it
	clientEcoFriendly, hasClientEcoFriendly := req.Metadata["isEcoFriendly"].(bool)

	points, ghg, calculatedEcoFriendly := CalculateTransportationImpact(distance, fuelConsumption, int(passengers), transportType)

	// Use client-provided eco-friendly flag if available, otherwise use calculated value
	var isEcoFriendly bool
	if hasClientEcoFriendly {
		isEcoFriendly = clientEcoFriendly
		log.Printf("Using client-provided eco-friendly flag: %v", isEcoFriendly)
	} else {
		isEcoFriendly = calculatedEcoFriendly
		log.Printf("Using calculated eco-friendly flag: %v", isEcoFriendly)
	}

	var userActionsCount int64
	db.DB.Model(&models.Action{}).Where("user_id = ?", userID).Count(&userActionsCount)
	log.Printf("User %d has performed %d actions.", userID, userActionsCount)

	if userActionsCount == 1 && isEcoFriendly { // First eco-friendly action
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

			userAchievement.IsAchieved = true
			achievedAt := time.Now()
			userAchievement.AchievedAt = &achievedAt
			db.DB.Save(&userAchievement)

			log.Printf("User %d achieved 'First Green Action'.", userID)
		}
	}

	// Save action to the database
	action := models.Action{
		UserID:     userID.(uint),
		ActionType: req.ActionType,
		Payload:    req.Payload,
		Metadata:   req.Metadata,
		Points:     points,
	}

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
				log.Printf("User %d location updated: %.6f, %.6f", userID, latitude, longitude)
			}
		}
	}

	// Format action title
	actionTitle := formatTransportationActionTitle(transportType, vehicle, distance, isEcoFriendly)

	// Create activity record
	activity := models.Activity{
		UserID: userID.(uint),
		Title:  actionTitle,
		Value:  points,
	}
	db.DB.Create(&activity)

	// FIX: Update weekly challenge progress synchronously ONLY ONCE
	updateWeeklyChallengeProgress(userID.(uint), "transportation", distance, float64(points), isEcoFriendly)

	// Check achievements asynchronously (without updating weekly challenge again)
	go checkUserAchievements(userID.(uint), "transportation", distance, float64(points), isEcoFriendly)

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
	log.Printf("User %d has performed %d green actions.", userID, userActionsCount)

	if userActionsCount == 1 { // First action
		var firstActionAchievement models.Achievement
		if err := db.DB.Where("title = ?", "First Green Action").First(&firstActionAchievement).Error; err == nil {
			var userAchievement models.UserAchievement
			db.DB.Where("user_id = ? AND achievement_id = ?", userID.(uint), firstActionAchievement.ID).First(&userAchievement)

			if userAchievement.ID == 0 { // Achievement not found
				userAchievement = models.UserAchievement{
					UserID:        userID.(uint),
					AchievementID: firstActionAchievement.ID,
					IsAchieved:    true,
				}

				achievedAt := time.Now()
				userAchievement.AchievedAt = &achievedAt
				db.DB.Create(&userAchievement)
				log.Printf("User %d achieved 'First Green Action'.", userID)
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

	if err := db.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload green action"})
		return
	}

	db.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", points))

	// Update location if available
	if locationArray, ok := req.Payload["location"].([]interface{}); ok && len(locationArray) == 2 {
		if latitude, ok := locationArray[0].(float64); ok {
			if longitude, ok := locationArray[1].(float64); ok {
				db.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
					"latitude":  latitude,
					"longitude": longitude,
				})
				log.Printf("User %d location updated: %.6f, %.6f", userID, latitude, longitude)
			}
		}
	}

	// Format action title
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

// Function to calculate points for green actions
func calculateGreenActionPoints(payload map[string]interface{}) int {
	option, _ := payload["option"].(string)

	log.Printf("Calculating points for green action: %s", option)

	switch option {
	case "planted_tree":
		return 50
	case "solar_power":
		return 30
	case "composting":
		return 25
	case "recycling":
		return 20
	case "water_conservation":
		return 15
	case "energy_saving":
		return 15
	case "waste_reduction":
		return 15
	case "used_bike":
		return 20
	case "public_transport":
		return 15
	case "lights_off":
		return 10
	default:
		return 10
	}
}

// Function to format transportation action title
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

// For transportation actions, use checkUserAchievements directly to avoid double counting
func checkAchievementsAndChallenges(userID uint, actionType string, value, points float64, isEcoFriendly bool) {
	log.Printf("Checking achievements and challenges for user %d", userID)
	// Check achievements
	checkUserAchievements(userID, actionType, value, points, isEcoFriendly)

	// Update progress of weekly challenge
	updateWeeklyChallengeProgress(userID, actionType, value, points, isEcoFriendly)
}
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
// Update the progress of the user's weekly challenge.
func updateWeeklyChallengeProgress(userID uint, actionType string, value, points float64, isEcoFriendly bool) {
	var userChallenge models.UserWeeklyChallenge
	if err := db.DB.Preload("WeeklyChallenge").Where("user_id = ? AND is_completed = ?", userID, false).First(&userChallenge).Error; err != nil {
		log.Printf("No active weekly challenge found for user %d: %v", userID, err)
		return // No active challenge
	}

	log.Printf("Found active challenge '%s' for user %d. Current progress: %.2f/%.2f",
		userChallenge.WeeklyChallenge.Title, userID, userChallenge.CurrentValue, userChallenge.WeeklyChallenge.TargetValue)

	// Initialize the progress increment
	var progressIncrement float64 = 0

	// Add the respective progress based on the challenge title
	switch userChallenge.WeeklyChallenge.Title {
	case "Eco Transport Week":
		if actionType == "transportation" && isEcoFriendly {
			progressIncrement = value // Add traveled distance
			log.Printf("Adding transportation distance: %.2f", progressIncrement)
		}
	case "Green Actions Week":
		if actionType == "green" {
			progressIncrement = 1 // Add 1 action for each green action
			log.Printf("Adding green action count: %.2f", progressIncrement)
		}
	case "Point Challenge":
		if points > 0 {
			progressIncrement = points // Add points to the challenge
			log.Printf("Adding points: %.2f", progressIncrement)
		}
	case "Carbon Reduction":
		if isEcoFriendly {
			progressIncrement = math.Max(value*0.1, 1) // Add carbon reduction progress
			log.Printf("Adding carbon reduction progress: %.2f", progressIncrement)
		}
	}

	// If there's progress, update the challenge
	if progressIncrement > 0 {
		userChallenge.CurrentValue += progressIncrement
		log.Printf("Updated progress for user %d: %.2f -> %.2f (added %.2f)",
			userID, userChallenge.CurrentValue-progressIncrement, userChallenge.CurrentValue, progressIncrement)

		// Check if the challenge is completed
		if userChallenge.CurrentValue >= userChallenge.WeeklyChallenge.TargetValue && !userChallenge.IsCompleted {
			userChallenge.IsCompleted = true
			completedAt := time.Now()
			userChallenge.CompletedAt = &completedAt
			log.Printf("Challenge '%s' completed for user %d!", userChallenge.WeeklyChallenge.Title, userID)

			// Add points for completing the challenge
			if userChallenge.WeeklyChallenge.Points > 0 {
				if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", userChallenge.WeeklyChallenge.Points)).Error; err != nil {
					log.Printf("Error updating user points for completing challenge: %v", err)
				} else {
					log.Printf("Added %d points to user %d for completing challenge", userChallenge.WeeklyChallenge.Points, userID)
				}
			}
		}

		// THIS IS THE CRITICAL PART THAT WAS MISSING - SAVE THE UPDATED PROGRESS
		if err := db.DB.Save(&userChallenge).Error; err != nil {
			log.Printf("Error saving weekly challenge progress for user %d: %v", userID, err)
		} else {
			log.Printf("Successfully saved weekly challenge progress for user %d", userID)
		}
	} else {
		log.Printf("No progress increment for challenge '%s' with action type '%s', eco-friendly: %v",
			userChallenge.WeeklyChallenge.Title, actionType, isEcoFriendly)
	}
}
