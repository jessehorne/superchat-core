package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jessehorne/superchat-core/database"
	"github.com/jessehorne/superchat-core/database/models"
	"github.com/jessehorne/superchat-core/util"
	"net/http"
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
