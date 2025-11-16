package interfaces

import "github.com/gin-gonic/gin"

type UserHandler interface {
	GetUsersGetReview(c *gin.Context)
	PostUsersSetIsActive(c *gin.Context)
}

type TeamHandler interface {
	PostTeamAdd(c *gin.Context)
	GetTeamGet(c *gin.Context)
}

type PullRequestHandler interface {
	PostPullRequestCreate(c *gin.Context)
	PostPullRequestMerge(c *gin.Context)
	PostPullRequestReassign(c *gin.Context)
}
