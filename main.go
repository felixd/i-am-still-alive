package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".") // look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("Config file not found")
		} else {
			// Config file was found but another error was produced
			log.Println("Config file was found but another error was produced")
		}
	}

	c := configuration{
		app_name:      "FlameIT Dead Person Switch",
		app_namespace: "flameit",
		app_acronym:   "fitdps",

		server_port: viper.GetString("SERVER_PORT"),
		server_host: viper.GetString("SERVER_HOST"),
		switch_life: viper.GetDuration("SWITCH_LIFE"),
	}

	r := gin.Default()

	// Load data from JSON file
	if err := LoadData(); err != nil {
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
