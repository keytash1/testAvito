package repository

import (
	"pr_reviewer_service_go/internal/db"
	"pr_reviewer_service_go/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(tx *gorm.DB, u *models.User) error {
	return tx.Create(u).Error
}

func (r *UserRepository) SetUserActiveStatus(userID string, isActive bool) error {
	return db.DB.Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("is_active", isActive).Error
}

func (r *UserRepository) GetUsersByTeam(teamName string) ([]models.User, error) {
	var users []models.User
	err := db.DB.Where("team_name = ?", teamName).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetActiveUsersByTeam(teamName string) ([]models.User, error) {
	var users []models.User
	err := db.DB.Where("team_name = ? AND is_active = ?", teamName, true).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetUsersReviews(userID string) ([]models.PullRequest, error) {
	var pullRequests []models.PullRequest
	err := db.DB.
		Where("assigned_reviewers @> ? AND status = ?", `["`+userID+`"]`, models.PullRequestStatusOPEN).
		Find(&pullRequests).Error
	return pullRequests, err
}

func (r *UserRepository) GetByID(userID string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
