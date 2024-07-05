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
	r.POST("/api/user", routes.UserCreate)
	r.GET("/api/user/token", routes.UserGetToken)

	r.POST("/api/room", middleware.AuthMiddleware, routes.RoomCreate)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))
}
