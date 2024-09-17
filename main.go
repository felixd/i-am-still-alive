package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	Config *Configuration // Pointer to Configuration
)

func main() {

	Config = NewConfiguration()

	if Config.AppEnv == "development" {
		fmt.Println(Config)
	}

	if Config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Load data from JSON file
	if err := LoadData(Config.DataFile); err != nil {
		panic(err)
	}

	// Routes
	r.POST("/signup", Signup)
	r.POST("/login", Login)

	authorized := r.Group("/switch")
	authorized.Use(AuthMiddleware())
	{
		authorized.POST("/create", CreateSwitch)   // Create Dead Person Switch
		authorized.GET("/checkin", CheckinSwitch)  // Update Switch timeout
		authorized.DELETE("/delete", DeleteSwitch) // Remove Switch
		authorized.PUT("/update", UpdateSwitch)    // Update Switch
	}

	// Periodic check for triggered switches
	go CheckExpiredSwitches()

	r.Run(":" + Config.ServerPort)
}
