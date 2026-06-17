package config

import (
	"fmt"
	"log"
	"os"

	"backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		sslmode,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erreur connexion PostgreSQL :", err)
	}

	DB = database

	err = DB.AutoMigrate(
		&models.User{},
		&models.Door{},
		&models.FaceProfile{},
		&models.AccessPermission{},
		&models.AccessLog{},
	)
	if err != nil {
		log.Fatal("Erreur migration :", err)
	}

	log.Println("Base de données connectée")
}
