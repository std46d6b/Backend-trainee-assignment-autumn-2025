package userservice

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/repository"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/store/postgres"
)

type TxManager interface {
	TxWrapper(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error
}

type RepoFactory interface {
	UserRepository(exec postgres.Execer) repository.UserRepository
}

type UserService struct {
	txManager TxManager
	repoFact  RepoFactory
	readExec  postgres.Execer
}

func NewUserService(
	txManager TxManager,
	readExec postgres.Execer,
	repoFact RepoFactory,
) *UserService {
	return &UserService{
		txManager: txManager,
		repoFact:  repoFact,
		readExec:  readExec,
	}
}

// POST /users/setIsActive
func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	var dbUser domain.User

	err := s.txManager.TxWrapper(ctx, func(ctx context.Context, tx pgx.Tx) error {
		localUserRepo := s.repoFact.UserRepository(tx)

		err := localUserRepo.SetIsActive(ctx, userID, isActive)
		if err != nil {
			return fmt.Errorf("service set is active: %w", err)
		}

		dbUser, err = localUserRepo.GetByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("service get by id: %w", err)
		}

		return nil
	})

	if err != nil {
		return dbUser, fmt.Errorf("service set is active transaction: %w", err)
	}

	return dbUser, nil
}

// GET /users/getReview
func (s *UserService) ListReviewPRs(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	locaUserRepo := s.repoFact.UserRepository(s.readExec)

	pullRequests, err := locaUserRepo.ListReviewPRs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service list review prs: %w", err)
	}

	return pullRequests, nil
}
