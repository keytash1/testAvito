package models

import (
	"errors"
	"time"
)

type ErrorResponseErrorCode string

const (
	NOCANDIDATE ErrorResponseErrorCode = "NO_CANDIDATE"
	NOTASSIGNED ErrorResponseErrorCode = "NOT_ASSIGNED"
	NOTFOUND    ErrorResponseErrorCode = "NOT_FOUND"
	PREXISTS    ErrorResponseErrorCode = "PR_EXISTS"
	PRMERGED    ErrorResponseErrorCode = "PR_MERGED"
	TEAMEXISTS  ErrorResponseErrorCode = "TEAM_EXISTS"
)

var (
	ErrTeamExists  = errors.New("team already exists")
	ErrPRExists    = errors.New("PR id already exists")
	ErrPRMerged    = errors.New("cannot reassign on merged PR")
	ErrNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate in team")
	ErrNotFound    = errors.New("resource not found")
)

type PullRequest struct {
	PullRequestID     string            `json:"pull_request_id" gorm:"primaryKey;type:varchar(100)"`
	PullRequestName   string            `json:"pull_request_name" gorm:"not null"`
	AuthorID          string            `json:"author_id" gorm:"index;not null"`
	Status            PullRequestStatus `json:"status" gorm:"type:varchar(20);not null"` // OPEN | MERGED
	AssignedReviewers []string          `json:"assigned_reviewers" gorm:"type:jsonb;serializer:json"`
	CreatedAt         time.Time         `json:"createdAt"`
	MergedAt          *time.Time        `json:"mergedAt,omitempty"`
}

type PullRequestStatus string

const (
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
)

type PullRequestShort struct {
	AuthorId        string            `json:"author_id"`
	PullRequestId   string            `json:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name"`
	Status          PullRequestStatus `json:"status"`
}

type Team struct {
	TeamName string       `json:"team_name" gorm:"primaryKey;type:varchar(100)"`
	Members  []TeamMember `json:"members" gorm:"type:jsonb;serializer:json"`
}

type TeamMember struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type User struct {
	UserID   string `json:"user_id" gorm:"primaryKey;type:varchar(100)"`
	Username string `json:"username" gorm:"not null"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
	TeamName string `json:"team_name" gorm:"index"`
}
