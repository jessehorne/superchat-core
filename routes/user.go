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
		GivenFields: models.GivenFields{
			ID: uuid.New().String(),
		},
		Email:        req.Email,
		Password:     hash,
		PasswordSalt: salt,
	}

	result := database.GDB.Create(&u)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":    "couldn't create user",
			"errors": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userID": u.ID,
	})
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
	sessionResult := database.GDB.First(&sesh, "user_id = ?", user.ID)
	sesh.Token = hash
	sesh.ExpiresAt = time.Now().Local().Add(1 * time.Hour)
	sesh.UserID = user.ID

	if sessionResult.RowsAffected == 0 {
		// we need to creat the record now
		sesh.ID = uuid.New().String()
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

type UserUpdateRequest struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

func UserUpdate(c *gin.Context) {
	var req UserUpdateRequest
	err, res := util.TryBind(&req, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if req.UserID == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing details",
		})
		return
	}

	// get user from request
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no auth user",
		})
		return
	}

	user := u.(models.User)

	if user.ID != req.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you can't do that dude",
		})
		return
	}

	salt, hash := util.ProcessPassword(req.Password)
	user.PasswordSalt = salt
	user.Password = hash

	saveResult := database.GDB.Save(&user)
	if saveResult.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db issue while saving user",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

type UserDeleteRequest struct {
	UserID string `json:"userID"`
}

func UserDelete(c *gin.Context) {
	var req UserUpdateRequest
	err, res := util.TryBind(&req, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing details",
		})
		return
	}

	// get user from request
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no auth user",
		})
		return
	}

	user := u.(models.User)

	if user.ID != req.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you can't do that dude",
		})
		return
	}

	deleteResult := database.GDB.Delete(&user)
	if deleteResult.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db issue while delete user",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
