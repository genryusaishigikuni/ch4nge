package utils

import (
	"github.com/genryusaishigikuni/ch4nge/models"
	"strconv"
)

func UserToResponse(user models.User) models.UserResponse {
	var friendIds []string
	for _, friend := range user.Friends {
		friendIds = append(friendIds, strconv.Itoa(int(friend.ID)))
	}

	return models.UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		ProfilePicURL: user.ProfilePicURL,
		Streak:        user.Streak,
		Points:        user.Points,
		GHGIndex:      user.GHGIndex,
		Location:      []float64{user.Latitude, user.Longitude},
		FriendsIds:    friendIds,
	}
}

func ParseUint(s string) uint {
	if id, err := strconv.Atoi(s); err == nil {
		return uint(id)
	}
	return 0
}
