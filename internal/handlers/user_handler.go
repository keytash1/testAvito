package handlers

import (
	"net/http"
	handlerInterfaces "pr_reviewer_service_go/internal/handlers/interfaces"
	"pr_reviewer_service_go/internal/models"
	repoInterfaces "pr_reviewer_service_go/internal/repository/interfaces"
	serviceInterfaces "pr_reviewer_service_go/internal/services/interfaces"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc    serviceInterfaces.UserService
	prRepo repoInterfaces.PullRequestRepository
}

func NewUserHandler(s serviceInterfaces.UserService, prRepo repoInterfaces.PullRequestRepository) handlerInterfaces.UserHandler {
	return &UserHandler{svc: s, prRepo: prRepo}
}

func (h *UserHandler) GetUsersGetReview(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}
	if _, err := h.svc.GetByID(userId); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": models.NOTFOUND, "message": err.Error()}})
		return
	}
	prs, err := h.svc.GetUserReviewPRs(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": userId, "pull_requests": prs})
}

func (h *UserHandler) PostUsersSetIsActive(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.SetUserActive(req.UserID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": models.NOTFOUND, "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
