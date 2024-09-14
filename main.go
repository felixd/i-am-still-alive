package main

import (
    "dead-man-switch/internal/config"
    "dead-man-switch/internal/db"
    "dead-man-switch/handlers"
    "log"
    "net/http"
)

func main() {
    // Load configuration
    config.LoadConfig()

    // Connect to the database
    db.Connect()

    // Setup HTTP routes
    http.HandleFunc("/register", handlers.RegisterHandler)
    http.HandleFunc("/login", handlers.LoginHandler)

    // Start server
    log.Printf("Server listening on port %s", config.Cfg.Server.Port)
    log.Fatal(http.ListenAndServe(config.Cfg.Server.Port, nil))
}
