package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jessehorne/superchat-core/database"
	"github.com/jessehorne/superchat-core/middleware"
	"github.com/jessehorne/superchat-core/routes"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	if _, err := database.InitDB(); err != nil {
		panic(err)
	}

	if _, err := database.InitGDB(); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/api/ping", routes.GetPing)

	/* User Routes */
	r.POST("/api/user", routes.UserCreate)
	r.GET("/api/user/token", routes.UserGetToken)
	r.PUT("/api/user", middleware.AuthMiddleware, routes.UserUpdate)

	/* Room Routes */
	r.POST("/api/room", middleware.AuthMiddleware, routes.RoomCreate)
	r.PUT("/api/room", middleware.AuthMiddleware, routes.RoomUpdate)
	r.DELETE("/api/room", middleware.AuthMiddleware, routes.RoomDelete)
	r.POST("/api/room/mod", middleware.AuthMiddleware, routes.RoomAddMod)
	r.PUT("/api/room/mod", middleware.AuthMiddleware, routes.RoomUpdateMod)
	r.DELETE("/api/room/mod", middleware.AuthMiddleware, routes.RoomDeleteMod)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))
}
