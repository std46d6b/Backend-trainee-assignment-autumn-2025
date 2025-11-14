package repository

import (
	"context"

	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, userID domain.UserID) (domain.User, error)
	UpsertUser(ctx context.Context, user domain.User) error
	SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) error
	ListReviewPRs(ctx context.Context, userID domain.UserID) ([]domain.PullRequest, error)
}
