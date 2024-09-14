package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    DB struct {
        Name string
    }
    Server struct {
        Port string
    }
}

var Cfg Config

func LoadConfig() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    Cfg.DB.Name = os.Getenv("DB_NAME") // Path to SQLite DB file
    Cfg.Server.Port = os.Getenv("SERVER_PORT")
}
