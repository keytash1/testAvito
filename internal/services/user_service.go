package services

import (
	"pr_reviewer_service_go/internal/models"
	"pr_reviewer_service_go/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService { return &UserService{repo: r} }

func (s *UserService) SetUserActive(userID string, isActive bool) (*models.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, models.ErrNotFound
	}

	if err := s.repo.SetUserActiveStatus(userID, isActive); err != nil {
		return nil, err
	}

	user.IsActive = isActive
	return user, nil
}

func (s *UserService) GetByID(userID string) (*models.User, error) {
	return s.repo.GetByID(userID)
}

func (s *UserService) GetUserReviewPRs(userID string) ([]models.PullRequestShort, error) {
	longprs, err := s.repo.GetUsersReviews(userID)
	if err != nil {
		return nil, err
	}
	shortprs := make([]models.PullRequestShort, 0, len(longprs))
	for _, longpr := range longprs {
		shortprs = append(shortprs, models.PullRequestShort{
			AuthorId:        longpr.AuthorID,
			PullRequestId:   longpr.PullRequestID,
			PullRequestName: longpr.PullRequestName,
			Status:          longpr.Status,
		})
	}
	return shortprs, nil
}
