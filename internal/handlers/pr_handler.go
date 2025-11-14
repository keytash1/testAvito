package handlers

import (
	"net/http"
	"pr_reviewer_service_go/internal/models"
	"pr_reviewer_service_go/internal/services"

	"github.com/gin-gonic/gin"
)

type PullRequestHandler struct {
	svc *services.PullRequestService
}

func NewPullRequestHandler(s *services.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{svc: s}
}

func (h *PullRequestHandler) PostPullRequestCreate(c *gin.Context) {
	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pr, err := h.svc.Create(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		case models.ErrPRExists:
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "PR_EXISTS", "message": err.Error()}})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"pr": pr})
}

func (h *PullRequestHandler) PostPullRequestMerge(c *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pr, err := h.svc.MergePullRequest(req.PullRequestID)
	if err != nil {
		if err == models.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

func (h *PullRequestHandler) PostPullRequestReassign(c *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newReviewer, pr, err := h.svc.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		case models.ErrPRMerged:
			///??
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "PR_MERGED", "message": err.Error()}})
		case models.ErrNotAssigned:
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "NOT_ASSIGNED", "message": err.Error()}})
		case models.ErrNoCandidate:
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "NO_CANDIDATE", "message": err.Error()}})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"pr": pr, "replaced_by": newReviewer})
}
