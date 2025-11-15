package main

import (
	"log"
	"os"
	"pr_reviewer_service_go/internal/db"
	"pr_reviewer_service_go/internal/models"
	"pr_reviewer_service_go/internal/router"
)

func main() {
	db.Connect()
	// db.DB.Migrator().DropTable(&models.User{}, &models.Team{}, &models.PullRequest{})
	if err := db.DB.AutoMigrate(&models.User{}, &models.Team{}, &models.PullRequest{}); err != nil {
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
