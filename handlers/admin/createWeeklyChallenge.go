package admin

import (
	"log"
	"net/http"
	"time"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func CreateWeeklyChallenge(c *gin.Context) {
	// Log the incoming request
	log.Println("Received request to create a new weekly challenge")

	// Bind the incoming JSON request to the WeeklyChallenge struct
	var req models.WeeklyChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the challenge creation attempt
	log.Printf("Creating weekly challenge: %s", req.Title)

	// Create the weekly challenge
	challenge := models.WeeklyChallenge{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Points:      req.Points,
		TargetValue: req.TargetValue,
		IsActive:    req.IsActive,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := db.DB.Create(&challenge).Error; err != nil {
		log.Printf("Error creating weekly challenge: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create weekly challenge"})
		return
	}

	// Log successful creation
	log.Printf("Weekly challenge %s created successfully, proceeding to assign to all users", challenge.Title)

	// Assign the weekly challenge to all users
	if err := assignWeeklyChallengeToAllUsers(challenge.ID); err != nil {
		log.Printf("Failed to assign weekly challenge %d to all users: %v", challenge.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Weekly challenge created but failed to assign to all users"})
		return
	}

	// Log the success of assignment
	log.Printf("Weekly challenge %d successfully assigned to all users", challenge.ID)

	// Respond with the created weekly challenge
	c.JSON(http.StatusCreated, challenge)
}

func assignWeeklyChallengeToAllUsers(challengeID uint) error {
	// Log the assignment process
	log.Printf("Assigning weekly challenge %d to all users", challengeID)

	var users []models.User
	// Fetch all users
	if err := db.DB.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return err
	}

	// Loop through each user and assign the weekly challenge
	for _, user := range users {
		userChallenge := models.UserWeeklyChallenge{
			UserID:            user.ID,
			WeeklyChallengeID: challengeID,
			CurrentValue:      0.0,
			AssignedAt:        time.Now(),
		}
		// Attempt to assign the weekly challenge
		if err := db.DB.FirstOrCreate(&userChallenge, models.UserWeeklyChallenge{
			UserID:            user.ID,
			WeeklyChallengeID: challengeID,
		}).Error; err != nil {
			log.Printf("Error assigning weekly challenge to user %d: %v", user.ID, err)
			return err
		}
	}

	// Log completion of the assignment
	log.Printf("Successfully assigned weekly challenge %d to all users", challengeID)
	return nil
}
