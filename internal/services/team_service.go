package services

import (
	"pr_reviewer_service_go/internal/models"
	"pr_reviewer_service_go/internal/repository"
)

type TeamService struct {
	teamRepo *repository.TeamRepository
	userRepo *repository.UserRepository
}

func NewTeamService(tr *repository.TeamRepository, ur *repository.UserRepository) *TeamService {
	return &TeamService{teamRepo: tr, userRepo: ur}
}

func (s *TeamService) Create(req *models.Team) (*models.Team, error) {

	err := s.teamRepo.CreateTeam(req)
	if err != nil {
		return nil, err
	}

	for _, member := range req.Members {
		user := &models.User{
			UserID:   member.UserId,
			Username: member.Username,
			IsActive: member.IsActive,
			TeamName: req.TeamName,
		}
		err = s.userRepo.CreateUser(user)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

func (s *TeamService) GetByName(name string) (*models.Team, error) {
	return s.teamRepo.GetTeamByName(name)
}
