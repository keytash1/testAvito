package interfaces

import "pr_reviewer_service_go/internal/models"

type UserService interface {
	SetUserActive(userID string, isActive bool) (*models.User, error)
	GetByID(userID string) (*models.User, error)
	GetUserReviewPRs(userID string) ([]models.PullRequestShort, error)
}

type TeamService interface {
	Create(req *models.Team) (*models.Team, error)
	GetByName(name string) (*models.Team, error)
}

type PullRequestService interface {
	Create(prID, title string, authorId string) (models.PullRequest, error)
	MergePullRequest(pullRequestId string) (*models.PullRequest, error)
	ReassignReviewer(pullRequestId string, oldReviewerID string) (string, *models.PullRequest, error)
}

type StatsService interface {
	ReviewerStats() ([]map[string]any, error)
}
