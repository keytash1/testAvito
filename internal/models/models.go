package models

import (
	"errors"
	"time"
)

type ErrorResponse struct {
	Error struct {
		Code    error  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

var (
	ErrNotFound    = errors.New("NOT_FOUND")
	ErrPRExists    = errors.New("PR_EXISTS")
	ErrPRMerged    = errors.New("PR_MERGED")
	ErrNoCandidate = errors.New("NO_CANDIDATE")
	ErrNotAssigned = errors.New("NOT_ASSIGNED")
	ErrTeamExists  = errors.New("TEAM_EXISTS")
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
