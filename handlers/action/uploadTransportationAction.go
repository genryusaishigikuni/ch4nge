package action

import (
	"net/http"

	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CalculateGHGIndex(distance, fuelConsumption float64, passengers int) float64 {
	if passengers <= 0 {
		passengers = 1
	}
	co2PerKm := (fuelConsumption / 100.0) * 2.3
	totalCO2 := co2PerKm * distance
	return totalCO2 / float64(passengers)
}

func CalculateGHGPoints(distance float64, fuelConsumption float64, passengers int) int {
	if passengers <= 0 {
		passengers = 1
	}

	co2PerKm := (fuelConsumption / 100.0) * 2.3
	totalCO2 := co2PerKm * distance
	co2PerPerson := totalCO2 / float64(passengers)

	switch {
	case co2PerPerson > 40:
		return -10
	case co2PerPerson > 25:
		return -5
	case co2PerPerson > 15:
		return 0
	case co2PerPerson > 5:
		return 5
	default:
		return 10
	}
}

func UploadTransportationAction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.TransportationActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	distance, _ := req.Payload["distance"].(float64)
	fuelConsumption, _ := req.Payload["fuelConsumption"].(float64)
	passengers, _ := req.Payload["passengers"].(float64)

	points := CalculateGHGPoints(distance, fuelConsumption, int(passengers))
	ghg := CalculateGHGIndex(distance, fuelConsumption, int(passengers))

	action := models.Action{
		UserID:     userID.(uint),
		ActionType: req.ActionType,
		Payload:    req.Payload,
		Metadata:   req.Metadata,
		Points:     points,
	}

	if err := database.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload transportation action"})
		return
	}

	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", action.Points))
	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("ghg_index", gorm.Expr("ghg_index + ?", ghg))

	// Optional: update user location
	if locationArray, ok := req.Payload["location"].([]interface{}); ok && len(locationArray) == 2 {
		if latitude, ok := locationArray[0].(float64); ok {
			if longitude, ok := locationArray[1].(float64); ok {
				database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
					"latitude":  latitude,
					"longitude": longitude,
				})
			}
		}
	}

	actionTitle := "Used transportation"
	if option, ok := req.Payload["option"].(string); ok && option != "" {
		if vehicle, ok := req.Payload["vehicle"].(string); ok && vehicle != "" {
			actionTitle = formatTransportationActionTitle(option, vehicle)
		} else {
			actionTitle = option
		}
	}

	activity := models.Activity{
		UserID: userID.(uint),
		Title:  actionTitle,
		Value:  points,
	}
	database.DB.Create(&activity)

	c.JSON(http.StatusCreated, gin.H{"message": "Transportation action uploaded successfully."})
}
