package database

import (
	"fmt"
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.GetDsn()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database. Error: %w", err)
	}

	if check := checkExistsDB(db, config.GetDBName()); check != nil {
		fmt.Println("DATABASE NOT FOUND. CREATING....")
		if createErr := createDatabase(db, config.GetDBName()); createErr != nil {
			return nil, createErr
		}

		db, err = gorm.Open(mysql.Open(config.GetDBDsn()), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("error reconnecting to database. Error: %w", err)
		}

		if err := db.AutoMigrate(models.User{}); err != nil {
			return nil, fmt.Errorf("error migrating user. Error: %w", err)
		}

	} else {
		fmt.Println("DATABASE FOUND. CONNECTING....")
		db, err = gorm.Open(mysql.Open(config.GetDBDsn()), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("error opening database. Error: %w", err)
		}
	}

	return db, nil
}

func checkExistsDB(db *gorm.DB, dbName string) error {
	var exists string
	err := db.Raw(config.ExistsDB, dbName).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}

	if exists == "" {
		return fmt.Errorf("database not found")
	}

	return nil
}
func createDatabase(db *gorm.DB, dbName string) error {
	query := fmt.Sprintf(config.CreateDB, dbName)

	if err := db.Exec(query).Error; err != nil {
		return fmt.Errorf("error creating database %s: %w", dbName, err)
	}

	fmt.Println("DATABASE CREATED SUCCESSFULLY")
	return nil
}
