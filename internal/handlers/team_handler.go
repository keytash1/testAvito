package handlers

import (
	"net/http"
	"pr_reviewer_service_go/internal/models"
	"pr_reviewer_service_go/internal/services"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	svc *services.TeamService
}

func NewTeamHandler(s *services.TeamService) *TeamHandler {
	return &TeamHandler{svc: s}
}

func (h *TeamHandler) PostTeamAdd(c *gin.Context) {
	var team models.Team
	if err := c.BindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdTeam, err := h.svc.Create(&team)
	if err != nil {
		if err == models.ErrTeamExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "TEAM_EXISTS", "message": err.Error()}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"team": createdTeam})
}

func (h *TeamHandler) GetTeamGet(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_name is required"})
		return
	}

	team, err := h.svc.GetByName(teamName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "team not found"}})
		return
	}
	//TeamMembers or Users
	c.JSON(http.StatusOK, team)
}
