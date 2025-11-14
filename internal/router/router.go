package router

import (
	"pr_reviewer_service_go/internal/handlers"
	"pr_reviewer_service_go/internal/repository"
	"pr_reviewer_service_go/internal/services"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository()
	teamRepo := repository.NewTeamRepository()
	prRepo := repository.NewPRRepository()

	teamSvc := services.NewTeamService(teamRepo, userRepo)
	userSvc := services.NewUserService(userRepo)
	prSvc := services.NewPRService(prRepo, userRepo, teamRepo)

	teamH := handlers.NewTeamHandler(teamSvc)
	userH := handlers.NewUserHandler(userSvc, prRepo)
	prH := handlers.NewPullRequestHandler(prSvc)

	api := r.Group("/")
	{
		// Teams
		api.POST("/team/add", teamH.PostTeamAdd)
		api.GET("/team/get", teamH.GetTeamGet)

		// Users
		api.POST("/users/setIsActive", userH.PostUsersSetIsActive)
		api.GET("/users/getReview", userH.GetUsersGetReview)

		// PullRequests
		api.POST("/pullRequest/create", prH.PostPullRequestCreate)
		api.POST("/pullRequest/merge", prH.PostPullRequestMerge)
		api.POST("/pullRequest/reassign", prH.PostPullRequestReassign)
	}

	return r
}
