package repository

import (
	"context"

	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type PullRequestRepository interface {
	InsertPullRequest(ctx context.Context, pullRequest domain.PullRequest) error
	GetByID(ctx context.Context, pullRequestID domain.PullRequestID) (domain.PullRequest, error)
	AddReviewer(ctx context.Context, pullRequestID domain.PullRequestID, reviewerID domain.UserID) error
	RemoveReviewer(ctx context.Context, pullRequestID domain.PullRequestID, reviewerID domain.UserID) error
	MergePullRequest(ctx context.Context, pullRequest domain.PullRequest) error
}
