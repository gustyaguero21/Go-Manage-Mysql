package config

import (
	"fmt"
)

// router params

const (
	Port = ":8080"
)

// Urls

const (
	BaseURL = "/api/go-manage"
)

//database configs

var (
	User      = "root"
	Password  = "gagueroferra21"
	Host      = "localhost"
	DBPort    = 3306
	DBName    = "USERS_DATA"
	Dsn       = fmt.Sprintf("%s:%s@tcp(%s:%d)/", User, Password, Host, DBPort)
	DsnWithDB = Dsn + DBName
)
