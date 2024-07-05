package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jessehorne/superchat-core/database"
	"github.com/jessehorne/superchat-core/database/models"
	"github.com/jessehorne/superchat-core/util"
	"log"
	"net/http"
)

type RoomCreateRequest struct {
	Name     string `json:"name" binding:"required,min=1"`
	Password string `json:"password"`
}

func RoomCreate(c *gin.Context) {
	var req RoomCreateRequest
	err, res := util.TryBind(&req, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// create room
	var newRoom models.Room
	newRoom.ID = uuid.New().String()
	newRoom.Name = req.Name

	if req.Password != "" {
		salt, hash := util.ProcessPassword(req.Password)
		newRoom.PasswordProtected = true
		newRoom.Password = hash
		newRoom.PasswordSalt = salt
	}

	result := database.GDB.Create(&newRoom)
	if result.RowsAffected != 1 {
		log.Println(result.Error.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}
