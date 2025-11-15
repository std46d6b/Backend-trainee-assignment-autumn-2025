package teamservice

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
	TeamRepository(exec postgres.Execer) repository.TeamRepository
	UserRepository(exec postgres.Execer) repository.UserRepository
}

type TeamService struct {
	txManager TxManager
	repoFact  RepoFactory
	readExec  postgres.Execer
}

func NewTeamService(
	txManager TxManager,
	readExec postgres.Execer,
	repoFact RepoFactory,
) *TeamService {
	return &TeamService{
		txManager: txManager,
		repoFact:  repoFact,
		readExec:  readExec,
	}
}

// POST /team/add
func (s *TeamService) CreateTeam(ctx context.Context, up domain.TeamUpsert) (domain.TeamUpsert, error) {
	err := s.txManager.TxWrapper(ctx, func(ctx context.Context, tx pgx.Tx) error {
		localTeamRepo := s.repoFact.TeamRepository(tx)
		localUserRepo := s.repoFact.UserRepository(tx)

		err := localTeamRepo.InsertTeam(ctx, up.Name)
		if err != nil {
			return fmt.Errorf("insert team: %w", err)
		}

		for _, member := range up.Members {
			upsUser := domain.User{
				ID:       member.UserID,
				Username: member.Username,
				TeamName: up.Name,
				IsActive: member.IsActive,
			}

			err = localUserRepo.UpsertUser(ctx, upsUser)
			if err != nil {
				return fmt.Errorf("upsert user %s: %w", member.UserID, err)
			}
		}

		return nil
	})

	if err != nil {
		return domain.TeamUpsert{}, fmt.Errorf("create team: %w", err)
	}

	return up, nil
}

// GET /team/get
func (s *TeamService) GetTeamWithMembers(ctx context.Context, teamName string) (domain.TeamUpsert, error) {
	localTeamRepo := s.repoFact.TeamRepository(s.readExec)

	domainTeam, err := localTeamRepo.GetTeamWithMembers(ctx, teamName)

	if err != nil {
		return domain.TeamUpsert{}, fmt.Errorf("service get team: %w", err)
	}

	return domainTeam, nil
}
