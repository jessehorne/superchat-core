package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jessehorne/superchat-core/database"
	"github.com/jessehorne/superchat-core/database/models"
	"github.com/jessehorne/superchat-core/util"
	"net/http"
	"time"
)

func AuthMiddleware(c *gin.Context) {
	userID := c.GetHeader("userID")
	token := c.GetHeader("Authorization")

	if userID == "" || token == "" {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	// check if session exists with that userID
	var sesh models.Session
	seshResult := database.GDB.First(&sesh, "user_id = ?", userID)
	if seshResult.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	// check if token is valid
	if !util.ValidateToken(token, sesh.Token) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		c.Abort()
		return
	}

	// check if expires at isn't past
	if time.Now().After(sesh.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "expired token",
		})
		c.Abort()
		return
	}

	// get user by userID
	var user models.User
	userResult := database.GDB.First(&user, "id = ?", userID)
	if userResult.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}
