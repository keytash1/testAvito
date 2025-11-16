package interfaces

import (
	"pr_reviewer_service_go/internal/models"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(tx *gorm.DB, u *models.User) error
	SetUserActiveStatus(userID string, isActive bool) error
	GetUsersByTeam(teamName string) ([]models.User, error)
	GetActiveUsersByTeam(teamName string) ([]models.User, error)
	GetUsersReviews(userID string) ([]models.PullRequest, error)
	GetByID(userID string) (*models.User, error)
	GetAll() ([]models.User, error)
}

type TeamRepository interface {
	CreateTeam(tx *gorm.DB, t *models.Team) error
	GetTeamByName(teamName string) (*models.Team, error)
}

type PullRequestRepository interface {
	CreatePullRequest(pr *models.PullRequest) error
	MergePullRequest(prID string, mergedAt *time.Time) error
	Save(pr *models.PullRequest) error
	GetByID(prID string) (*models.PullRequest, error)
	CountOpenReviewsByReviewer(userID string) int64
	CountAllReviewsByReviewer(userID string) int64
}

type TransactionRepository interface {
	Transaction(fn func(tx *gorm.DB) error) error
}
