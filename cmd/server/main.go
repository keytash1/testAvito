package main

import (
	"log"
	"pr_reviewer_service_go/internal/db"
	"pr_reviewer_service_go/internal/models"
	"pr_reviewer_service_go/internal/router"
)

func main() {
	db.Connect()
	//db.DB.Migrator().DropTable(&models.User{}, &models.Team{}, &models.PullRequest{})
	if err := db.DB.AutoMigrate(&models.User{}, &models.Team{}, &models.PullRequest{}); err != nil {
		log.Fatal("migrate:", err)
	}

	r := router.New()
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

//
