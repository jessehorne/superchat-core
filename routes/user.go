package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jessehorne/superchat-core/database"
	"github.com/jessehorne/superchat-core/database/models"
	"github.com/jessehorne/superchat-core/util"
	"net/http"
	"time"
)

type UserCreateRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password,min=8,max=255"`
}

func UserCreate(c *gin.Context) {
	var req UserCreateRequest
	err, res := util.TryBind(&req, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, res)
		return
	}

	salt, hash := util.ProcessPassword(req.Password)

	// attempt to create user
	u := models.User{
		Email:        req.Email,
		Password:     hash,
		PasswordSalt: salt,
	}
	u.ID = uuid.New().String()

	result := database.GDB.Create(&u)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":    "couldn't create user",
			"errors": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

type UserGetTokenRequest struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password,gte=8,lte=255"`
}

func UserGetToken(c *gin.Context) {
	var req UserGetTokenRequest
	err, res := util.TryBind(&req, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// get user by email
	var user models.User
	result := database.GDB.First(&user, "email = ?", req.Email)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "no user found",
		})
		return
	}

	// validate password
	if !util.ComparePassword(req.Password, user.PasswordSalt, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid credentials",
		})
		return
	}

	// user is GOOD TO GO!
	// generate token to send to user and store hashed token in database
	token, hash := util.CreateToken()

	// check if session already exists...if so, update it
	var sesh models.Session
	sessionResult := database.GDB.First(&sesh, "id = ?", user.ID)

	sesh.Token = hash
	sesh.ExpiresAt = time.Now().Local().Add(1 * time.Hour)
	sesh.UserID = user.ID

	if sessionResult.RowsAffected == 0 {
		// we need to creat the record now
		database.GDB.Create(&sesh)
	} else {
		// we need to update it
		database.GDB.Save(&sesh)
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"userID":    user.ID,
		"expiresAt": sesh.ExpiresAt.Format(time.RFC3339),
	})
}
