package services

import (
	"math/rand"
	"pr_reviewer_service_go/internal/models"
	"pr_reviewer_service_go/internal/repository"
	"time"
)

type PullRequestService struct {
	prRepo   *repository.PullRequestRepository
	userRepo *repository.UserRepository
	teamRepo *repository.TeamRepository
}

func NewPRService(pr *repository.PullRequestRepository, ur *repository.UserRepository, tr *repository.TeamRepository) *PullRequestService {
	return &PullRequestService{prRepo: pr, userRepo: ur, teamRepo: tr}
}

func (s *PullRequestService) Create(prID, title string, authorId string) (models.PullRequest, error) {
	author, err := s.userRepo.GetByID(authorId)
	if err != nil {
		return models.PullRequest{}, models.ErrNotFound
	}

	if _, err := s.prRepo.GetByID(prID); err == nil {
		return models.PullRequest{}, models.ErrPRExists
	}

	revs, err := s.assignReviewers(*author)
	if err != nil {
		return models.PullRequest{}, err
	}

	pr := models.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   title,
		AuthorID:          authorId,
		AssignedReviewers: revs,
		Status:            models.PullRequestStatusOPEN,
	}

	if err := s.prRepo.CreatePullRequest(&pr); err != nil {
		return pr, err
	}

	return pr, nil
}

func (s *PullRequestService) MergePullRequest(pullRequestId string) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetByID(pullRequestId)
	if err != nil {
		return nil, models.ErrNotFound
	}
	if pr.Status == models.PullRequestStatusMERGED {
		return pr, nil
	}
	now := time.Now().UTC()
	if err := s.prRepo.MergePullRequest(pullRequestId, &now); err != nil {
		return nil, err
	}
	pr.Status = models.PullRequestStatusMERGED
	pr.MergedAt = &now
	return pr, nil
}

func (s *PullRequestService) assignReviewers(author models.User) ([]string, error) {
	activeUsers, err := s.userRepo.GetActiveUsersByTeam(author.TeamName)
	if err != nil {
		return nil, err
	}

	var candidates []models.User
	for _, u := range activeUsers {
		if u.UserID != author.UserID {
			candidates = append(candidates, u)
		}
	}

	numToPick := 2
	if len(candidates) < 2 {
		numToPick = len(candidates)
	}
	if numToPick == 0 {
		return []string{}, nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	revs := make([]string, 0, 2)
	for i := 0; i < numToPick; i++ {
		revs = append(revs, candidates[i].UserID)
	}
	return revs, nil
}

func (s *PullRequestService) ReassignReviewer(pullRequestId string, oldReviewerID string) (string, *models.PullRequest, error) {
	pr, err := s.prRepo.GetByID(pullRequestId)
	if err != nil {
		return "", nil, models.ErrNotFound
	}
	if pr.Status == models.PullRequestStatusMERGED {
		return "", nil, models.ErrPRMerged
	}

	oldReviewer, err := s.userRepo.GetByID(oldReviewerID)
	if err != nil {
		return "", nil, models.ErrNotFound
	}

	activeUsers, err := s.userRepo.GetActiveUsersByTeam(oldReviewer.TeamName)
	if err != nil {
		return "", nil, err
	}

	exclude := map[string]struct{}{}
	exclude[pr.AuthorID] = struct{}{}
	for _, r := range pr.AssignedReviewers {
		exclude[r] = struct{}{}
	}

	candidates := make([]string, 0)
	for _, u := range activeUsers {
		if _, ok := exclude[u.UserID]; !ok {
			candidates = append(candidates, u.UserID)
		}
	}

	if len(candidates) == 0 {
		return "", nil, models.ErrNoCandidate
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	newReviewer := candidates[r.Intn(len(candidates))]

	replaced := false
	for i, r := range pr.AssignedReviewers {
		if r == oldReviewerID {
			pr.AssignedReviewers[i] = newReviewer
			replaced = true
			break
		}
	}
	if !replaced {
		return "", nil, models.ErrNotAssigned
	}

	if err := s.prRepo.Save(pr); err != nil {
		return "", nil, err
	}

	return newReviewer, pr, nil
}
