package main

import (
	"log"
	"os"
	"pr_reviewer_service_go/internal/db"
	"pr_reviewer_service_go/internal/models"
	"pr_reviewer_service_go/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Println("Warning: error loading .env file:", err)
	}
	db.Connect()

	if os.Getenv("IS_TEST") == "true" {
		if err := db.DB.Migrator().DropTable(
			&models.PullRequest{},
			&models.User{},
			&models.Team{},
		); err != nil {
			log.Fatal("failed to drop tables:", err)
		}
	}

	if err := db.DB.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.PullRequest{},
	); err != nil {
		log.Fatal("migrate:", err)
	}

	r := router.New()
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		log.Fatal("SERVER_URL environment variable is required")
	}

	if err := r.Run(serverURL); err != nil {
		log.Fatal(err)
	}
}
