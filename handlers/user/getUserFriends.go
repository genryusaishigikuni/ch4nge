package User

import (
	db "github.com/genryusaishigikuni/ch4nge/database"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/genryusaishigikuni/ch4nge/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserFriends(c *gin.Context) {
	userID := c.Param("userId")

	var user models.User
	if err := db.DB.Preload("Friends").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	friendResponses := make([]models.UserResponse, 0)
	for _, friend := range user.Friends {
		friendResponses = append(friendResponses, utils.UserToResponse(*friend))
	}

	c.JSON(http.StatusOK, friendResponses)
}
