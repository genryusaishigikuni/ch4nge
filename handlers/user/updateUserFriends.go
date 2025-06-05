package User

import (
	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/utils"

	"github.com/genryusaishigikuni/ch4nge/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UpdateUserFriends(c *gin.Context) {
	userID := c.Param("userId")

	var req models.UpdateFriendsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Clear existing friendships
	db.DB.Where("user_id = ?", userID).Delete(&models.UserFriend{})

	// Add new friendships
	for _, friendIDStr := range req.FriendIds {
		friendID, err := strconv.Atoi(friendIDStr)
		if err != nil {
			continue
		}

		db.DB.Create(&models.UserFriend{
			UserID:   utils.ParseUint(userID),
			FriendID: uint(friendID),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friends updated successfully"})
}
