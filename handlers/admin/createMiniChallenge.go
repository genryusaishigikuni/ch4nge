package admin

import (
	"log"
	"net/http"
	"time"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
)

func CreateMiniChallenge(c *gin.Context) {
	// Log the incoming request
	log.Println("Received request to create new mini challenge")

	// Bind the incoming JSON request to the MiniChallenge struct
	var challenge models.MiniChallenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the creation attempt
	log.Printf("Creating mini challenge: %s", challenge.Title)

	// Create the mini challenge
	if err := db.DB.Create(&challenge).Error; err != nil {
		log.Printf("Error creating mini challenge: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mini challenge"})
		return
	}

	// Log successful creation
	log.Printf("Mini challenge %s created successfully, proceeding to assign to all users", challenge.Title)

	// Assign the mini challenge to all users
	if err := assignMiniChallengeToAllUsers(challenge.ID); err != nil {
		log.Printf("Failed to assign mini challenge %d to all users: %v", challenge.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Mini challenge created but failed to assign to all users"})
		return
	}

	// Log the success of assignment
	log.Printf("Mini challenge %d successfully assigned to all users", challenge.ID)

	// Respond with the created mini challenge
	c.JSON(http.StatusCreated, challenge)
}

func assignMiniChallengeToAllUsers(challengeID uint) error {
	// Log the assignment process
	log.Printf("Assigning mini challenge %d to all users", challengeID)

	var users []models.User
	// Fetch all users
	if err := db.DB.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return err
	}

	// Loop through each user and assign the mini challenge
	for _, user := range users {
		userChallenge := models.UserMiniChallenge{
			UserID:          user.ID,
			MiniChallengeID: challengeID,
			AssignedAt:      time.Now(),
		}
		// Attempt to assign the mini challenge
		if err := db.DB.FirstOrCreate(&userChallenge, models.UserMiniChallenge{
			UserID:          user.ID,
			MiniChallengeID: challengeID,
		}).Error; err != nil {
			log.Printf("Error assigning mini challenge to user %d: %v", user.ID, err)
			return err
		}
	}

	// Log completion of the assignment
	log.Printf("Successfully assigned mini challenge %d to all users", challengeID)
	return nil
}
