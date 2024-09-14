package db

import (
    "database/sql"
    "log"
    "dead-man-switch/internal/config"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Connect() {
    connStr := config.Cfg.DB.Name
    var err error
    DB, err = sql.Open("sqlite3", connStr)
    if err != nil {
        log.Fatalf("Failed to connect to the SQLite database: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Failed to ping the SQLite database: %v", err)
    }

    createTables()
}

func createTables() {
    createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`

    createSwitchTable := `
    CREATE TABLE IF NOT EXISTS dead_man_switch (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        message TEXT NOT NULL,
        trigger_at DATETIME NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`

    _, err := DB.Exec(createUsersTable)
    if err != nil {
        log.Fatalf("Failed to create users table: %v", err)
    }

    _, err = DB.Exec(createSwitchTable)
    if err != nil {
        log.Fatalf("Failed to create dead_man_switch table: %v", err)
    }
}
