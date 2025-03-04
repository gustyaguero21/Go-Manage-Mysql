package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// router params
const (
	Port    = ":8080"
	BaseURL = "/api/go-manage"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

}

func GetUser() string {
	return os.Getenv("DB_USER")
}

func GetPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func GetHost() string {
	return os.Getenv("DB_HOST")
}

func GetDBPort() string {
	return os.Getenv("DB_PORT")
}

func GetDBName() string {
	return os.Getenv("DB_NAME")
}

func GetDsn() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	return dsn
}

func GetDBDsn() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	return dsn
}
