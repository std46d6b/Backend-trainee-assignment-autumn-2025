package dto

import (
	"fmt"
	"time"

	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type PullRequestDTO struct {
	ID                string   `json:"pull_request_id"`
	Name              string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type PullRequestShortDTO struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}

func PullRequestDomainToDTO(pr domain.PullRequest) PullRequestDTO {
	var createdAtPtr *string
	var mergedAtPtr *string

	if !pr.CreatedAt.IsZero() {
		s := pr.CreatedAt.Format(time.RFC3339)
		createdAtPtr = &s
	}

	if pr.MergedAt != nil {
		s := pr.MergedAt.Format(time.RFC3339)
		mergedAtPtr = &s
	}

	return PullRequestDTO{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         createdAtPtr,
		MergedAt:          mergedAtPtr,
	}
}

func PullRequestDTOToDomain(pr PullRequestDTO) (domain.PullRequest, error) {
	var (
		createdAt time.Time
		mergedAt  *time.Time
		err       error
	)

	if pr.CreatedAt != nil && *pr.CreatedAt != "" {
		createdAt, err = time.Parse(time.RFC3339, *pr.CreatedAt)
		if err != nil {
			return domain.PullRequest{}, fmt.Errorf("error parsing createdAt: %w", err)
		}
	}

	if pr.MergedAt != nil && *pr.MergedAt != "" {
		m, err2 := time.Parse(time.RFC3339, *pr.MergedAt)
		if err2 != nil {
			return domain.PullRequest{}, fmt.Errorf("error parsing mergedAt: %w", err2)
		}
		mergedAt = &m
	}

	return domain.PullRequest{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            domain.PullRequestStatus(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}, nil
}

func PullRequestShortDomainToDTO(pr domain.PullRequest) PullRequestShortDTO {
	return PullRequestShortDTO{
		ID:       pr.ID,
		Name:     pr.Name,
		AuthorID: pr.AuthorID,
		Status:   string(pr.Status),
	}
}

func PullRequestShortDTOToDomain(pr PullRequestShortDTO) domain.PullRequest {
	return domain.PullRequest{
		ID:       pr.ID,
		Name:     pr.Name,
		AuthorID: pr.AuthorID,
		Status:   domain.PullRequestStatus(pr.Status),
	}
}
