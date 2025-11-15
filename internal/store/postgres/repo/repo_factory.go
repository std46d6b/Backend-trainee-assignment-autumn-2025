package postgresrepo

import (
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/repository"
	pg "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/store/postgres"
)

type PostgreRepoFactory struct{}

func NewRepoFactory() *PostgreRepoFactory {
	return &PostgreRepoFactory{}
}

func (r *PostgreRepoFactory) TeamRepository(exec pg.Execer) repository.TeamRepository {
	return NewTeamRepo(exec)
}

func (r *PostgreRepoFactory) UserRepository(exec pg.Execer) repository.UserRepository {
	return NewUserRepo(exec)
}

func (r *PostgreRepoFactory) PullRequestRepository(exec pg.Execer) repository.PullRequestRepository {
	return NewPullRequestRepo(exec)
}
