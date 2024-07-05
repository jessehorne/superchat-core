package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jessehorne/superchat-core/database"
	"github.com/jessehorne/superchat-core/database/models"
	"github.com/jessehorne/superchat-core/util"
	"net/http"
	"time"
)

type UserCreateRequest struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password,gte=8,lte=255"`
}

func UserCreate(c *gin.Context) {
	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := err.(validator.ValidationErrors)
		errs := []string{}
		for _, e := range errors {
			errs = append(errs, e.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":    "invalid",
			"errors": errs,
		})
		return
	}

	salt, hash := util.ProcessPassword(req.Password)

	// attempt to create user
	u := models.User{
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

	c.JSON(http.StatusOK, gin.H{})
}

type UserGetTokenRequest struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password,gte=8,lte=255"`
}

func UserGetToken(c *gin.Context) {
	var req UserGetTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := err.(validator.ValidationErrors)
		errs := []string{}
		for _, e := range errors {
			errs = append(errs, e.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":    "invalid",
			"errors": errs,
		})
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
		"expiresAt": sesh.ExpiresAt.Format(time.RFC3339),
	})
}
