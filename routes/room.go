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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "db error",
		})
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

	c.JSON(http.StatusOK, gin.H{
		"roomID": newRoom.ID,
	})
}

type RoomUpdateRequest struct {
	RoomID   string `json:"roomID" binding:"required"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func RoomUpdate(c *gin.Context) {
	var req RoomUpdateRequest
	err, res := util.TryBind(&req, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if req.RoomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing roomID",
		})
		return
	}

	// get room or return error if it doesn't exist
	var room models.Room
	roomResult := database.GDB.First(&room, "id = ?", req.RoomID)
	if roomResult.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "room not found",
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

	// make sure user is an owner
	var roomMod models.RoomMod
	roomModResult := database.GDB.Where("room_id = ?", req.RoomID).First(&roomMod, "user_id = ?", user.ID)
	if roomModResult.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "are you even a mod bro",
		})
		return
	}

	if roomMod.Role != models.RoomModRoleOwner {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you're not the owner",
		})
		return
	}

	// update room name if the field was found
	if req.Name != "" {
		room.Name = req.Name
	}

	// update room password if password field was given
	if req.Password != "" {
		salt, hash := util.ProcessPassword(req.Password)
		room.PasswordProtected = true
		room.Password = hash
		room.PasswordSalt = salt
	}

	updateResult := database.GDB.Save(&room)
	if updateResult.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong with the db",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

type RoomDeleteRequest struct {
	RoomID string `json:"roomID"`
}

func RoomDelete(c *gin.Context) {
	var req RoomDeleteRequest
	err, res := util.TryBind(&req, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if req.RoomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing roomID",
		})
		return
	}

	// get room or return error if it doesn't exist
	var room models.Room
	roomResult := database.GDB.First(&room, "id = ?", req.RoomID)
	if roomResult.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "room not found",
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

	// make sure user is an owner
	var roomMod models.RoomMod
	roomModResult := database.GDB.Where("room_id = ?", req.RoomID).First(&roomMod, "user_id = ?", user.ID)
	if roomModResult.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "are you even a mod bro",
		})
		return
	}

	if roomMod.Role != models.RoomModRoleOwner {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you're not the owner",
		})
		return
	}

	// delete room
	deleteResult := database.GDB.Delete(&room)
	if deleteResult.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong deleting the room",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

type RoomAddModRequest struct {
	RoomID string `json:"roomID"`
	UserID string `json:"userID"`
	Role   int    `json:"role,default=-1"`
}

func RoomAddMod(c *gin.Context) {
	var req RoomAddModRequest
	err, res := util.TryBind(&req, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if req.RoomID == "" || req.UserID == "" || req.Role == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing details",
		})
		return
	}

	if req.Role != 0 && req.Role != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid role",
		})
		return
	}

	// get room or return error if it doesn't exist
	var room models.Room
	roomResult := database.GDB.First(&room, "id = ?", req.RoomID)
	if roomResult.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "room not found",
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

	// make sure user is an owner
	var roomMod models.RoomMod
	roomModResult := database.GDB.Where("room_id = ?", req.RoomID).First(&roomMod, "user_id = ?", user.ID)
	if roomModResult.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "are you even a mod bro",
		})
		return
	}

	// make sure a mod doesn't already exist
	var existingRoomMod models.RoomMod
	existingRoomModResult := database.GDB.Where("room_id = ?", req.RoomID).First(&existingRoomMod, "user_id = ?", req.UserID)
	if existingRoomModResult.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "mod already exists",
		})
		return
	}

	// get target user
	var targetUser models.User
	targetUserResult := database.GDB.First(&targetUser, "id = ?", req.UserID)
	if targetUserResult.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no target user",
		})
		return
	}

	newRoomMod := models.RoomMod{
		GivenFields: models.GivenFields{
			ID: uuid.New().String(),
		},
		UserID: req.UserID,
		RoomID: req.RoomID,
		Role:   req.Role,
	}

	newRoomModResult := database.GDB.Create(&newRoomMod)
	if newRoomModResult.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error saving room mod",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
