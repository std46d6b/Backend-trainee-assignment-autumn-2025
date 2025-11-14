package repository

import (
	"context"

	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type TeamRepository interface {
	InsertTeam(ctx context.Context, teamName domain.TeamName) error
	GetTeamWithMembers(ctx context.Context, teamName domain.TeamName) (domain.TeamUpsert, error)
}
