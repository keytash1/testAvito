package services

import (
	repoInterfaces "pr_reviewer_service_go/internal/repository/interfaces"
	serviceSnterfaces "pr_reviewer_service_go/internal/services/interfaces"
)

type statsService struct {
	prRepo   repoInterfaces.PullRequestRepository
	userRepo repoInterfaces.UserRepository
}

func NewStatsService(pr repoInterfaces.PullRequestRepository, ur repoInterfaces.UserRepository) serviceSnterfaces.StatsService {
	return &statsService{prRepo: pr, userRepo: ur}
}

func (s *statsService) ReviewerStats() ([]map[string]any, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	stats := make([]map[string]any, 0, len(users))
	for _, u := range users {
		open := s.prRepo.CountOpenReviewsByReviewer(u.UserID)
		total := s.prRepo.CountAllReviewsByReviewer(u.UserID)

		stats = append(stats, map[string]any{
			"user_id":            u.UserID,
			"open_reviews":       open,
			"total_reviews_ever": total,
		})
	}
	return stats, nil
}
