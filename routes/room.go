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

	// get user from request
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no auth user",
		})
		return
	}

	user := u.(models.User)

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

	roomResult := database.GDB.Create(&newRoom)
	if roomResult.RowsAffected != 1 {
		log.Println(roomResult.Error.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// add user as mod (owner)
	newMod := models.RoomMod{
		GivenFields: models.GivenFields{
			ID: uuid.New().String(),
		},
		UserID: user.ID,
		RoomID: newRoom.ID,
		Role:   models.RoomModRoleOwner,
	}
	newModResult := database.GDB.Create(&newMod)
	if newModResult.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error creating owner record",
		})
		return
	}

	// add user to room
	newRoomUser := models.RoomUser{
		GivenFields: models.GivenFields{
			ID: uuid.New().String(),
		},
		RoomID: newRoom.ID,
		UserID: user.ID,
		Muted:  false,
	}
	newRoomUserResult := database.GDB.Create(&newRoomUser)
	if newRoomUserResult.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error creating room user record",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
