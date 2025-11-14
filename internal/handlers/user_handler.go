package handlers

import (
	"net/http"
	"pr_reviewer_service_go/internal/repository"
	"pr_reviewer_service_go/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc    *services.UserService
	prRepo *repository.PullRequestRepository
}

func NewUserHandler(s *services.UserService, prRepo *repository.PullRequestRepository) *UserHandler {
	return &UserHandler{svc: s, prRepo: prRepo}
}

func (h *UserHandler) GetUsersGetReview(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}
	if _, err := h.svc.GetByID(userId); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "user not found"}})
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
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
