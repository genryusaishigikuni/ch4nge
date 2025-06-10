package admin

import (
	"log"
	"net/http"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func CreateAchievement(c *gin.Context) {
	// Log the incoming request
	log.Println("Received request to create new achievement")

	// Bind the incoming JSON request to the Achievement struct
	var achievement models.Achievement
	if err := c.ShouldBindJSON(&achievement); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the creation attempt
	log.Printf("Creating achievement: %s", achievement.Title)

	// Create the achievement
	if err := db.DB.Create(&achievement).Error; err != nil {
		log.Printf("Error creating achievement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create achievement"})
		return
	}

	// Log successful creation
	log.Printf("Achievement %s created successfully, proceeding to assign to all users", achievement.Title)

	// Assign the achievement to all users
	if err := assignAchievementToAllUsers(achievement.ID); err != nil {
		log.Printf("Failed to assign achievement %d to all users: %v", achievement.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Achievement created but failed to assign to all users"})
		return
	}

	// Log the success of assignment
	log.Printf("Achievement %d successfully assigned to all users", achievement.ID)

	// Respond with the created achievement
	c.JSON(http.StatusCreated, achievement)
}

func assignAchievementToAllUsers(achievementID uint) error {
	// Log the assignment process
	log.Printf("Assigning achievement %d to all users", achievementID)

	var users []models.User
	// Fetch all users
	if err := db.DB.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return err
	}

	// Loop through each user and assign the achievement
	for _, user := range users {
		userAchievement := models.UserAchievement{
			UserID:        user.ID,
			AchievementID: achievementID,
		}
		// Attempt to assign the achievement
		if err := db.DB.FirstOrCreate(&userAchievement, models.UserAchievement{
			UserID:        user.ID,
			AchievementID: achievementID,
		}).Error; err != nil {
			log.Printf("Error assigning achievement to user %d: %v", user.ID, err)
			return err
		}
	}

	// Log completion of the assignment
	log.Printf("Successfully assigned achievement %d to all users", achievementID)
	return nil
}
