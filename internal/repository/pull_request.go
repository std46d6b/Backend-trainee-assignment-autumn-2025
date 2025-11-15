package repository

import (
	"context"

	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type PullRequestRepository interface {
	InsertPullRequest(ctx context.Context, pullRequest domain.PullRequest) error
	GetByID(ctx context.Context, pullRequestID string) (domain.PullRequest, error)
	AddReviewer(ctx context.Context, pullRequestID string, reviewerID string) error
	RemoveReviewer(ctx context.Context, pullRequestID string, reviewerID string) error
	MergePullRequest(ctx context.Context, pullRequest domain.PullRequest) error
}
