package main

import (
	"github.com/gin-gonic/gin"
)

var (
	Config Configuration
)

func main() {

	Config := NewConfiguration()
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
		authorized.POST("/create", CreateSwitch)
		authorized.GET("/checkin", CheckinSwitch)
		authorized.DELETE("/delete", DeleteSwitch)
		authorized.PUT("/update", UpdateSwitch)
	}

	// Periodic check for triggered switches
	go CheckExpiredSwitches()

	r.Run(":8080")
}
