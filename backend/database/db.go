package database

import (
	"fmt"
	"permit-app/helper"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	host := helper.GetEnv("DB_HOST")
	port := helper.GetEnv("DB_PORT")
	user := helper.GetEnv("DB_USERNAME")
	password := helper.GetEnv("DB_PASSWORD")
	dbName := helper.GetEnv("DB_NAME")
	sslmode := helper.GetEnv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, dbName, port, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
