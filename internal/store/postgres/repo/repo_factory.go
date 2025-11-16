package postgresrepo

import (
	"github.com/Masterminds/squirrel"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/repository"
	pg "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/store/postgres"
)

type PostgreRepoFactory struct {
	builder squirrel.StatementBuilderType
}

func NewRepoFactory(builder squirrel.StatementBuilderType) *PostgreRepoFactory {
	return &PostgreRepoFactory{builder: builder}
}

func (r *PostgreRepoFactory) TeamRepository(exec pg.Execer) repository.TeamRepository {
	return NewTeamRepo(exec, r.builder)
}

func (r *PostgreRepoFactory) UserRepository(exec pg.Execer) repository.UserRepository {
	return NewUserRepo(exec, r.builder)
}

func (r *PostgreRepoFactory) PullRequestRepository(exec pg.Execer) repository.PullRequestRepository {
	return NewPullRequestRepo(exec, r.builder)
}
