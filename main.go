package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var (
	Config *Configuration // Pointer to Configuration
)

func main() {

	Config = NewConfiguration()

	if Config.AppEnv == "development" {
		log.Println(Config)
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
	r.GET("/health")
	// Update Switch (CheckIn) using Token (URL Token Update)
	r.GET("/checkin/:token", CheckinToken)

	authorized := r.Group("/switch")
	authorized.Use(AuthMiddleware())
	{
		authorized.POST("/create", SwitchCreate)      // Create Dead Person Switch
		authorized.GET("/checkin", SwitchCheckin)     // Update Switch timeout
		authorized.GET("/token", SwitchGenerateToken) // Generate Token for Checkins (without Auth)
		authorized.DELETE("/delete", SwitchDelete)    // Remove Switch
		authorized.PUT("/update", SwitchUpdate)       // Update Switch
	}

	// Periodic check for triggered switches
	go CheckExpiredSwitches()

	r.Run(":" + Config.ServerPort)
}
