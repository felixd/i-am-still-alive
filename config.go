package main

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	AppEnv       string `mapstructure:"APP_ENV"`
	AppName      string
	AppNamespace string
	AppAcronym   string

	DatabaseHost string `mapstructure:"DATABASE_HOST"`
	DatabasePort string `mapstructure:"DATABASE_PORT"`
	DatabaseUser string `mapstructure:"DATABASE_USER"`
	DatabasePass string `mapstructure:"DATABASE_PASS"`
	DatabaseName string `mapstructure:"DATABASE_NAME"`

	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort string `mapstructure:"SERVER_PORT"`

	SmtpHost      string `mapstructure:"SMTP_HOST"`
	SmtpPort      string `mapstructure:"SMTP_PORT"`
	SmtpUser      string `mapstructure:"SMTP_USER"`
	SmtpPass      string `mapstructure:"SMTP_PASS"`
	SmtpFromEmail string `mapstructure:"SMTP_FROM_EMAIL"`
	SmtpFromLabel string `mapstructure:"SMTP_FROM_LABEL"`

	DataFile string `mapstructure:"DATA_FILE"`

	JwtSecretKey string `mapstructure:"JWT_SECRET_KEY"`

	// 1814400000000000 ns -> 21 days
	SwitchDefaultLife time.Duration
}

func NewConfiguration() *Configuration {

	c := Configuration{}

	c.AppName = "FlameIT Dead Person Switch"
	c.AppNamespace = "flameit"
	c.AppAcronym = "fitdps"
	c.AppEnv = "development"
	c.DataFile = "data.json"

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".") // look for config in the working directory
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	t, err := time.ParseDuration(viper.GetString("SWITCH_DEFAULT_LIFE"))
	if err != nil {
		log.Fatal("Error parsing Switch Default Life Time : ", err)
	}

	if t > 0 {
		c.SwitchDefaultLife = t
	}

	if c.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &c

}
