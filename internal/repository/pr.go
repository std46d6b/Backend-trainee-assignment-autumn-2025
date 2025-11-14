package repository

import (
	"context"

	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type PullRequestRepository interface {
	InsertPullRequest(ctx context.Context, pullRequest domain.PullRequest) error
	AddReviewer(ctx context.Context, pullRequestID domain.PullRequestID, reviewerID domain.UserID) error
	ClearReviewers(ctx context.Context, pullRequestID domain.PullRequestID) error
	MergePullRequest(ctx context.Context, pullRequest domain.PullRequest) error
}
