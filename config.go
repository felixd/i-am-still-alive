package main

import "time"

type configuration struct {
	app_name      string
	app_namespace string
	app_acronym   string

	database_host string
	database_port string
	database_user string
	database_pass string
	database_name string

	server_host string
	server_port string

	smtp_host       string
	smtp_port       string
	smtp_user       string
	smtp_pass       string
	smtp_from_email string
	smtp_from_label string

	data_file string

	jwt_secret_key string

	// 1814400000000000 ns -> 21 days
	switch_life time.Duration
}
