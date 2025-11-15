package services

import (
	"pr_reviewer_service_go/internal/models"
	"pr_reviewer_service_go/internal/repository"

	"gorm.io/gorm"
)

type TeamService struct {
	teamRepo        *repository.TeamRepository
	userRepo        *repository.UserRepository
	transactionRepo *repository.TransactionRepository
}

func NewTeamService(tr *repository.TeamRepository, ur *repository.UserRepository, transRepo *repository.TransactionRepository) *TeamService {
	return &TeamService{teamRepo: tr, userRepo: ur, transactionRepo: transRepo}
}

func (s *TeamService) Create(req *models.Team) (*models.Team, error) {
	err := s.transactionRepo.Transaction(func(tx *gorm.DB) error {
		if err := s.teamRepo.CreateTeam(tx, req); err != nil {
			return err
		}

		for _, member := range req.Members {
			user := &models.User{
				UserID:   member.UserId,
				Username: member.Username,
				IsActive: member.IsActive,
				TeamName: req.TeamName,
			}

			if err := s.userRepo.CreateUser(tx, user); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *TeamService) GetByName(name string) (*models.Team, error) {
	return s.teamRepo.GetTeamByName(name)
}
