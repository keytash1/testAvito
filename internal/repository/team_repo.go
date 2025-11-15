package repository

import (
	"pr_reviewer_service_go/internal/db"
	"pr_reviewer_service_go/internal/models"

	"gorm.io/gorm"
)

type TeamRepository struct{}

func NewTeamRepository() *TeamRepository { return &TeamRepository{} }

func (r *TeamRepository) CreateTeam(tx *gorm.DB, t *models.Team) error {
	var existing models.Team
	if err := tx.Where("team_name = ?", t.TeamName).First(&existing).Error; err == nil {
		return models.ErrTeamExists
	}
	return tx.Create(t).Error
}

func (r *TeamRepository) GetTeamByName(teamName string) (*models.Team, error) {
	var team models.Team
	if err := db.DB.Where("team_name = ?", teamName).First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}
