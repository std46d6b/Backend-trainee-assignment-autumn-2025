package postgres

import "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/repository"

type PostgreRepoFactory struct{}

func NewRepoFactory() *PostgreRepoFactory {
	return &PostgreRepoFactory{}
}

func (r *PostgreRepoFactory) TeamRepository(exec Execer) repository.TeamRepository {
	return NewTeamRepo(exec)
}

func (r *PostgreRepoFactory) UserRepository(exec Execer) repository.UserRepository {
	return NewUserRepo(exec)
}

func (r *PostgreRepoFactory) PullRequestRepository(exec Execer) repository.PullRequestRepository {
	return NewPullRequestRepo(exec)
}
