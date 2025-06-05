package User

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/genryusaishigikuni/ch4nge/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Preload("Friends").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, utils.UserToResponse(user))
	}

	c.JSON(http.StatusOK, userResponses)
}
