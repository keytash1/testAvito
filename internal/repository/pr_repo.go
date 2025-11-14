package repository

import (
	"pr_reviewer_service_go/internal/db"
	"pr_reviewer_service_go/internal/models"
	"time"
)

type PullRequestRepository struct{}

func NewPRRepository() *PullRequestRepository { return &PullRequestRepository{} }

func (r *PullRequestRepository) CreatePullRequest(pr *models.PullRequest) error {
	return db.DB.Create(pr).Error
}

func (r *PullRequestRepository) MergePullRequest(prID string, mergedAt *time.Time) error {
	return db.DB.Model(&models.PullRequest{}).
		Where("pull_request_id = ?", prID).
		Updates(map[string]interface{}{
			"status":    "MERGED",
			"merged_at": mergedAt,
		}).Error
}

func (r *PullRequestRepository) Save(pr *models.PullRequest) error {
	return db.DB.Save(pr).Error
}

func (r *PullRequestRepository) GetByID(prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	err := db.DB.Where("pull_request_id = ?", prID).First(&pr).Error
	if err != nil {
		return nil, err
	}
	return &pr, nil
}
