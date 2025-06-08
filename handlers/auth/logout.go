package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
